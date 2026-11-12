use alloy_consensus::BlockHeader;
use alloy_primitives::B256;
use alloy_rpc_types_engine::{ExecutionPayloadEnvelopeV2, ExecutionPayloadV1};
use op_alloy_rpc_types_engine::{
    OpExecutionData, OpExecutionPayloadEnvelope, OpExecutionPayloadEnvelopeV3,
    OpExecutionPayloadEnvelopeV4,
};
use reth_consensus::ConsensusError;
use reth_mind_chain_consensus::isthmus;
use reth_mind_chain_forks::OpHardforks;
use reth_mind_chain_payload_builder::{
    OpExecData, OpExecutionPayloadValidator, OpPayloadAttrs, OpPayloadTypes,
};
use reth_mind_chain_primitives::{L2_TO_L1_MESSAGE_PASSER_ADDRESS, OpBlock};
use reth_node_api::{
    BuiltPayload, EngineApiValidator, EngineTypes, InsertBlockErrorKind, NodePrimitives,
    PayloadValidator,
    payload::{
        EngineApiMessageVersion, EngineObjectValidationError, MessageValidationKind,
        NewPayloadError, PayloadOrAttributes, PayloadTypes, VersionSpecificValidationError,
        validate_parent_beacon_block_root_presence,
    },
    validate_version_specific_fields,
};
use reth_primitives_traits::{Block, RecoveredBlock, SealedBlock, SealedHeader, SignedTransaction};
use reth_provider::{ProviderResult, StateProviderBox, StateProviderFactory};
use reth_trie_common::{HashedPostState, KeyHasher};
use std::{marker::PhantomData, sync::Arc};

/// The types used in the optimism beacon consensus engine.
#[derive(Debug, Default, Clone, serde::Deserialize, serde::Serialize)]
#[non_exhaustive]
pub struct OpEngineTypes<T: PayloadTypes = OpPayloadTypes> {
    _marker: PhantomData<T>,
}

impl<T: PayloadTypes<ExecutionData = OpExecData>> PayloadTypes for OpEngineTypes<T>
where
    OpExecData: From<<T as PayloadTypes>::BuiltPayload>,
{
    type ExecutionData = T::ExecutionData;
    type BuiltPayload = T::BuiltPayload;
    type PayloadAttributes = T::PayloadAttributes;

    fn block_to_payload(
        block: SealedBlock<
            <<Self::BuiltPayload as BuiltPayload>::Primitives as NodePrimitives>::Block,
        >,
        _bal: Option<alloy_primitives::Bytes>,
    ) -> <T as PayloadTypes>::ExecutionData {
        OpExecData(OpExecutionData::from(
            OpExecutionPayloadEnvelope::from_block_unchecked(
                block.hash(),
                &block.into_block().into_ethereum_block(),
            )
            .expect("built OP blocks must normalize"),
        ))
    }
}

impl<T: PayloadTypes<ExecutionData = OpExecData>> EngineTypes for OpEngineTypes<T>
where
    OpExecData: From<<T as PayloadTypes>::BuiltPayload>,
    T::BuiltPayload: BuiltPayload<Primitives: NodePrimitives<Block = OpBlock>>
        + TryInto<ExecutionPayloadV1>
        + TryInto<ExecutionPayloadEnvelopeV2>
        + TryInto<OpExecutionPayloadEnvelopeV3>
        + TryInto<OpExecutionPayloadEnvelopeV4>,
{
    type ExecutionPayloadEnvelopeV1 = ExecutionPayloadV1;
    type ExecutionPayloadEnvelopeV2 = ExecutionPayloadEnvelopeV2;
    type ExecutionPayloadEnvelopeV3 = OpExecutionPayloadEnvelopeV3;
    type ExecutionPayloadEnvelopeV4 = OpExecutionPayloadEnvelopeV4;
    type ExecutionPayloadEnvelopeV5 = OpExecutionPayloadEnvelopeV4;
    type ExecutionPayloadEnvelopeV6 = OpExecutionPayloadEnvelopeV4;
}

/// Validator for Optimism engine API.
#[derive(Debug)]
pub struct OpEngineValidator<P, Tx, ChainSpec> {
    inner: OpExecutionPayloadValidator<ChainSpec>,
    provider: P,
    hashed_addr_l2tol1_msg_passer: B256,
    phantom: PhantomData<Tx>,
}

