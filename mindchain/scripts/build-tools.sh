#!/usr/bin/env bash
# Builds op-deployer v0.8.1 and cannon v1.6.0 from releases/ into ./bin
# Needs: go, just, forge (versions: releases/*/mise.toml)
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
source "$ROOT/mindchain/versions.env"
mkdir -p "$ROOT/bin"
GITCOMMIT=$(git -C "$ROOT" rev-parse HEAD); GITDATE=$(git -C "$ROOT" show -s --format=%ct)

cd "$ROOT/releases/cannon-v1.6.0/cannon"
make cannon GITCOMMIT="$GITCOMMIT" GITDATE="$GITDATE" VERSION="$CANNON_VERSION"
cp bin/cannon "$ROOT/bin/cannon"
echo "cannon $CANNON_VERSION -> bin/cannon"

cd "$ROOT/releases/op-deployer-v0.8.1/op-deployer"
# Build the chain-config bundle from v0.8.1's own superchain-registry (git would
# otherwise point the sync script at the root repo's registry).
SUPERCHAIN_REGISTRY_DIR="$ROOT/releases/op-deployer-v0.8.1/superchain-registry" \
  bash ../op-core/superchain/sync-superchain.sh
just GITCOMMIT="$GITCOMMIT" GITDATE="$GITDATE" VERSION="$OP_DEPLOYER_VERSION" build
cp bin/op-deployer "$ROOT/bin/op-deployer"
echo "op-deployer $OP_DEPLOYER_VERSION -> bin/op-deployer"
