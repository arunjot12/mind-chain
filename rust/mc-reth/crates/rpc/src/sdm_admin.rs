//! Local SDM `PostExec` opt-in admin RPC for mc-reth's payload builder.
//!
//! The protocol gate (hardfork activation) is a consensus rule shared by every node. This module
//! adds an orthogonal *operator* gate: even on an SDM-active chain, mc-reth's standard payload
//! builder produces `PostExec` txs only when the operator has opted in via
//! `admin_setOperatorSdmOptIn`. Both must be true in order for SDM to be active.
//!
//! State is in-memory and starts disabled on every process boot; persistence is deliberately
//! out of scope.

use alloy_hardforks::ForkCondition;
use jsonrpsee::{core::RpcResult, proc_macros::rpc};
use reth_mind_chain_forks::{OpHardfork, OpHardforks};
use reth_mind_chain_payload_builder::config::OperatorSdmOptIn;
use serde::{Deserialize, Serialize};
use std::{
    sync::Arc,
    time::{SystemTime, UNIX_EPOCH},
};

/// Status snapshot returned by `admin_sdmStatus`. Shape matches mc-node so a single client view
/// can render either participant's state.
#[derive(Debug, Clone, Copy, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct SdmStatus {
    /// Whether the operator has opted in via `admin_setOperatorSdmOptIn`.
    pub operator_sdm_opt_in: bool,
    /// Whether SDM is active per the chain spec at `query_timestamp`.
    pub protocol_active: bool,
    /// AND of the two gates.
    pub effective: bool,
    /// Activation timestamp of the protocol gate if scheduled.
    pub activation_time: Option<u64>,
}

#[cfg_attr(not(test), rpc(server, namespace = "admin"))]
#[cfg_attr(test, rpc(server, client, namespace = "admin"))]
pub trait SdmAdminApi {
    /// Toggle local `PostExec` production. Starts disabled on process boot.
    #[method(name = "setOperatorSdmOptIn")]
    fn set_operator_sdm_opt_in(&self, enabled: bool) -> RpcResult<()>;

    /// Report the local opt-in flag, the chain-spec gate at `query_timestamp`, and the AND.
    /// If `query_timestamp` is omitted, uses the current wall-clock time.
    #[method(name = "sdmStatus")]
    fn sdm_status(&self, query_timestamp: Option<u64>) -> RpcResult<SdmStatus>;
}

/// Concrete `admin_` SDM RPC server holding the shared opt-in flag and the chain spec used
/// to evaluate the protocol gate for status queries.
#[derive(Debug, Clone)]
pub struct OpSdmAdminApi<ChainSpec> {
    opt_in: OperatorSdmOptIn,
    chain_spec: Arc<ChainSpec>,
}

impl<ChainSpec> OpSdmAdminApi<ChainSpec> {
    /// Construct a handler that mutates `opt_in` and reads the protocol gate from `chain_spec`.
    /// The shared `OperatorSdmOptIn` should also be threaded into the payload builder's
    /// `OpBuilderConfig` so writes here take effect on the next produced block.
    pub const fn new(opt_in: OperatorSdmOptIn, chain_spec: Arc<ChainSpec>) -> Self {
        Self { opt_in, chain_spec }
    }
}

impl<ChainSpec> SdmAdminApiServer for OpSdmAdminApi<ChainSpec>
where
    ChainSpec: OpHardforks + Send + Sync + 'static,
{
    fn set_operator_sdm_opt_in(&self, enabled: bool) -> RpcResult<()> {
        self.opt_in.set(enabled);
        Ok(())
    }

    fn sdm_status(&self, query_timestamp: Option<u64>) -> RpcResult<SdmStatus> {
        let timestamp = query_timestamp.unwrap_or_else(current_unix_timestamp);
        let operator_sdm_opt_in = self.opt_in.enabled();
        let protocol_active =
            reth_mind_chain_evm::is_sdm_active_at_timestamp(&self.chain_spec, timestamp);
        let activation_time = match self.chain_spec.op_fork_activation(OpHardfork::Lagoon) {
            ForkCondition::Timestamp(timestamp) => Some(timestamp),
            _ => None,
        };
        let effective = operator_sdm_opt_in && protocol_active;

        Ok(SdmStatus { operator_sdm_opt_in, protocol_active, effective, activation_time })
    }
}

fn current_unix_timestamp() -> u64 {
    SystemTime::now().duration_since(UNIX_EPOCH).map_or(0, |duration| duration.as_secs())
}

#[cfg(test)]
mod tests {
    use super::*;

    // Any timestamp works for the gate checks below: SDM's gating hardfork (Lagoon) is activated
    // at genesis on the active spec, and never on the inactive one.
    /// The status serializes with camelCase JSON keys.
    #[test]
    fn sdm_status_serializes_with_camel_case_keys() {
        let json = serde_json::to_value(SdmStatus {
            operator_sdm_opt_in: true,
            protocol_active: false,
            effective: false,
            activation_time: Some(123),
        })
        .unwrap();

        assert_eq!(json["operatorSdmOptIn"], serde_json::json!(true));
        assert_eq!(json["protocolActive"], serde_json::json!(false));
        assert_eq!(json["effective"], serde_json::json!(false));
        assert_eq!(json["activationTime"], serde_json::json!(123));
    }
}
