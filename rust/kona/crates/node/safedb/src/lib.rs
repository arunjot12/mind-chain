#![doc = include_str!("../README.md")]
#![cfg_attr(docsrs, feature(doc_cfg))]

mod disabled;
mod error;
mod traits;

#[cfg(feature = "rocksdb")]
mod encoding;
#[cfg(feature = "rocksdb")]
mod safe_db;

#[cfg(test)]
mod tests;

pub use disabled::DisabledDatabase;
pub use error::SafeDbError;
pub use traits::{SafeDb, SafeHeadRecord};

#[cfg(feature = "rocksdb")]
pub use safe_db::SafeDatabase;
