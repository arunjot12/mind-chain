//! Optimism Consensus implementation.

#![doc(
    html_logo_url = "https://raw.githubusercontent.com/paradigmxyz/reth/main/assets/reth-docs.png",
    html_favicon_url = "https://avatars0.githubusercontent.com/u/97369466?s=256",
    issue_tracker_base_url = "https://github.com/paradigmxyz/reth/issues/"
)]
#![cfg_attr(docsrs, feature(doc_cfg))]
#![cfg_attr(not(feature = "std"), no_std)]
#![cfg_attr(not(test), warn(unused_crate_dependencies))]

extern crate alloc;

use alloc::{format, sync::Arc};
use alloy_consensus::{
    BlockHeader as _, EMPTY_OMMER_ROOT_HASH, constants::MAXIMUM_EXTRA_DATA_SIZE,
};
use alloy_primitives::{B64, B256};
use core::fmt::Debug;
use op_alloy_consensus::{OpTransaction, decode_holocene_extra_data, decode_jovian_extra_data};
use reth_chainspec::EthChainSpec;
use reth_consensus::{Consensus, ConsensusError, FullConsensus, HeaderValidator, ReceiptRootBloom};
use reth_consensus_common::validation::{
    validate_against_parent_eip1559_base_fee, validate_against_parent_hash_number,
    validate_against_parent_timestamp, validate_cancun_gas, validate_header_base_fee,
    validate_header_extra_data, validate_header_gas,
};
use reth_execution_types::BlockExecutionResult;
use reth_mind_chain_forks::{OpHardfork, OpHardforks};
use reth_mind_chain_primitives::DepositReceipt;
use reth_primitives_traits::{
    Block, BlockBody, BlockHeader, GotExpected, NodePrimitives, RecoveredBlock, SealedBlock,
    SealedHeader,
};

mod proof;
pub use proof::calculate_receipt_root_optimism;

#[cfg(feature = "std")]
mod sdm_metrics;

pub mod validation;
pub use validation::{canyon, isthmus, validate_block_post_execution};

pub mod error;
pub use error::OpConsensusError;

/// Validates that a header's `extraData` has the encoding required by the fork active at its
/// timestamp: the Jovian or Holocene encoding, or empty before Holocene.
///
/// Before Holocene, the rollup genesis block may carry arbitrary `extraData`. That is the Bedrock
/// transition block: block 0, or on chains migrated from a legacy chain (OP Mainnet) the first
/// Bedrock block, whose chain-spec genesis is the legacy block 0.
fn validate_op_header_extra_data<H, ChainSpec>(
    header: &H,
    chain_spec: &ChainSpec,
) -> Result<(), ConsensusError>
where
    H: BlockHeader,
    ChainSpec: OpHardforks,
{
    if chain_spec.is_jovian_active_at_timestamp(header.timestamp()) {
        decode_jovian_extra_data(header.extra_data()).map_err(ConsensusError::other)?;
    } else if chain_spec.is_holocene_active_at_timestamp(header.timestamp()) {
        decode_holocene_extra_data(header.extra_data()).map_err(ConsensusError::other)?;
    } else if !header.extra_data().is_empty() &&
        !chain_spec.op_fork_activation(OpHardfork::Bedrock).transitions_at_block(header.number())
    {
        return Err(ConsensusError::msg("extraData must be empty before Holocene"));
    }

    Ok(())
}

/// Optimism consensus implementation.
///
/// Provides basic checks as outlined in the execution specs.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct OpBeaconConsensus<ChainSpec> {
    /// Configuration
    chain_spec: Arc<ChainSpec>,
    /// Maximum allowed extra data size in bytes
    max_extra_data_size: usize,
}

impl<ChainSpec> OpBeaconConsensus<ChainSpec> {
    /// Create a new instance of [`OpBeaconConsensus`]
    pub const fn new(chain_spec: Arc<ChainSpec>) -> Self {
        Self { chain_spec, max_extra_data_size: MAXIMUM_EXTRA_DATA_SIZE }
    }

    /// Returns the maximum allowed extra data size.
    pub const fn max_extra_data_size(&self) -> usize {
        self.max_extra_data_size
    }

    /// Sets the maximum allowed extra data size and returns the updated instance.
    pub const fn with_max_extra_data_size(mut self, size: usize) -> Self {
        self.max_extra_data_size = size;
        self
    }
}

impl<N, ChainSpec> FullConsensus<N> for OpBeaconConsensus<ChainSpec>
where
    N: NodePrimitives<Receipt: DepositReceipt, SignedTx: OpTransaction>,
    ChainSpec: EthChainSpec<Header = N::BlockHeader> + OpHardforks + Debug + Send + Sync,
{
    fn validate_block_post_execution(
        &self,
        block: &RecoveredBlock<N::Block>,
        result: &BlockExecutionResult<N::Receipt>,
        receipt_root_bloom: Option<ReceiptRootBloom>,
        _block_access_list_hash: Option<B256>,
    ) -> Result<(), ConsensusError> {
        validate_block_post_execution(
            block.header(),
            &self.chain_spec,
            result,
            receipt_root_bloom,
        )?;
        #[cfg(feature = "std")]
        sdm_metrics::record_validated_block(block.header().gas_used(), block.body().transactions());
        Ok(())
    }
}

