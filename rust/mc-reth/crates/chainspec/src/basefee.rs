//! Base fee related utilities for Optimism chains.

use core::cmp::max;

use alloy_consensus::BlockHeader;
use alloy_eips::calc_next_block_base_fee;
use op_alloy_consensus::{EIP1559ParamError, decode_holocene_extra_data, decode_jovian_extra_data};
use reth_chainspec::BaseFeeParams;

/// Extracts the Holocene 1599 parameters from the encoded extra data from the parent header.
///
/// Caution: Caller must ensure that holocene is active in the parent header.
///
/// See also [Base fee computation](https://github.com/ethereum-optimism/specs/blob/main/specs/protocol/holocene/exec-engine.md#base-fee-computation)
pub fn decode_holocene_base_fee<H>(parent: &H) -> Result<u64, EIP1559ParamError>
where
    H: BlockHeader,
{
    // The denominator and elasticity are guaranteed non-zero: `decode_holocene_extra_data` rejects
    // a zero of either, per the Holocene header rules.
    let (elasticity, denominator) = decode_holocene_extra_data(parent.extra_data())?;
    let base_fee_params = BaseFeeParams::new(denominator as u128, elasticity as u128);

    Ok(parent.next_block_base_fee(base_fee_params).unwrap_or_default())
}

/// Extracts the Jovian 1599 parameters from the encoded extra data from the parent header.
/// Additionally to [`decode_holocene_base_fee`], checks if the next block base fee is less than the
/// minimum base fee, then the minimum base fee is returned.
///
/// Caution: Caller must ensure that jovian is active in the parent header.
///
/// See also [Base fee computation](https://github.com/ethereum-optimism/specs/blob/main/specs/protocol/jovian/exec-engine.md#base-fee-computation)
/// and [Minimum base fee in block header](https://github.com/ethereum-optimism/specs/blob/main/specs/protocol/jovian/exec-engine.md#minimum-base-fee-in-block-header)
pub fn compute_jovian_base_fee<H>(parent: &H) -> Result<u64, EIP1559ParamError>
where
    H: BlockHeader,
{
    // The denominator and elasticity are guaranteed non-zero: `decode_jovian_extra_data` rejects a
    // zero of either, per the Holocene header rules that Jovian inherits.
    let (elasticity, denominator, min_base_fee) = decode_jovian_extra_data(parent.extra_data())?;
    let base_fee_params = BaseFeeParams::new(denominator as u128, elasticity as u128);

    // Starting from Jovian, we use the maximum of the gas used and the blob gas used to calculate
    // the next base fee.
    let gas_used = max(parent.gas_used(), parent.blob_gas_used().unwrap_or_default());

    let next_base_fee = calc_next_block_base_fee(
        gas_used,
        parent.gas_limit(),
        parent.base_fee_per_gas().unwrap_or_default(),
        base_fee_params,
    );

    if next_base_fee < min_base_fee {
        return Ok(min_base_fee);
    }

    Ok(next_base_fee)
}
