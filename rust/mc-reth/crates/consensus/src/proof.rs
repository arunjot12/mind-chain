//! Helper function for Receipt root calculation for Optimism hardforks.

use alloc::vec::Vec;
use alloy_consensus::ReceiptWithBloom;
use alloy_eips::eip2718::Encodable2718;
use alloy_primitives::B256;
use alloy_trie::root::ordered_trie_root_with_encoder;
use reth_mind_chain_forks::OpHardforks;
use reth_mind_chain_primitives::DepositReceipt;

/// Calculates the receipt root for a header.
pub fn calculate_receipt_root_optimism<R: DepositReceipt>(
    receipts: &[ReceiptWithBloom<&R>],
    chain_spec: impl OpHardforks,
    timestamp: u64,
) -> B256 {
    // There is a minor bug in op-geth and op-erigon where in the Regolith hardfork,
    // the receipt root calculation does not include the deposit nonce in the receipt
    // encoding. In the Regolith Hardfork, we must strip the deposit nonce from the
    // receipts before calculating the receipt root. This was corrected in the Canyon
    // hardfork.
    if chain_spec.is_regolith_active_at_timestamp(timestamp) &&
        !chain_spec.is_canyon_active_at_timestamp(timestamp)
    {
        let receipts = receipts
            .iter()
            .map(|receipt| {
                let mut receipt = receipt.clone().map_receipt(|r| r.clone());
                if let Some(receipt) = receipt.receipt.as_deposit_receipt_mut() {
                    receipt.deposit_nonce = None;
                }
                receipt
            })
            .collect::<Vec<_>>();

        return ordered_trie_root_with_encoder(receipts.as_slice(), |r, buf| r.encode_2718(buf));
    }

    ordered_trie_root_with_encoder(receipts, |r, buf| r.encode_2718(buf))
}
