//! Verification of blocks w.r.t. Optimism hardforks.

pub mod canyon;
pub mod isthmus;

// Re-export the decode_holocene_base_fee function for compatibility
use reth_execution_types::BlockExecutionResult;
pub use reth_mind_chain_chainspec::decode_holocene_base_fee;

use crate::proof::calculate_receipt_root_optimism;
use alloc::vec::Vec;
use alloy_consensus::{BlockHeader, EMPTY_OMMER_ROOT_HASH, TxReceipt};
use alloy_eips::Encodable2718;
use alloy_primitives::{B256, Bloom, Bytes};
use alloy_trie::EMPTY_ROOT_HASH;
use reth_consensus::ConsensusError;
use reth_ethereum_consensus::compare_receipts_root_and_logs_bloom;
use reth_mind_chain_forks::OpHardforks;
use reth_mind_chain_primitives::DepositReceipt;
use reth_primitives_traits::{BlockBody, GotExpected, receipt::gas_spent_by_transactions};

/// UPSTREAM-MIRROR(copy): reth@rev:8d359d0
/// `reth_consensus_common::validation::validate_body_against_header`
///
/// Copies the generic body/header checks and replaces withdrawals-root handling after Isthmus.
/// Ensures the block response data matches the header.
///
/// This ensures the body response items match the header's hashes:
///   - ommer hash
///   - transaction root
///   - withdrawals root: the body's withdrawals root must only match the header's before isthmus
pub fn validate_body_against_header_op<B, H>(
    chain_spec: impl OpHardforks,
    body: &B,
    header: &H,
) -> Result<(), ConsensusError>
where
    B: BlockBody,
    H: reth_primitives_traits::BlockHeader,
{
    let ommers_hash = body.calculate_ommers_root();
    if Some(header.ommers_hash()) != ommers_hash {
        return Err(ConsensusError::BodyOmmersHashDiff(
            GotExpected {
                got: ommers_hash.unwrap_or(EMPTY_OMMER_ROOT_HASH),
                expected: header.ommers_hash(),
            }
            .into(),
        ));
    }

    let tx_root = body.calculate_tx_root();
    if header.transactions_root() != tx_root {
        return Err(ConsensusError::BodyTransactionRootDiff(
            GotExpected { got: tx_root, expected: header.transactions_root() }.into(),
        ));
    }

    match (header.withdrawals_root(), body.calculate_withdrawals_root()) {
        (Some(header_withdrawals_root), Some(withdrawals_root)) => {
            // after isthmus, the withdrawals root field is repurposed and no longer mirrors the
            // withdrawals root computed from the body
            if chain_spec.is_isthmus_active_at_timestamp(header.timestamp()) {
                // After isthmus we only ensure that the body has empty withdrawals
                if withdrawals_root != EMPTY_ROOT_HASH {
                    return Err(ConsensusError::BodyWithdrawalsRootDiff(
                        GotExpected { got: withdrawals_root, expected: EMPTY_ROOT_HASH }.into(),
                    ));
                }
            } else {
                // before isthmus we ensure that the header root matches the body
                if withdrawals_root != header_withdrawals_root {
                    return Err(ConsensusError::BodyWithdrawalsRootDiff(
                        GotExpected { got: withdrawals_root, expected: header_withdrawals_root }
                            .into(),
                    ));
                }
            }
        }
        (None, None) => {
            // this is ok because we assume the fork is not active in this case
        }
        _ => return Err(ConsensusError::WithdrawalsRootUnexpected),
    }

    Ok(())
}

/// UPSTREAM-MIRROR(copy): reth@rev:8d359d0
/// `reth_ethereum_consensus::validation::validate_block_post_execution`
///
/// Copies upstream post-execution validation with OP receipt roots and Jovian DA accounting.
/// Validate a block with regard to execution results:
///
/// - Compares the receipts root in the block header to the block body
/// - Compares the gas used in the block header to the actual gas usage after execution
///
/// If `receipt_root_bloom` is provided, the pre-computed receipt root and logs bloom are used
/// instead of computing them from the receipts.
pub fn validate_block_post_execution<R: DepositReceipt>(
    header: impl BlockHeader,
    chain_spec: impl OpHardforks,
    result: &BlockExecutionResult<R>,
    receipt_root_bloom: Option<(B256, Bloom)>,
) -> Result<(), ConsensusError> {
    // Validate that the blob gas used is present and correctly computed if Jovian is active.
    if chain_spec.is_jovian_active_at_timestamp(header.timestamp()) {
        let computed_blob_gas_used = result.blob_gas_used;
        let header_blob_gas_used =
            header.blob_gas_used().ok_or(ConsensusError::BlobGasUsedMissing)?;

        if computed_blob_gas_used != header_blob_gas_used {
            return Err(ConsensusError::BlobGasUsedDiff(GotExpected {
                got: computed_blob_gas_used,
                expected: header_blob_gas_used,
            }));
        }
    }

    let receipts = &result.receipts;

    // Before Byzantium, receipts contained state root that would mean that expensive
    // operation as hashing that is required for state root got calculated in every
    // transaction This was replaced with is_success flag.
    // See more about EIP here: https://eips.ethereum.org/EIPS/eip-658
    if chain_spec.is_byzantium_active_at_block(header.number()) {
        let result = if let Some((receipts_root, logs_bloom)) = receipt_root_bloom {
            compare_receipts_root_and_logs_bloom(
                receipts_root,
                logs_bloom,
                header.receipts_root(),
                header.logs_bloom(),
            )
        } else {
            verify_receipts_optimism(
                header.receipts_root(),
                header.logs_bloom(),
                receipts,
                chain_spec,
                header.timestamp(),
            )
        };

        if let Err(error) = result {
            let receipts = receipts
                .iter()
                .map(|r| Bytes::from(r.with_bloom_ref().encoded_2718()))
                .collect::<Vec<_>>();
            tracing::debug!(%error, ?receipts, "receipts verification failed");
            return Err(error);
        }
    }

    // Check if gas used matches the value set in header.
    let cumulative_gas_used =
        receipts.last().map(|receipt| receipt.cumulative_gas_used()).unwrap_or(0);
    if header.gas_used() != cumulative_gas_used {
        return Err(ConsensusError::BlockGasUsed {
            gas: GotExpected { got: cumulative_gas_used, expected: header.gas_used() },
            gas_spent_by_tx: gas_spent_by_transactions(receipts),
        });
    }

    Ok(())
}

/// UPSTREAM-MIRROR(copy): reth@rev:8d359d0 `reth_ethereum_consensus::validation::verify_receipts`
///
/// Copies the upstream helper, substituting OP receipt-root calculation.
/// Verify the calculated receipts root against the expected receipts root.
fn verify_receipts_optimism<R: DepositReceipt>(
    expected_receipts_root: B256,
    expected_logs_bloom: Bloom,
    receipts: &[R],
    chain_spec: impl OpHardforks,
    timestamp: u64,
) -> Result<(), ConsensusError> {
    // Calculate receipts root.
    let receipts_with_bloom = receipts.iter().map(TxReceipt::with_bloom_ref).collect::<Vec<_>>();
    let receipts_root =
        calculate_receipt_root_optimism(&receipts_with_bloom, chain_spec, timestamp);

    // Calculate header logs bloom.
    let logs_bloom = receipts_with_bloom.iter().fold(Bloom::ZERO, |bloom, r| bloom | r.bloom_ref());

    compare_receipts_root_and_logs_bloom(
        receipts_root,
        logs_bloom,
        expected_receipts_root,
        expected_logs_bloom,
    )?;

    Ok(())
}
