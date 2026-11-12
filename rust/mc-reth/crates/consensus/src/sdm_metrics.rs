//! SDM refund observation counters for blocks this node executed and accepted, recorded from
//! post-execution validation. They follow imported blocks rather than locally built ones: blocks
//! later reorged out are included, and a block re-executed after an unwind is counted again. The
//! payload builder's `sdm_refund_gas_total` measures the same quantity for blocks this node
//! produced, a different population, so the two must not be summed.
//!
//! Handles are resolved on every record rather than cached, so a record before the recorder is
//! installed drops that one block instead of silencing the counters for the process.

use op_alloy_consensus::{OpTransaction, PostExecPayload};
use reth_metrics::{
    Metrics,
    metrics::{self, Counter},
};

const NO_LABELS: &[(&str, &str)] = &[];

/// One struct for all three counters on purpose: a block without refunds is a real observation, so
/// a node on a chain without SDM publishing zeros for the refund series is correct, and the alert
/// ratio needs `block_gas_total` from every node either way.
#[derive(Metrics)]
#[metrics(scope = "op_sdm")]
struct SdmImportMetrics {
    /// Canonical (post-refund) gas used by validated blocks.
    block_gas_total: Counter,
    /// Warming-refund entries carried by validated blocks' post-exec transactions.
    refund_entries_total: Counter,
    /// Warming-refund gas carried by validated blocks' post-exec transactions.
    refund_gas_total: Counter,
}

/// Records one validated block: its canonical gas always, and its refunds when it carries a
/// post-exec transaction.
pub(crate) fn record_validated_block<T: OpTransaction>(block_gas_used: u64, transactions: &[T]) {
    let metrics = SdmImportMetrics::new_with_labels(NO_LABELS);
    metrics.block_gas_total.increment(block_gas_used);
    if let Some(payload) = post_exec_payload(transactions) {
        metrics.refund_entries_total.increment(payload.gas_refund_entries.len() as u64);
        metrics.refund_gas_total.increment(payload.total_gas_refund());
    }
}

/// The post-exec payload of a block that passed execution. The executor has already enforced that
/// a post-exec transaction is unique and final, so only the last transaction is inspected.
fn post_exec_payload<T: OpTransaction>(transactions: &[T]) -> Option<&PostExecPayload> {
    transactions.last().and_then(OpTransaction::as_post_exec).map(|tx| &tx.inner().payload)
}

#[cfg(test)]
mod tests {
    use super::*;

    use alloy_consensus::Sealable;

    use op_alloy_consensus::{SDMGasEntry, TxDeposit, build_post_exec_tx};

    use reth_mind_chain_primitives::OpTransactionSigned;

    use rstest::rstest;

    fn deposit_tx() -> OpTransactionSigned {
        OpTransactionSigned::Deposit(TxDeposit::default().seal_slow())
    }

    fn post_exec_tx(block_number: u64, refunds: &[u64]) -> OpTransactionSigned {
        let entries = refunds
            .iter()
            .enumerate()
            .map(|(index, gas_refund)| SDMGasEntry { index: index as u64, gas_refund: *gas_refund })
            .collect();
        OpTransactionSigned::PostExec(build_post_exec_tx(block_number, entries).seal_slow())
    }

    #[rstest]
    #[case::trailing(vec![deposit_tx(), post_exec_tx(1, &[3])], Some(3))]
    #[case::absent(vec![deposit_tx()], None)]
    #[case::not_last(vec![post_exec_tx(1, &[3]), deposit_tx()], None)]
    fn post_exec_payload_reads_the_final_transaction(
        #[case] transactions: Vec<OpTransactionSigned>,
        #[case] expected_refund: Option<u64>,
    ) {
        let refund = post_exec_payload(&transactions).map(PostExecPayload::total_gas_refund);

        assert_eq!(refund, expected_refund);
    }
}
