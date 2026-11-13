#!/usr/bin/env bash
# Removes every build output that build-all.sh produces, so the next build-all.sh starts from scratch.
#
#   bash mindchain/scripts/build-clean.sh            clean all build outputs in this repo
#   bash mindchain/scripts/build-clean.sh --caches   also clear the Go build cache (shared, outside the repo)
#
# Only git-ignored files under the build-output paths below are deleted (git clean -X), so tracked
# and new, not-yet-committed source files are never touched.
set -euo pipefail
unset GOROOT
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
CACHES=false
for a in "$@"; do
  case "$a" in
    --caches) CACHES=true ;;
    -h|--help) sed -n 2,8p "$0"; exit 0 ;;
    *) echo "ERROR: unknown option '$a' (use --caches or --help)"; exit 1 ;;
  esac
done

say() { echo; echo "==> $*"; }
# clean <dir> <paths...>: delete git-ignored files under the given paths of the repo in <dir>
clean() { local d="$1"; shift; (cd "$d" && git clean -fdX -q -- "$@"); }

say "bin/ (built binaries)"
clean "$ROOT" bin

say "Go service binaries: mc-node mc-batcher mc-proposer mc-challenger"
clean "$ROOT" mc-node/bin mc-batcher/bin mc-proposer/bin mc-challenger/bin

say "cannon (releases/cannon-v1.6.0)"
clean "$ROOT/releases/cannon-v1.6.0" cannon/bin cannon/multicannon/embeds

say "mc-deployer + contracts (releases/mc-deployer-v0.8.1)"
D="$ROOT/releases/mc-deployer-v0.8.1"
clean "$D" mc-deployer/bin mc-deployer/pkg/deployer/artifacts/forge-artifacts
clean "$D" packages/contracts-bedrock/forge-artifacts packages/contracts-bedrock/artifacts packages/contracts-bedrock/cache

say "Rust: mc-reth + kona-host (cargo clean)"
if command -v cargo >/dev/null 2>&1; then
  (cd "$ROOT/rust" && cargo clean)
else
  clean "$ROOT" rust/target
fi

if $CACHES; then
  say "Go build cache"
  go clean -cache
fi

say "Clean done. Now run: bash mindchain/scripts/build-all.sh"