impl<B, ChainSpec> Consensus<B> for OpBeaconConsensus<ChainSpec>
where
    B: Block,
    ChainSpec: EthChainSpec<Header = B::Header> + OpHardforks + Debug + Send + Sync,
{
    fn validate_body_against_header(
        &self,
        body: &B::Body,
        header: &SealedHeader<B::Header>,
    ) -> Result<(), ConsensusError> {
        validation::validate_body_against_header_op(&self.chain_spec, body, header.header())
    }

    fn validate_block_pre_execution(&self, block: &SealedBlock<B>) -> Result<(), ConsensusError> {
        // Check ommers hash
        let ommers_hash = block.body().calculate_ommers_root();
        if Some(block.ommers_hash()) != ommers_hash {
            return Err(ConsensusError::BodyOmmersHashDiff(
                GotExpected {
                    got: ommers_hash.unwrap_or(EMPTY_OMMER_ROOT_HASH),
                    expected: block.ommers_hash(),
                }
                .into(),
            ));
        }

        // Check transaction root
        if let Err(error) = block.ensure_transaction_root_valid() {
            return Err(ConsensusError::BodyTransactionRootDiff(error.into()));
        }

        // Check empty shanghai-withdrawals
        if self.chain_spec.is_canyon_active_at_timestamp(block.timestamp()) {
            canyon::ensure_empty_shanghai_withdrawals(block.body()).map_err(|err| {
                ConsensusError::msg(format!("failed to verify block {}: {err}", block.number()))
            })?
        } else {
            return Ok(());
        }

        // Blob gas used validation
        // In Jovian, the blob gas used computation has changed. We are moving the blob base fee
        // validation to post-execution since the DA footprint calculation is stateful.
        // Pre-execution we only validate that the blob gas used is present in the header.
        if self.chain_spec.is_jovian_active_at_timestamp(block.timestamp()) {
            block.blob_gas_used().ok_or(ConsensusError::BlobGasUsedMissing)?;
        } else if self.chain_spec.is_ecotone_active_at_timestamp(block.timestamp()) {
            validate_cancun_gas(block)?;
        }

        // Check withdrawals root field in header
        if self.chain_spec.is_isthmus_active_at_timestamp(block.timestamp()) {
            // storage root of withdrawals pre-deploy is verified post-execution
            isthmus::ensure_withdrawals_storage_root_is_some(block.header()).map_err(|err| {
                ConsensusError::msg(format!("failed to verify block {}: {err}", block.number()))
            })?
        } else {
            // canyon is active, else would have returned already
            canyon::ensure_empty_withdrawals_root(block.header())?
        }

        Ok(())
    }
}

impl<H, ChainSpec> HeaderValidator<H> for OpBeaconConsensus<ChainSpec>
where
    H: BlockHeader,
    ChainSpec: EthChainSpec<Header = H> + OpHardforks + Debug + Send + Sync,
{
    fn validate_header(&self, header: &SealedHeader<H>) -> Result<(), ConsensusError> {
        let header = header.header();
        // with OP-stack Bedrock activation number determines when TTD (eth Merge) has been reached.
        debug_assert!(
            self.chain_spec.is_bedrock_active_at_block(header.number()),
            "manually import OVM blocks"
        );

        if header.nonce() != Some(B64::ZERO) {
            return Err(ConsensusError::TheMergeNonceIsNotZero);
        }

        if header.ommers_hash() != EMPTY_OMMER_ROOT_HASH {
            return Err(ConsensusError::TheMergeOmmerRootIsNotEmpty);
        }

        // Post-merge, the consensus layer is expected to perform checks such that the block
        // timestamp is a function of the slot. This is different from pre-merge, where blocks
        // are only allowed to be in the future (compared to the system's clock) by a certain
        // threshold.
        //
        // Block validation with respect to the parent should ensure that the block timestamp
        // is greater than its parent timestamp.

        // validate header extra data for all networks post merge
        validate_header_extra_data(header, self.max_extra_data_size)?;
        validate_op_header_extra_data(header, &self.chain_spec)?;
        validate_header_gas(header)?;
        validate_header_base_fee(header, &self.chain_spec)
    }

    fn validate_header_against_parent(
        &self,
        header: &SealedHeader<H>,
        parent: &SealedHeader<H>,
    ) -> Result<(), ConsensusError> {
        validate_against_parent_hash_number(header.header(), parent)?;

        if self.chain_spec.is_bedrock_active_at_block(header.number()) {
            validate_against_parent_timestamp(header.header(), parent.header())?;
        }

        validate_against_parent_eip1559_base_fee(
            header.header(),
            parent.header(),
            &self.chain_spec,
        )?;

        // Ensure that the blob gas fields for this block are correctly set.
        // In the op-stack, the excess blob gas is always 0 for all blocks after ecotone.
        // The blob gas used and the excess blob gas should both be set after ecotone.
        // After Jovian, the blob gas used contains the current DA footprint.
        if self.chain_spec.is_ecotone_active_at_timestamp(header.timestamp()) {
            let blob_gas_used = header.blob_gas_used().ok_or(ConsensusError::BlobGasUsedMissing)?;

            // Before Jovian and after ecotone, the blob gas used should be 0.
            if !self.chain_spec.is_jovian_active_at_timestamp(header.timestamp()) &&
                blob_gas_used != 0
            {
                return Err(ConsensusError::BlobGasUsedDiff(GotExpected {
                    got: blob_gas_used,
                    expected: 0,
                }));
            }

            let excess_blob_gas =
                header.excess_blob_gas().ok_or(ConsensusError::ExcessBlobGasMissing)?;
            if excess_blob_gas != 0 {
                return Err(ConsensusError::ExcessBlobGasDiff {
                    diff: GotExpected { got: excess_blob_gas, expected: 0 },
                    parent_excess_blob_gas: parent.excess_blob_gas().unwrap_or(0),
                    parent_blob_gas_used: parent.blob_gas_used().unwrap_or(0),
                });
            }
        }

        Ok(())
    }
}
