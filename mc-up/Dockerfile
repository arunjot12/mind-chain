FROM debian:bookworm-20240812-slim
ENTRYPOINT ["/mc-up"]
COPY mc-up /mc-up

# Install ca-certificates so that HTTPS requests work
RUN apt-get -o Acquire::Retries=8 update && apt-get -o Acquire::Retries=8 install -y ca-certificates

# Symlink onto the PATH
RUN ln -s /mc-up /usr/local/bin/mc-up
