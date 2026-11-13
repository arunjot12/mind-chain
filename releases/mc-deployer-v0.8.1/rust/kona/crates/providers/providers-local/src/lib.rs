#![doc = include_str!("../README.md")]
#![cfg_attr(docsrs, feature(doc_cfg))]

mod buffer;
pub use buffer::{CacheStats, CachedBlock, ChainBufferError, ChainStateBuffer, ChainStateEvent};

mod buffered;
pub use buffered::{BufferedL2Provider, BufferedProviderError};

#[cfg(feature = "metrics")]
mod metrics;
#[cfg(feature = "metrics")]
pub use metrics::Metrics;
