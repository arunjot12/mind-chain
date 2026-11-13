#![doc = include_str!("../README.md")]
#![cfg_attr(docsrs, feature(doc_cfg))]

#[macro_use]
extern crate tracing;

mod admin;
pub use admin::{AdminRpc, NetworkAdminQuery};

mod client;
pub use client::{EngineRpcClient, SequencerAdminAPIClient, SequencerAdminAPIError};

mod config;
pub use config::RpcBuilder;

mod net;
pub use net::P2pRpc;

mod p2p;

mod response;
pub use response::SafeHeadResponse;

mod output;
pub use output::OutputResponse;

mod dev;
pub use dev::DevEngineRpc;

mod jsonrpsee;
pub use jsonrpsee::{
    AdminApiServer, DevEngineApiServer, HealthzApiServer, MinerApiExtServer, OpAdminApiServer,
    OpP2PApiServer, RollupNodeApiServer, WsServer,
};

#[cfg(feature = "client")]
pub use jsonrpsee::RollupNodeApiClient;

mod rollup;
pub use rollup::RollupRpc;

mod l1_watcher;
pub use l1_watcher::{L1State, L1WatcherQueries, L1WatcherQuerySender};

mod ws;
pub use ws::WsRPC;

mod health;
pub use health::{HealthzResponse, HealthzRpc};
