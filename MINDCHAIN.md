# Mind Chain

One repository for all Mind Chain OP Stack services, based on ethereum-optimism/optimism.
Upstream versions: `mindchain/versions.env`.

| Path | Contents |
|---|---|
| repo root | op-node v1.20.0, op-batcher v1.18.0, op-reth v2.8.0, op-proposer v1.17.0, op-challenger v1.11.0, kona-host v1.9.0 |
| releases/op-deployer-v0.8.1 | exact op-deployer v0.8.1 source (separate Go module) |
| releases/cannon-v1.6.0 | exact cannon v1.6.0 source (separate Go module) |

## Clone
    git clone --recurse-submodules https://github.com/arunjot12/mind-chain.git

## Daily work
    git switch main && git pull
    git switch -c <your-name>/<change>
    # ...edit...
    git add -A && git commit -m "op-node: what changed"
    git push -u origin HEAD                  # then open a Pull Request to main

## Build locally (needs Go, Rust, just, forge, Docker; versions in mise.toml)
    just op-node            # -> op-node/bin/op-node
    just op-batcher         # -> op-batcher/bin/op-batcher
    just op-proposer        # -> op-proposer/bin/op-proposer
    just op-challenger      # -> op-challenger/bin/op-challenger
    cd rust && cargo build --release -p op-reth --bin op-reth
    cd rust && cargo build --release --bin kona-host
    mindchain/scripts/build-tools.sh              # op-deployer v0.8.1 + cannon v1.6.0 from releases/ -> ./bin
    mindchain/scripts/install-tools.sh            # or: download the official binaries instead

## Docker images (CI)
Every PR builds the images. Every push to main and every v* tag pushes them to
ghcr.io/arunjot12/mind-chain/<service>:<tag>.

## Release
    git tag v1.1.0 && git push origin v1.1.0
