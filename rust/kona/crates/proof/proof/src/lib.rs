#![doc = include_str!("../README.md")]
#![cfg_attr(docsrs, feature(doc_cfg))]
#![no_std]

extern crate alloc;

#[macro_use]
extern crate tracing;

pub mod l1;

pub mod l2;

pub mod sync;

pub mod errors;

pub mod executor;

mod hint;
pub use hint::{Hint, HintType};

pub mod boot;
pub use boot::BootInfo;

mod caching_oracle;
pub use caching_oracle::{CachingOracle, FlushableCache};

mod blocking_runtime;
pub use blocking_runtime::block_on;

mod eip2935;
pub use eip2935::eip_2935_history_lookup;
