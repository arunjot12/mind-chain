#!/usr/bin/env bash
# Builds the Mind Chain binaries and puts them in ./bin (git-ignored, never pushed).
#
#   bash mindchain/scripts/build-all.sh                  build everything
#   bash mindchain/scripts/build-all.sh mc-node mc-reth  build only these
#
# Components: mc-node mc-batcher mc-proposer mc-challenger mc-deployer cannon mc-reth kona-host
# Needs: go, just, make, cargo (Rust), forge (Foundry). Versions: mise.toml.
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
source "$ROOT/mindchain/versions.env"
BIN="$ROOT/bin"
ALL=(mc-node mc-batcher mc-proposer mc-challenger mc-deployer cannon mc-reth kona-host)
if [ $# -gt 0 ]; then COMPONENTS=("$@"); else COMPONENTS=("${ALL[@]}"); fi   # bash 3.2 safe

export GITCOMMIT GITDATE
GITCOMMIT=$(git -C "$ROOT" rev-parse HEAD)
GITDATE=$(git -C "$ROOT" show -s --format=%ct)

say()  { echo; echo "==> $*"; }
need() { command -v "$1" >/dev/null 2>&1 || { echo "ERROR: '$1' not found. Install it first (see mise.toml)."; exit 1; }; }
# copy to a temp file, then rename: replaces the file without touching a binary that is running
put()  { cp "$1" "$BIN/.$2.tmp"; chmod +x "$BIN/.$2.tmp"; mv -f "$BIN/.$2.tmp" "$BIN/$2"; echo "    -> bin/$2"; }

for c in "${COMPONENTS[@]}"; do
  case " ${ALL[*]} " in *" $c "*) ;; *) echo "ERROR: unknown component '$c'. Use: ${ALL[*]}"; exit 1;; esac
done
need git; need just
case " ${COMPONENTS[*]} " in *" mc-reth "*|*" kona-host "*) need cargo;; esac
case " ${COMPONENTS[*]} " in *" mc-node "*|*" mc-batcher "*|*" mc-proposer "*|*" mc-challenger "*|*" mc-deployer "*|*" cannon "*) need go;; esac
case " ${COMPONENTS[*]} " in *" mc-deployer "*) need forge;; esac
case " ${COMPONENTS[*]} " in *" cannon "*) need make;; esac
mkdir -p "$BIN"

for c in "${COMPONENTS[@]}"; do
  case "$c" in
    mc-node)
      say "mc-node $OP_NODE_VERSION"
      (cd "$ROOT" && VERSION="$OP_NODE_VERSION" just mc-node); put "$ROOT/mc-node/bin/mc-node" mc-node ;;
    mc-batcher)
      say "mc-batcher $OP_BATCHER_VERSION"
      (cd "$ROOT" && VERSION="$OP_BATCHER_VERSION" just mc-batcher); put "$ROOT/mc-batcher/bin/mc-batcher" mc-batcher ;;
    mc-proposer)
      say "mc-proposer $OP_PROPOSER_VERSION"
      (cd "$ROOT" && VERSION="$OP_PROPOSER_VERSION" just mc-proposer); put "$ROOT/mc-proposer/bin/mc-proposer" mc-proposer ;;
    mc-challenger)
      say "mc-challenger $OP_CHALLENGER_VERSION"
      (cd "$ROOT" && VERSION="$OP_CHALLENGER_VERSION" just mc-challenger); put "$ROOT/mc-challenger/bin/mc-challenger" mc-challenger ;;
    cannon)
      say "cannon $CANNON_VERSION (releases/cannon-v1.6.0)"
      (cd "$ROOT/releases/cannon-v1.6.0/cannon" && make cannon GITCOMMIT="$GITCOMMIT" GITDATE="$GITDATE" VERSION="$CANNON_VERSION")
      put "$ROOT/releases/cannon-v1.6.0/cannon/bin/cannon" cannon ;;
    mc-deployer)
      say "mc-deployer $OP_DEPLOYER_VERSION (releases/mc-deployer-v0.8.1)"
      D="$ROOT/releases/mc-deployer-v0.8.1"
      # chain-config bundle from v0.8.1's own superchain-registry
      (cd "$D/mc-deployer" && SUPERCHAIN_REGISTRY_DIR="$D/superchain-registry" bash ../op-core/superchain/sync-superchain.sh)
      (cd "$D/mc-deployer" && just GITCOMMIT="$GITCOMMIT" GITDATE="$GITDATE" VERSION="$OP_DEPLOYER_VERSION" build)
      put "$D/mc-deployer/bin/mc-deployer" mc-deployer ;;
    mc-reth)
      say "mc-reth $OP_RETH_VERSION (Rust, slow the first time)"
      (cd "$ROOT/rust" && cargo build --release -p mc-reth --bin mc-reth); put "$ROOT/rust/target/release/mc-reth" mc-reth ;;
    kona-host)
      say "kona-host $KONA_HOST_VERSION (Rust)"
      (cd "$ROOT/rust" && cargo build --release --bin kona-host); put "$ROOT/rust/target/release/kona-host" kona-host ;;
  esac
done

say "Built binaries in $BIN"
for c in "${COMPONENTS[@]}"; do
  printf '    %-14s ' "$c"; "$BIN/$c" --version 2>/dev/null | head -1 || echo "(no --version)"
done
git -C "$ROOT" check-ignore -q bin/mc-node 2>/dev/null || git -C "$ROOT" check-ignore -q "bin/${COMPONENTS[0]}" || echo "WARNING: bin/ is not git-ignored"
