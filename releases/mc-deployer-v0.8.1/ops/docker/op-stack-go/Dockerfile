# automatically set by buildkit, can be changed with --platform flag
# see https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope
# TARGETOS
# TARGETARCH
# TARGETPLATFORM
# BUILDPLATFORM

# All target images use this as base image, and add the final build results.
# It will default to the target platform.
ARG TARGET_BASE_IMAGE=chainguard/static:latest

# mc-challenger runs on Wolfi: glibc for kona-host, busybox for the chart's
# /bin/sh entrypoint, plus a bundled musl loader for the musl-linked cannon embeds.
ARG CHALLENGER_TARGET_BASE_IMAGE=chainguard/wolfi-base:latest@sha256:02dab76bd852a70556b5b2002195c8a5fdab77d323c433bf6642aab080489795

# builder_foundry does not need to be built on $BUILDPLATFORM, as foundry produces static binaries.
FROM alpine:3.21 AS builder_foundry

# Build args are only-available in top-level by default, so we'll need to declare it inside the stage too
ARG TARGETARCH

# apk has no built-in download retry; loop so a transient CDN drop doesn't flake CI (~5 min budget).
RUN n=0; until apk add --no-cache curl bash git jq; do n=$((n+1)); [ "$n" -ge 15 ] && exit 1; echo "apk add retry $n/15 in 20s" >&2; sleep 20; done

COPY ./mc-deployer/pkg/deployer/forge/version.json /tmp/mc-deployer-versions.json

# Verify forge checksums
RUN <<-EOF
    #!/bin/bash
    set -e

    FOUNDRY_VERSION="$(jq -r '.forge' /tmp/mc-deployer-versions.json)"
    FOUNDRY_PLATFORM=alpine_${TARGETARCH}
    FOUNDRY_URL="https://github.com/foundry-rs/foundry/releases/download/${FOUNDRY_VERSION}/foundry_${FOUNDRY_VERSION}_${FOUNDRY_PLATFORM}.tar.gz"
    DOWNLOAD_PATH="/tmp/forge.tar.gz"

    # Download the artifact archive
    curl -fL --retry 5 --retry-all-errors --retry-delay 2 -o ${DOWNLOAD_PATH} "${FOUNDRY_URL}"

    # Ensure the checksum matches our expectations
    ACTUAL_SHA="$(sha256sum "${DOWNLOAD_PATH}" | awk '{print $1}')"
    EXPECTED_SHA="$(jq -r --arg arch "${FOUNDRY_PLATFORM}" '.checksums[$arch]' /tmp/mc-deployer-versions.json)"

    [ "${ACTUAL_SHA}" == "${EXPECTED_SHA}" ] || (echo "Forge checksum mismatch: Expected $EXPECTED_SHA, got $ACTUAL_SHA" && exit 1)

    # Extract the archive
    tar -xzf "${DOWNLOAD_PATH}" -C /usr/local/bin forge
EOF

# Display installed version for visibility
RUN forge --version

# We may be cross-building for another platform. Specify which platform we need as builder.
FROM --platform=$BUILDPLATFORM golang:1.26.5-alpine3.23 AS builder

# Enforce that builds use the committed go.mod/go.sum exactly and never
# silently mutate them. Fails the build if the lockfile is out of date.
ENV GOFLAGS=-mod=readonly

# apk has no built-in download retry; loop so a transient CDN drop doesn't flake CI (~5 min budget).
RUN n=0; until apk add --no-cache curl tar gzip make gcc musl-dev linux-headers git jq yq-go zip bash just; do n=$((n+1)); [ "$n" -ge 15 ] && exit 1; echo "apk add retry $n/15 in 20s" >&2; sleep 20; done

ARG TARGETARCH

# We copy the go.mod/sum first, so the `go mod download` does not have to re-run if dependencies do not change.
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

RUN echo "go mod cache: $(go env GOMODCACHE)"
RUN echo "go build cache: $(go env GOCACHE)"

# warm-up the cache
# go mod download has no built-in retry; loop so a transient module proxy drop doesn't flake CI.
# Already-downloaded modules are cached, so a retry only re-fetches what's missing.
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build n=0; until go mod download; do n=$((n+1)); [ "$n" -ge 5 ] && exit 1; echo "go mod download retry $n/5 in 10s" >&2; sleep 10; done