impl<P, Tx, ChainSpec> OpEngineValidator<P, Tx, ChainSpec> {
    /// Instantiates a new validator.
    pub fn new<KH: KeyHasher>(chain_spec: Arc<ChainSpec>, provider: P) -> Self {
        let hashed_addr_l2tol1_msg_passer = KH::hash_key(L2_TO_L1_MESSAGE_PASSER_ADDRESS);
        Self {
            inner: OpExecutionPayloadValidator::new(chain_spec),
            provider,
            hashed_addr_l2tol1_msg_passer,
            phantom: PhantomData,
        }
    }
}

impl<P, Tx, ChainSpec> Clone for OpEngineValidator<P, Tx, ChainSpec>
where
    P: Clone,
    ChainSpec: OpHardforks,
{
    fn clone(&self) -> Self {
        Self {
            inner: OpExecutionPayloadValidator::new(self.inner.clone()),
            provider: self.provider.clone(),
            hashed_addr_l2tol1_msg_passer: self.hashed_addr_l2tol1_msg_passer,
            phantom: Default::default(),
        }
    }
}

impl<P, Tx, ChainSpec> OpEngineValidator<P, Tx, ChainSpec>
where
    ChainSpec: OpHardforks,
{
    /// Returns the chain spec used by the validator.
    #[inline]
    pub fn chain_spec(&self) -> &ChainSpec {
        self.inner.chain_spec()
    }
}

impl<P, Tx, ChainSpec, Types> PayloadValidator<Types> for OpEngineValidator<P, Tx, ChainSpec>
where
    P: StateProviderFactory + Unpin + 'static,
    Tx: SignedTransaction + Unpin + 'static,
    ChainSpec: OpHardforks + Send + Sync + 'static,
    Types: PayloadTypes<ExecutionData = OpExecData>,
{
    type Block = alloy_consensus::Block<Tx>;

    fn validate_block_post_execution_with_hashed_state<'a>(
        &self,
        state_updates: impl FnOnce() -> &'a HashedPostState,
        block: &RecoveredBlock<Self::Block>,
        _parent_header: &SealedHeader<<Self::Block as Block>::Header>,
        parent_state: impl FnOnce() -> ProviderResult<StateProviderBox>,
    ) -> Result<(), InsertBlockErrorKind> {
        if self.chain_spec().is_isthmus_active_at_timestamp(block.timestamp()) {
            let state = parent_state()?;
            let predeploy_storage_updates = state_updates()
                .storages
                .get(&self.hashed_addr_l2tol1_msg_passer)
                .cloned()
                .unwrap_or_default();
            isthmus::verify_withdrawals_root_prehashed(
                predeploy_storage_updates,
                state,
                block.header(),
            )
            .map_err(|err| {
                ConsensusError::msg(format!("failed to verify block post-execution: {err}"))
            })?
        }

        Ok(())
    }

    fn convert_payload_to_block(
        &self,
        payload: OpExecData,
    ) -> Result<SealedBlock<Self::Block>, NewPayloadError> {
        self.inner.ensure_well_formed_payload(payload.0).map_err(NewPayloadError::other)
    }
}

impl<Types, P, Tx, ChainSpec> EngineApiValidator<Types> for OpEngineValidator<P, Tx, ChainSpec>
where
    Types: PayloadTypes<
            PayloadAttributes = OpPayloadAttrs,
            ExecutionData = OpExecData,
            BuiltPayload: BuiltPayload<Primitives: NodePrimitives<SignedTx = Tx>>,
        >,
    P: StateProviderFactory + Unpin + 'static,
    Tx: SignedTransaction + Unpin + 'static,
    ChainSpec: OpHardforks + Send + Sync + 'static,
{
    fn validate_version_specific_fields(
        &self,
        version: EngineApiMessageVersion,
        payload_or_attrs: PayloadOrAttributes<
            '_,
            Types::ExecutionData,
            <Types as PayloadTypes>::PayloadAttributes,
        >,
    ) -> Result<(), EngineObjectValidationError> {
        validate_withdrawals_presence(
            self.chain_spec(),
            version,
            payload_or_attrs.message_validation_kind(),
            payload_or_attrs.timestamp(),
            payload_or_attrs.withdrawals().is_some(),
        )?;
        validate_parent_beacon_block_root_presence(
            self.chain_spec(),
            version,
            payload_or_attrs.message_validation_kind(),
            payload_or_attrs.timestamp(),
            payload_or_attrs.parent_beacon_block_root().is_some(),
        )
    }

    fn ensure_well_formed_attributes(
        &self,
        version: EngineApiMessageVersion,
        attributes: &<Types as PayloadTypes>::PayloadAttributes,
    ) -> Result<(), EngineObjectValidationError> {
        ensure_well_formed_payload_attributes(self.chain_spec(), version, attributes)
    }
}

