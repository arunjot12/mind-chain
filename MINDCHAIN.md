# Mind Chain

One repository for all Mind Chain OP Stack services, based on ethereum-optimism/optimism.
Upstream versions: `mindchain/versions.env`.

| Path | Contents |
|---|---|
| repo root | mc-node v1.20.0, mc-batcher v1.18.0, mc-reth v2.8.0, mc-proposer v1.17.0, mc-challenger v1.11.0, kona-host v1.9.0 |
| releases/mc-deployer-v0.8.1 | exact mc-deployer v0.8.1 source (separate Go module) |
| releases/cannon-v1.6.0 | exact cannon v1.6.0 source (separate Go module) |

## Clone
    git clone --recurse-submodules https://github.com/arunjot12/mind-chain.git

## Daily work
    git switch main && git pull
    git switch -c <your-name>/<change>
    # ...edit...
    git add -A && git commit -m "mc-node: what changed"
    git push -u origin HEAD                  # then open a Pull Request to main

## Build locally
All binaries go to `./bin` (git-ignored, never pushed). Needs Go, Rust/cargo, just, make, forge (versions in mise.toml).

    bash mindchain/scripts/build-all.sh                  # everything
    bash mindchain/scripts/build-all.sh mc-node mc-reth  # only some

Builds: mc-node, mc-batcher, mc-proposer, mc-challenger (repo root), mc-deployer v0.8.1 and
cannon v1.6.0 (from releases/), mc-reth and kona-host (rust/). Each binary reports its upstream
version from mindchain/versions.env.

## Docker images (CI)
Every PR builds the images. Every push to main and every v* tag pushes them to
ghcr.io/arunjot12/mind-chain/<service>:<tag>.

## Release
    git tag v1.1.0 && git push origin v1.1.0


## Component names
| Mind Chain name | Upstream Optimism name |
|---|---|
| mc-node | op-node |
| mc-batcher | op-batcher |
| mc-proposer | op-proposer |
| mc-challenger | op-challenger |
| mc-deployer | op-deployer |
| mc-reth | op-reth |


Environment variables (OP_NODE_*, OP_BATCHER_*, ...) and metrics names keep the upstream names.
