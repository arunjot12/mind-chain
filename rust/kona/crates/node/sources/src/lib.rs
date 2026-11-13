#![doc = include_str!("../README.md")]
#![cfg_attr(not(test), warn(unused_crate_dependencies))]
#![cfg_attr(docsrs, feature(doc_cfg))]

mod signer;
pub use signer::{
    BlockSigner, BlockSignerError, BlockSignerHandler, BlockSignerStartError, RemoteSigner,
    RemoteSignerError, RemoteSignerHandler, RemoteSignerStartError,
};

mod tls;
pub use tls::{
    CertificateError, ClientCert, ReloadingRpcClient, ReloadingRpcClientError, TlsPaths,
    redacted_url,
};