/// Checks that payload attributes sent with an `engine_forkchoiceUpdated` call of message version
/// `version` are well formed for the forks active at their timestamp:
///
/// - the engine API's version-specific fields (withdrawals, parent beacon block root, slot number)
///   are present exactly when `version` and the active forks require them, and the timestamp falls
///   in a fork `version` supports;
/// - a gas limit is set;
/// - from Holocene, EIP-1559 parameters are set, with elasticity and denominator either both zero
///   or both non-zero;
/// - from Jovian, a minimum base fee is set; before Jovian, none is.
///
/// These are the checks [`OpEngineValidator`] runs on the attributes of an
/// `engine_forkchoiceUpdated` call, taking `version` from the method called. They read no state.
pub fn ensure_well_formed_payload_attributes<ChainSpec: OpHardforks>(
    chain_spec: &ChainSpec,
    version: EngineApiMessageVersion,
    attributes: &OpPayloadAttrs,
) -> Result<(), EngineObjectValidationError> {
    validate_version_specific_fields(
        chain_spec,
        version,
        PayloadOrAttributes::<OpExecData, OpPayloadAttrs>::PayloadAttributes(attributes),
    )?;

    if attributes.gas_limit.is_none() {
        return Err(EngineObjectValidationError::InvalidParams(
            "MissingGasLimitInPayloadAttributes".to_string().into(),
        ));
    }

    if chain_spec.is_holocene_active_at_timestamp(attributes.payload_attributes.timestamp) {
        let (elasticity, denominator) = attributes.decode_eip_1559_params().ok_or_else(|| {
            EngineObjectValidationError::InvalidParams(
                "MissingEip1559ParamsInPayloadAttributes".to_string().into(),
            )
        })?;

        if elasticity != 0 && denominator == 0 {
            return Err(EngineObjectValidationError::InvalidParams(
                "Eip1559ParamsDenominatorZero".to_string().into(),
            ));
        } else if denominator != 0 && elasticity == 0 {
            return Err(EngineObjectValidationError::InvalidParams(
                "Eip1559ParamsElasticityZero".to_string().into(),
            ));
        }
    }

    if chain_spec.is_jovian_active_at_timestamp(attributes.payload_attributes.timestamp) {
        if attributes.min_base_fee.is_none() {
            return Err(EngineObjectValidationError::InvalidParams(
                "MissingMinBaseFeeInPayloadAttributes".to_string().into(),
            ));
        }
    } else if attributes.min_base_fee.is_some() {
        return Err(EngineObjectValidationError::InvalidParams(
            "MinBaseFeeNotAllowedBeforeJovian".to_string().into(),
        ));
    }

    Ok(())
}

/// Validates the presence of the `withdrawals` field according to the payload timestamp.
///
/// After Canyon, withdrawals field must be [Some].
/// Before Canyon, withdrawals field must be [None];
///
/// Canyon activates the Shanghai EIPs, see the Canyon specs for more details:
/// <https://github.com/mindchain/mind-chain/blob/ab926c5fd1e55b5c864341c44842d6d1ca679d99/specs/superchain-upgrades.md#canyon>
pub fn validate_withdrawals_presence(
    chain_spec: impl OpHardforks,
    version: EngineApiMessageVersion,
    message_validation_kind: MessageValidationKind,
    timestamp: u64,
    has_withdrawals: bool,
) -> Result<(), EngineObjectValidationError> {
    let is_shanghai = chain_spec.is_canyon_active_at_timestamp(timestamp);

    match version {
        EngineApiMessageVersion::V1 => {
            if has_withdrawals {
                return Err(message_validation_kind
                    .to_error(VersionSpecificValidationError::WithdrawalsNotSupportedInV1));
            }
            if is_shanghai {
                return Err(message_validation_kind
                    .to_error(VersionSpecificValidationError::NoWithdrawalsPostShanghai));
            }
        }
        EngineApiMessageVersion::V2 |
        EngineApiMessageVersion::V3 |
        EngineApiMessageVersion::V4 |
        EngineApiMessageVersion::V5 |
        EngineApiMessageVersion::V6 => {
            if is_shanghai && !has_withdrawals {
                return Err(message_validation_kind
                    .to_error(VersionSpecificValidationError::NoWithdrawalsPostShanghai));
            }
            if !is_shanghai && has_withdrawals {
                return Err(message_validation_kind
                    .to_error(VersionSpecificValidationError::HasWithdrawalsPreShanghai));
            }
        }
    };

    Ok(())
}

#[cfg(test)]
mod test {}
