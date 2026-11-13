#![doc = include_str!("../README.md")]
#![cfg_attr(docsrs, feature(doc_cfg))]
#![cfg_attr(not(feature = "arbitrary"), no_std)]

extern crate alloc;

mod pre_state;
pub use pre_state::{
    INVALID_TRANSITION, INVALID_TRANSITION_HASH, OptimisticBlock, PreState,
    TRANSITION_STATE_MAX_STEPS, TransitionState,
};

mod hint;
pub use hint::HintType;

mod provider;
pub use provider::{
    ChainScopedHinter, LogCommitmentSource, OracleInteropProvider, receipts_from_trie,
};

pub mod boot;
pub use boot::BootInfo;

mod consolidation;
pub use consolidation::{ConsolidationError, SuperchainConsolidator};
