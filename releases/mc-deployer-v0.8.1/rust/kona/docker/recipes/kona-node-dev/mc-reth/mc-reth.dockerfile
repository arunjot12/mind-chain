FROM us-docker.pkg.dev/oplabs-tools-artifacts/images/op-reth:nightly AS reth

FROM ubuntu:latest

COPY --from=reth /usr/local/bin/mc-reth /usr/local/bin/mc-reth

WORKDIR /

COPY jwttoken/jwt.hex /

ENTRYPOINT [ "mc-reth" ]
