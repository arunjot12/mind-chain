use super::{
    ExecutionInfo, build_post_exec_recovered_tx, is_invalid_derived_payload_error,
    try_include_post_exec_tx,
};
use crate::error::OpPayloadBuilderError;
use alloy_consensus::{Typed2718, transaction::Recovered};
use alloy_evm::{RecoveredTx, block::BlockExecutionError as AlloyBlockExecutionError};
use alloy_op_evm::block::OpBlockExecutionError;
use alloy_primitives::{Address, B256, U256};
use op_alloy_consensus::{PostExecPayloadValidationError, SDMGasEntry};
use reth_evm::execute::{BlockExecutionError, BlockValidationError};
use reth_mind_chain_primitives::OpTransactionSigned;
use reth_payload_builder_primitives::PayloadBuilderError;
use reth_storage_api::errors::ProviderError;
use revm::context::result::{EVMError, InvalidTransaction};
use std::cell::Cell;

fn entries(specs: &[(u64, u64)]) -> Vec<SDMGasEntry> {
    specs.iter().map(|&(index, gas_refund)| SDMGasEntry { index, gas_refund }).collect()
}

fn unwrap_post_exec(tx: Recovered<OpTransactionSigned>) -> (u8, u64, Vec<SDMGasEntry>) {
    let ty = tx.tx().ty();
    let OpTransactionSigned::PostExec(signed) = tx.into_inner() else {
        panic!("expected post-exec transaction");
    };
    (ty, signed.inner().payload.block_number, signed.inner().payload.gas_refund_entries.clone())
}

#[test]
fn execution_info_pre_refund_limit_uses_evm_gas_not_canonical_gas() {
    let mut info = ExecutionInfo::new();
    info.cumulative_gas_used = 50;
    info.cumulative_evm_gas_used = 90;

    assert!(
        !info.is_tx_over_limits(0, 100, None, None, 10, None, 0, None),
        "tx exactly filling the remaining pre-refund budget should fit"
    );
    assert!(
        info.is_tx_over_limits(0, 100, None, None, 11, None, 0, None),
        "tx that fits canonical gas but exceeds pre-refund gas must be skipped"
    );
}

#[test]
fn is_tx_over_limits_enforces_max_uncompressed_block_size() {
    let mut info = ExecutionInfo::new();
    info.cumulative_uncompressed_bytes = 100;

    // No limit configured: never over the uncompressed-size limit.
    assert!(!info.is_tx_over_limits(0, u64::MAX, None, None, 0, None, 1_000, None));

    // Exactly filling the remaining budget fits (100 + 50 == 150).
    assert!(!info.is_tx_over_limits(0, u64::MAX, None, None, 0, None, 50, Some(150)));

    // One byte over the limit is rejected (100 + 51 > 150).
    assert!(info.is_tx_over_limits(0, u64::MAX, None, None, 0, None, 51, Some(150)));
}

#[test]
fn build_post_exec_recovered_tx_wraps_entries_in_post_exec_tx() {
    let block_number = 42;
    let payload_entries = entries(&[(3, 17), (5, 23)]);
    let recovered =
        build_post_exec_recovered_tx::<OpTransactionSigned>(block_number, payload_entries.clone());

    assert_eq!(recovered.signer(), Address::ZERO);
    let (ty, decoded_block, decoded_entries) = unwrap_post_exec(recovered);
    assert_eq!(ty, op_alloy_consensus::POST_EXEC_TX_TYPE_ID);
    assert_eq!(decoded_block, block_number);
    assert_eq!(decoded_entries, payload_entries);
}

#[test]
fn try_include_post_exec_tx_skips_when_no_entries() {
    let called = Cell::new(false);
    let result = try_include_post_exec_tx::<OpTransactionSigned, _>(1, Vec::new(), |_tx| {
        called.set(true);
        Ok::<_, BlockExecutionError>(0)
    });
    assert!(matches!(result, Ok(false)));
    assert!(!called.get(), "execute must not run when there are no entries");
}

#[test]
fn try_include_post_exec_tx_executes_post_exec_tx_on_happy_path() {
    let block_number = 99;
    let payload_entries = entries(&[(0, 7)]);
    let captured: Cell<Option<(u8, u64, Vec<SDMGasEntry>)>> = Cell::new(None);

    let result = try_include_post_exec_tx::<OpTransactionSigned, _>(
        block_number,
        payload_entries.clone(),
        |tx| {
            captured.set(Some(unwrap_post_exec(tx)));
            Ok::<_, BlockExecutionError>(21_000)
        },
    );

    assert!(matches!(result, Ok(true)));
    let (ty, decoded_block, decoded_entries) = captured.take().expect("execute closure ran");
    assert_eq!(ty, op_alloy_consensus::POST_EXEC_TX_TYPE_ID);
    assert_eq!(decoded_block, block_number);
    assert_eq!(decoded_entries, payload_entries);
}