# NOTE: the Dockerfile.dockerignore file effectively describes all dependencies
COPY . /app

# Most op-stack-go services embed mc-core/superchain/superchain-configs.zip via
# //go:embed. It's gitignored and regenerated from the superchain-registry
# submodule, whose data is copied into the build context (see Dockerfile.dockerignore).
# Build it here in the shared builder stage — git-free, verifying against the
# committed .sha256 — so every Go build stage below links it.
RUN bash mc-core/superchain/sync-superchain.sh

# We avoid copying the full .git dir into the build for just some metadata.
# Instead, specify:
# --build-arg GIT_COMMIT=$(git rev-parse HEAD)
# --build-arg GIT_DATE=$(git show -s --format='%ct')
ARG GIT_COMMIT
ARG GIT_DATE

ARG TARGETOS

# Show just version
RUN just --version

# separate docker-builds:
# - proxyd
# - any JS/TS/smart-contract builds

# Build the Go services, utilizing caches and share the many common packages.
# The "id" defaults to the value of "target", the cache will thus be reused during this build.
# "sharing" defaults to "shared", the cache will thus be available to other concurrent docker builds.

FROM --platform=$BUILDPLATFORM us-docker.pkg.dev/oplabs-tools-artifacts/images/cannon:v1.0.0 AS cannon-builder-v1-0-0
FROM --platform=$BUILDPLATFORM us-docker.pkg.dev/oplabs-tools-artifacts/images/cannon:v1.2.0 AS cannon-builder-v1-2-0
FROM --platform=$BUILDPLATFORM us-docker.pkg.dev/oplabs-tools-artifacts/images/cannon:v1.3.0 AS cannon-builder-v1-3-0
FROM --platform=$BUILDPLATFORM us-docker.pkg.dev/oplabs-tools-artifacts/images/cannon:v1.4.0 AS cannon-builder-v1-4-0
FROM --platform=$BUILDPLATFORM us-docker.pkg.dev/oplabs-tools-artifacts/images/cannon:v1.6.0 AS cannon-builder-v1-6-0

FROM --platform=$BUILDPLATFORM builder AS cannon-builder
ARG CANNON_VERSION=v0.0.0
# Copy cannon binaries from previous versions
COPY --from=cannon-builder-v1-0-0 /usr/local/bin/cannon ./cannon/multicannon/embeds/cannon-0
COPY --from=cannon-builder-v1-3-0 /usr/local/bin/cannon-1 ./cannon/multicannon/embeds/cannon-1
COPY --from=cannon-builder-v1-4-0 /usr/local/bin/cannon-2 ./cannon/multicannon/embeds/cannon-2
COPY --from=cannon-builder-v1-2-0 /usr/local/bin/cannon-3 ./cannon/multicannon/embeds/cannon-3
COPY --from=cannon-builder-v1-3-0 /usr/local/bin/cannon-4 ./cannon/multicannon/embeds/cannon-4
COPY --from=cannon-builder-v1-4-0 /usr/local/bin/cannon-5 ./cannon/multicannon/embeds/cannon-5
COPY --from=cannon-builder-v1-4-0 /usr/local/bin/cannon-6 ./cannon/multicannon/embeds/cannon-6
COPY --from=cannon-builder-v1-6-0 /usr/local/bin/cannon-7 ./cannon/multicannon/embeds/cannon-7
# Build current binaries
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd cannon && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$CANNON_VERSION" just cannon

FROM --platform=$BUILDPLATFORM builder AS op-wheel-builder
ARG OP_WHEEL_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-wheel && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_WHEEL_VERSION" just mc-wheel

FROM --platform=$BUILDPLATFORM builder AS mc-node-builder
ARG OP_NODE_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-node && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_NODE_VERSION" just mc-node

FROM --platform=$BUILDPLATFORM builder AS mc-challenger-builder
ARG OP_CHALLENGER_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-challenger && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_CHALLENGER_VERSION" just mc-challenger

FROM --platform=$BUILDPLATFORM builder AS op-dispute-mon-builder
ARG OP_DISPUTE_MON_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-dispute-mon && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_DISPUTE_MON_VERSION" just mc-dispute-mon

FROM --platform=$BUILDPLATFORM builder AS mc-batcher-builder
ARG OP_BATCHER_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-batcher && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_BATCHER_VERSION" just mc-batcher

