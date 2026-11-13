#![doc = include_str!("../README.md")]
#![cfg_attr(docsrs, feature(doc_cfg))]

mod cli;
mod version;

fn main() {
    use clap::Parser;

    cli::Cli::parse().run();
}
