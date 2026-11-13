#![doc = include_str!("../README.md")]
#![cfg_attr(docsrs, feature(doc_cfg))]

pub mod cli;
pub mod cli_metrics;
pub mod commands;
pub mod flags;

pub(crate) mod version;

fn main() {
    use clap::Parser;

    kona_cli::sigsegv_handler::install();
    kona_cli::backtrace::enable();

    if let Err(err) = cli::Cli::parse().run() {
        eprintln!("Error: {err:?}");
        std::process::exit(1);
    }
}
