#!/usr/bin/env bash
# Builds mc-deployer v0.8.1 and cannon v1.6.0 from releases/ into ./bin
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

cd "$ROOT/releases/mc-deployer-v0.8.1/mc-deployer"
just GITCOMMIT="$GITCOMMIT" GITDATE="$GITDATE" VERSION="$OP_DEPLOYER_VERSION" build
cp bin/mc-deployer "$ROOT/bin/mc-deployer"
echo "mc-deployer $OP_DEPLOYER_VERSION -> bin/mc-deployer"
