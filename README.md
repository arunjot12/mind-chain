<div align="center">
  <br />
  <h1>Mind Chain</h1>
  <h3>The Mind Chain rollup stack: node, batcher, proposer, challenger, execution client and contracts.</h3>
  <br />
</div>

**Table of Contents**

<!--TOC-->

- [What is Mind Chain?](#what-is-mind-chain)
- [Getting Started](#getting-started)
- [Directory Structure](#directory-structure)
- [Releases](#releases)
- [Downloading & Shallow-Cloning](#downloading--shallow-cloning)
- [License](#license)

<!--TOC-->

## What is Mind Chain?

Mind Chain is an Ethereum Layer 2 rollup. This repository holds every Mind Chain service in one place:
the rollup node (`mc-node`), the batch submitter (`mc-batcher`), the output proposer (`mc-proposer`),
the dispute game agent (`mc-challenger`), the deployer (`mc-deployer`), the execution client (`mc-reth`)
and the L1/L2 smart contracts.

See [MINDCHAIN.md](./MINDCHAIN.md) for how to clone, build, run and release the services.

## Getting Started

[CONTRIBUTING.md](./CONTRIBUTING.md) explains how to set up a development environment and contribute.
Use the [Developer Quick Start](./CONTRIBUTING.md#development-quick-start) to install the toolchain.

## Directory Structure

<pre>
├── <a href="./cannon">cannon</a>: Onchain MIPS instruction emulator for fault proofs
├── <a href="./mindchain">mindchain</a>: Mind Chain version pins and helper scripts
├── <a href="./mc-acceptance-tests">mc-acceptance-tests</a>: Acceptance tests
├── <a href="./mc-alt-da">mc-alt-da</a>: Alternative Data Availability mode (beta)
├── <a href="./mc-batcher">mc-batcher</a>: L2-Batch Submitter, submits bundles of batches to L1
├── <a href="./mc-chain-ops">mc-chain-ops</a>: State surgery utilities
├── <a href="./mc-challenger">mc-challenger</a>: Dispute game challenge agent
├── <a href="./mc-conductor">mc-conductor</a>: High-availability sequencer service
├── <a href="./mc-deployer">mc-deployer</a>: CLI tool for deploying and upgrading the Mind Chain smart contracts
├── <a href="./mc-devstack">mc-devstack</a>: Flexible test frontend for integration and acceptance testing
├── <a href="./mc-dispute-mon">mc-dispute-mon</a>: Off-chain service to monitor dispute games
├── <a href="./mc-dripper">mc-dripper</a>: Controlled token distribution service
├── <a href="./mc-e2e">mc-e2e</a>: End-to-End testing of all components in Go
├── <a href="./mc-faucet">mc-faucet</a>: Dev-faucet with support for multiple chains
├── <a href="./mc-fetcher">mc-fetcher</a>: Data fetching utilities
├── <a href="./mc-interop-mon">mc-interop-mon</a>: Interoperability monitoring service
├── <a href="./mc-node">mc-node</a>: Rollup consensus-layer client
├── <a href="./mc-preimage">mc-preimage</a>: Go bindings for Preimage Oracle
├── <a href="./mc-proposer">mc-proposer</a>: L2-Output Submitter, submits proposals to L1
├── <a href="./mc-service">mc-service</a>: Common codebase utilities
├── <a href="./mc-sync-tester">mc-sync-tester</a>: Sync testing utilities
├── <a href="./mc-test-sequencer">mc-test-sequencer</a>: Test sequencer for development
├── <a href="./mc-up">mc-up</a>: Deployment and management utilities
├── <a href="./mc-wheel">mc-wheel</a>: Database utilities
├── <a href="./ops">ops</a>: Various operational packages
├── <a href="./packages">packages</a>
│   └── <a href="./packages/contracts-bedrock">contracts-bedrock</a>: Mind Chain smart contracts
├── <a href="./releases">releases</a>: Pinned sources of mc-deployer and cannon releases
└── <a href="./rust">rust</a>: Rust components — a unified Cargo workspace
    ├── <a href="./rust/kona">kona</a>: State-transition (fault proof) program and rollup node in Rust
    ├── <a href="./rust/mc-reth">mc-reth</a>: Mind Chain execution client built on reth
    ├── <a href="./rust/mc-reth-test-engine">mc-reth-test-engine</a>: Minimal execution-layer engine for Go action tests (WIP)
    ├── <a href="./rust/op-revm">op-revm</a>: Rollup EVM implementation (revm fork)
    ├── <a href="./rust/op-alloy">op-alloy</a>: Rollup types and providers for the alloy ecosystem
    ├── <a href="./rust/alloy-op-evm">alloy-op-evm</a>: Rollup EVM support for the alloy ecosystem
    ├── <a href="./rust/alloy-op-hardforks">alloy-op-hardforks</a>: Hardfork definitions for the alloy ecosystem
    └── <a href="./rust/lokahi">lokahi</a>: Rust multi-chain consensus-layer host (WIP)
</pre>

## Releases

Every push to `main` and every `v*` tag builds and pushes the service images to
`ghcr.io/arunjot12/mind-chain/<service>:<tag>`. To cut a release:

```
git tag v1.1.0 && git push origin v1.1.0
```

The upstream versions each service is based on are pinned in [mindchain/versions.env](./mindchain/versions.env).

## Downloading & Shallow-Cloning

To fetch a specific commit, branch or tag `$REF` without the full history:

```
git clone --no-checkout --depth 1 --shallow-submodules https://github.com/arunjot12/mind-chain.git
cd mind-chain
git fetch --depth 1 origin "$REF"
git checkout "$REF"
```

## License

MIT, see [LICENSE](./LICENSE). Mind Chain is derived from the open-source OP Stack by Optimism,
which is also MIT-licensed; its copyright notice is kept in [LICENSE](./LICENSE) as that license requires.
