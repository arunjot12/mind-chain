#!/usr/bin/env bash
# Quick alternative to build-tools.sh: downloads the published mc-deployer release binary into ./bin
# and builds cannon from releases/cannon-v1.6.0 (needs go + make).
# mc-deployer releases are published by goreleaser (mc-deployer/.goreleaser.yaml) under the tag
# mc-deployer/<version>, e.g. mc-deployer/v0.8.1. Override the source repo with MIND_CHAIN_REPO=owner/name.
set -euo pipefail
cd "$(dirname "$0")/../.."
ROOT=$(pwd)
source mindchain/versions.env
REPO="${MIND_CHAIN_REPO:-arunjot12/mind-chain}"
mkdir -p bin
os=$(uname -s | tr '[:upper:]' '[:lower:]'); arch=$(uname -m); [ "$arch" = x86_64 ] && arch=amd64; [ "$arch" = aarch64 ] && arch=arm64
v=${OP_DEPLOYER_VERSION#v}
curl -fL --retry 5 "https://github.com/${REPO}/releases/download/mc-deployer%2Fv${v}/mc-deployer-${v}-${os}-${arch}.tar.gz" | tar -xz -C bin --strip-components=1
GITCOMMIT=$(git rev-parse HEAD); GITDATE=$(git show -s --format=%ct)
(cd releases/cannon-v1.6.0/cannon && make cannon GITCOMMIT="$GITCOMMIT" GITDATE="$GITDATE" VERSION="$CANNON_VERSION")
cp releases/cannon-v1.6.0/cannon/bin/cannon "$ROOT/bin/cannon"
echo "done: bin/mc-deployer bin/cannon"