/// Consensus-critical: if the post-exec tx fails to execute, the payload build MUST abort.
/// Returning `Ok(_)` (e.g. an empty block, or silently dropping the tx) would diverge the
/// producer from any honest verifier, because the verifier observes refunds from the normal
/// txs and expects a matching post-exec tx.
#[test]
fn try_include_post_exec_tx_aborts_when_execution_fails() {
    let called = Cell::new(false);
    let result = try_include_post_exec_tx::<OpTransactionSigned, _>(1, entries(&[(0, 7)]), |_tx| {
        called.set(true);
        Err::<u64, _>(BlockExecutionError::msg("forced post-exec tx failure"))
    });

    assert!(called.get(), "execute must be invoked so its error can propagate");
    match result {
        Err(PayloadBuilderError::EvmExecutionError(err)) => {
            assert!(err.to_string().contains("forced post-exec tx failure"));
        }
        other => panic!("expected EvmExecutionError, got {other:?}"),
    }
}

/// Only errors caused by the payload's transactions may be classified as invalid: a state-access
/// error reported as invalid would make the rollup node derive a wrong deposit-only block.
#[test]
fn derived_payload_error_classification() {
    let validation =
        |err: BlockValidationError| PayloadBuilderError::from(BlockExecutionError::from(err));
    let op_execution =
        |err: OpBlockExecutionError| validation(BlockValidationError::Other(Box::new(err)));
    let op_payload = PayloadBuilderError::other::<OpPayloadBuilderError>;
    let invalid_post_exec = || OpBlockExecutionError::InvalidPostExecPayload("x".into());
    let cases = [
        (
            validation(BlockValidationError::InvalidTx {
                hash: B256::ZERO,
                error: Box::new(InvalidTransaction::NonceTooHigh { tx: 5, state: 0 }),
            }),
            true,
        ),
        (validation(BlockValidationError::BlockGasExceeded), true),
        (
            validation(BlockValidationError::TransactionGasLimitMoreThanAvailableBlockGas {
                transaction_gas_limit: 2,
                block_available_gas: 1,
            }),
            true,
        ),
        (op_execution(invalid_post_exec()), true),
        (
            op_execution(OpBlockExecutionError::TransactionDaFootprintAboveGasLimit {
                transaction_da_footprint: 2,
                available_block_da_footprint: 1,
            }),
            true,
        ),
        (op_execution(OpBlockExecutionError::UnexpectedNonDepositTxInForkActivationBlock), true),
        (
            op_execution(OpBlockExecutionError::PostExecSettlementUnderflow {
                address: Address::ZERO,
                delta: U256::from(1),
            }),
            true,
        ),
        (PayloadBuilderError::other(alloy_rlp::Error::UnexpectedLength), true),
        (
            PayloadBuilderError::other(PostExecPayloadValidationError::UnexpectedPostExecTx {
                tx_index: 0,
            }),
            true,
        ),
        (op_payload(OpPayloadBuilderError::TransactionEcRecoverFailed), true),
        (op_payload(OpPayloadBuilderError::BlobTransactionRejected), true),
        (op_payload(OpPayloadBuilderError::L1BlockInfoParseFailed), false),
        (op_payload(OpPayloadBuilderError::AccountLoadFailed(Address::ZERO)), false),
        (op_payload(OpPayloadBuilderError::ForceCreate2DeployerFail), false),
        (PayloadBuilderError::invalid_attributes(alloy_rlp::Error::UnexpectedLength), true),
        (PayloadBuilderError::MissingParentHeader(B256::ZERO), false),
        // Pins only the classifier. The node test's "empty refund entries" case, rejected during
        // pre-execution, is what catches `OpBuilder::build` reporting that error as `Internal`.
        (
            PayloadBuilderError::Internal(
                BlockExecutionError::from(BlockValidationError::Other(Box::new(
                    invalid_post_exec(),
                )))
                .into(),
            ),
            false,
        ),
        // The parent's state is unavailable.
        (ProviderError::StateForHashNotFound(B256::ZERO).into(), false),
        (BlockExecutionError::other(ProviderError::StateForHashNotFound(B256::ZERO)).into(), false),
        // A state read fails while a transaction executes.
        (
            PayloadBuilderError::from(AlloyBlockExecutionError::evm(
                EVMError::<ProviderError, InvalidTransaction>::Database(
                    ProviderError::StateForHashNotFound(B256::ZERO),
                ),
                B256::ZERO,
            )),
            false,
        ),
    ];

    for (err, expected) in cases {
        assert_eq!(is_invalid_derived_payload_error(&err), expected, "{err:?}");
    }
}
