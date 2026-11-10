#!/usr/bin/env bash
# Quick alternative to build-tools.sh: downloads the official op-deployer and cannon binaries into ./bin
set -euo pipefail
cd "$(dirname "$0")/../.."
source mindchain/versions.env
mkdir -p bin
os=$(uname -s | tr '[:upper:]' '[:lower:]'); arch=$(uname -m); [ "$arch" = x86_64 ] && arch=amd64; [ "$arch" = aarch64 ] && arch=arm64
v=${OP_DEPLOYER_VERSION#v}
curl -fL --retry 5 "https://github.com/ethereum-optimism/optimism/releases/download/op-deployer/${OP_DEPLOYER_VERSION}/op-deployer-${v}-${os}-${arch}.tar.gz" | tar -xz -C bin --strip-components=1
img=us-docker.pkg.dev/oplabs-tools-artifacts/images/cannon:${CANNON_VERSION}
docker pull "$img"; id=$(docker create "$img"); docker cp "$id:/usr/local/bin/." bin/; docker rm "$id" >/dev/null
echo "done: bin/op-deployer bin/cannon"