FROM --platform=$BUILDPLATFORM builder AS mc-proposer-builder
ARG OP_PROPOSER_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-proposer && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_PROPOSER_VERSION" just mc-proposer

FROM --platform=$BUILDPLATFORM builder AS op-conductor-builder
ARG OP_CONDUCTOR_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-conductor && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_CONDUCTOR_VERSION" just mc-conductor

FROM --platform=$BUILDPLATFORM builder AS da-server-builder
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-alt-da && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE just da-server

FROM --platform=$BUILDPLATFORM builder AS op-supernode-builder
ARG OP_SUPERNODE_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-supernode && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_SUPERNODE_VERSION" just mc-supernode

FROM --platform=$BUILDPLATFORM builder AS op-interop-filter-builder
ARG OP_INTEROP_FILTER_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-interop-filter && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_INTEROP_FILTER_VERSION" just mc-interop-filter

FROM --platform=$BUILDPLATFORM builder AS op-test-sequencer-builder
ARG OP_TEST_SEQUENCER_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-test-sequencer && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_TEST_SEQUENCER_VERSION" just mc-test-sequencer

# Compile contracts in a glibc-based image (solc requires glibc, crashes on Alpine musl).
# Uses BUILDPLATFORM so it runs natively (no QEMU). Contract artifacts are platform-independent.
FROM --platform=$BUILDPLATFORM golang:1.26.5-bookworm AS mc-deployer-contracts
ARG BUILDARCH
ENV GOFLAGS=-mod=readonly
RUN apt-get -o Acquire::Retries=8 update && apt-get -o Acquire::Retries=8 install -y --no-install-recommends curl jq git && rm -rf /var/lib/apt/lists/*
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
WORKDIR /app
# go mod download has no built-in retry; loop so a transient module proxy drop doesn't flake CI.
RUN --mount=type=cache,target=/go/pkg/mod n=0; until go mod download; do n=$((n+1)); [ "$n" -ge 5 ] && exit 1; echo "go mod download retry $n/5 in 10s" >&2; sleep 10; done
COPY . /app
# Install forge for the build platform (linux/glibc variant)
RUN <<-EOF
    set -e
    FOUNDRY_VERSION="$(jq -r '.forge' mc-deployer/pkg/deployer/forge/version.json)"
    FOUNDRY_PLATFORM=linux_${BUILDARCH}
    FOUNDRY_URL="https://github.com/foundry-rs/foundry/releases/download/${FOUNDRY_VERSION}/foundry_${FOUNDRY_VERSION}_${FOUNDRY_PLATFORM}.tar.gz"
    DOWNLOAD_PATH="/tmp/forge.tar.gz"
    curl -fL --retry 5 --retry-all-errors --retry-delay 2 -o "${DOWNLOAD_PATH}" "${FOUNDRY_URL}"
    ACTUAL_SHA="$(sha256sum "${DOWNLOAD_PATH}" | awk '{print $1}')"
    EXPECTED_SHA="$(jq -r --arg arch "${FOUNDRY_PLATFORM}" '.checksums[$arch]' mc-deployer/pkg/deployer/forge/version.json)"
    [ "${ACTUAL_SHA}" = "${EXPECTED_SHA}" ] || (echo "Forge checksum mismatch: Expected $EXPECTED_SHA, got $ACTUAL_SHA" && exit 1)
    tar -xzf "${DOWNLOAD_PATH}" -C /usr/local/bin forge
    rm "${DOWNLOAD_PATH}"
EOF
# Pre-install solc versions serially (matching .circleci install-solc-compilers) so
# forge's parallel jobs never race installing them (foundry-rs/foundry#4736).
# `forge build --use <version>` in an empty dir installs one version, compiles nothing.
RUN <<-EOF
    set -e
    cd "$(mktemp -d)"
    # Retry: forge doesn't retry the solc download itself.
    for version in 0.8.15 0.8.19 0.8.25 0.8.28; do
      n=0
      until forge build --use "${version}"; do
        n=$((n+1))
        [ "${n}" -ge 5 ] && exit 1
        echo "solc ${version} install retry ${n}/5 in 2s" >&2
        sleep 2
      done
    done
EOF
# Compile contracts and create artifacts tarball
RUN cd packages/contracts-bedrock && forge build --skip "/**/test/**"
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  go run ./mc-deployer/pkg/deployer/artifacts/cmd/mktar \
    -base ./packages/contracts-bedrock \
    -out ./mc-deployer/pkg/deployer/artifacts/forge-artifacts/artifacts.tzst \
    --exclude="*.t.sol"

FROM --platform=$BUILDPLATFORM builder AS mc-deployer-builder
ARG OP_DEPLOYER_VERSION=v0.0.0
# Copy pre-compiled contract artifacts from the glibc-based contracts stage
COPY --from=mc-deployer-contracts /app/mc-deployer/pkg/deployer/artifacts/forge-artifacts/artifacts.tzst mc-deployer/pkg/deployer/artifacts/forge-artifacts/artifacts.tzst
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  just GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_DEPLOYER_VERSION" mc-deployer/build-go

FROM --platform=$BUILDPLATFORM builder AS op-dripper-builder
ARG OP_DRIPPER_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  cd mc-dripper && GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_DRIPPER_VERSION" just mc-dripper

FROM --platform=$BUILDPLATFORM builder AS op-faucet-builder
ARG OP_FAUCET_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_FAUCET_VERSION" \
  just mc-faucet/mc-faucet

FROM --platform=$BUILDPLATFORM builder AS op-interop-mon-builder
ARG OP_INTEROP_MON_VERSION=v0.0.0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
  GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_INTEROP_MON_VERSION" \
  just mc-interop-mon/mc-interop-mon

# The Rust version must match rust/rust-toolchain.toml. We don't use "latest" to ensure reproducibility
FROM --platform=$BUILDPLATFORM rust:1.95 AS kona-host-builder
ARG TARGETARCH
# Install build dependencies and cross-compilation toolchains
RUN apt-get -o Acquire::Retries=8 update && apt-get -o Acquire::Retries=8 install -y --no-install-recommends \
  build-essential \
  libssl-dev \
  clang \
  pkg-config \
  gcc-aarch64-linux-gnu \
  gcc-x86-64-linux-gnu \
  && rm -rf /var/lib/apt/lists/*
# Copy the entire rust workspace since kona depends on sibling crates (op-alloy, alloy-op-evm, etc.)
COPY ./rust /rust
# kona-hardforks build.rs probes ancestor paths for the embedded NUT bundle JSON.
# Keep mc-core as a sibling of /rust so Docker builds see the monorepo layout.
COPY ./mc-core/nuts/bundles /mc-core/nuts/bundles
WORKDIR /rust
# Map TARGETARCH to Rust target triple, install it, and build
RUN --mount=type=cache,target=/usr/local/cargo/registry --mount=type=cache,target=/rust/target \
  case "$TARGETARCH" in \
    amd64) RUST_TARGET=x86_64-unknown-linux-gnu ;; \
    arm64) RUST_TARGET=aarch64-unknown-linux-gnu ;; \
    *) echo "Unsupported architecture: $TARGETARCH" && exit 1 ;; \
  esac && \
  rustup target add "$RUST_TARGET" && \
  RUSTFLAGS="-C target-cpu=generic" cargo build --locked --bin kona-host --profile release-perf --target "$RUST_TARGET" && \
  cp "/rust/target/$RUST_TARGET/release-perf/kona-host" /kona-host

FROM dhi.io/alpine-base:3.23 AS cannon-target
USER 0
COPY --from=cannon-builder /app/cannon/bin/cannon /usr/local/bin/
COPY --from=cannon-builder /app/cannon/multicannon/embeds/* /usr/local/bin/
CMD ["cannon"]

FROM $TARGET_BASE_IMAGE AS op-wheel-target
USER 0
COPY --from=op-wheel-builder /app/mc-wheel/bin/mc-wheel /usr/local/bin/
CMD ["mc-wheel"]

FROM $TARGET_BASE_IMAGE AS mc-node-target
USER 0
COPY --from=mc-node-builder /app/mc-node/bin/mc-node /usr/local/bin/
CMD ["mc-node"]

# The historical multicannon embeds are Alpine/musl-linked and Wolfi ships no musl
# package, so bundle the loader+libc (as the melange multicannon pipeline does).
FROM alpine:3.21@sha256:48b0309ca019d89d40f670aa1bc06e426dc0931948452e8491e3d65087abc07d AS musl-libs
RUN mkdir -p /musl && cp -a /lib/ld-musl-*.so.1 /lib/libc.musl-*.so.1 /musl/

# Also produce an mc-challenger loaded with kona
FROM $CHALLENGER_TARGET_BASE_IMAGE AS mc-challenger-target
RUN attempt=1; delay=1; until apk add --no-cache ca-certificates-bundle ca-certificates busybox glibc libgcc libstdc++ openssl; do [ "$attempt" -ge 11 ] && exit 1; echo "apk add attempt $attempt/11 failed; retrying in ${delay}s" >&2; sleep "$delay"; attempt=$((attempt+1)); delay=$((delay*2)); [ "$delay" -le 60 ] || delay=60; done
COPY --from=musl-libs /musl/ /lib/
COPY --from=mc-challenger-builder /app/mc-challenger/bin/mc-challenger /usr/local/bin/
# Copy in cannon
COPY --from=cannon-builder /app/cannon/bin/cannon /usr/local/bin/
ENV OP_CHALLENGER_CANNON_BIN=/usr/local/bin/cannon
# Copy in kona-host (built from local kona/ source)
COPY --from=kona-host-builder /kona-host /usr/local/bin/kona-host
ENV OP_CHALLENGER_ASTERISC_KONA_SERVER=/usr/local/bin/kona-host
ENV OP_CHALLENGER_CANNON_KONA_SERVER=/usr/local/bin/kona-host
CMD ["mc-challenger"]

FROM $TARGET_BASE_IMAGE AS op-dispute-mon-target
USER 0
COPY --from=op-dispute-mon-builder /app/mc-dispute-mon/bin/mc-dispute-mon /usr/local/bin/
CMD ["mc-dispute-mon"]

FROM $TARGET_BASE_IMAGE AS mc-batcher-target
USER 0
COPY --from=mc-batcher-builder /app/mc-batcher/bin/mc-batcher /usr/local/bin/
CMD ["mc-batcher"]

FROM $TARGET_BASE_IMAGE AS mc-proposer-target
USER 0
COPY --from=mc-proposer-builder /app/mc-proposer/bin/mc-proposer /usr/local/bin/
CMD ["mc-proposer"]

FROM $TARGET_BASE_IMAGE AS op-conductor-target
USER 0
COPY --from=op-conductor-builder /app/mc-conductor/bin/mc-conductor /usr/local/bin/
CMD ["mc-conductor"]

FROM $TARGET_BASE_IMAGE AS da-server-target
USER 0
COPY --from=da-server-builder /app/mc-alt-da/bin/da-server /usr/local/bin/
CMD ["da-server"]

FROM $TARGET_BASE_IMAGE AS op-supernode-target
USER 0
COPY --from=op-supernode-builder /app/mc-supernode/bin/mc-supernode /usr/local/bin/
CMD ["mc-supernode"]

FROM $TARGET_BASE_IMAGE AS op-interop-filter-target
USER 0
COPY --from=op-interop-filter-builder /app/mc-interop-filter/bin/mc-interop-filter /usr/local/bin/
CMD ["mc-interop-filter"]

FROM $TARGET_BASE_IMAGE AS op-test-sequencer-target
USER 0
COPY --from=op-test-sequencer-builder /app/mc-test-sequencer/bin/mc-test-sequencer /usr/local/bin/
CMD ["mc-test-sequencer"]

FROM dhi.io/alpine-base:3.23 AS mc-deployer-target
USER 0
COPY --from=builder_foundry /usr/local/bin/forge /usr/local/bin/
COPY --from=mc-deployer-builder /app/mc-deployer/bin/mc-deployer /usr/local/bin/
ENV FORGE_ENV=alpine
CMD ["mc-deployer"]

FROM $TARGET_BASE_IMAGE AS op-dripper-target
USER 0
COPY --from=op-dripper-builder /app/mc-dripper/bin/mc-dripper /usr/local/bin/
CMD ["mc-dripper"]

FROM $TARGET_BASE_IMAGE AS op-faucet-target
USER 0
COPY --from=op-faucet-builder /app/mc-faucet/bin/mc-faucet /usr/local/bin/
CMD ["mc-faucet"]

FROM $TARGET_BASE_IMAGE AS op-interop-mon-target
USER 0
COPY --from=op-interop-mon-builder /app/mc-interop-mon/bin/mc-interop-mon /usr/local/bin/
CMD ["mc-interop-mon"]
