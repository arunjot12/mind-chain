# mc-interop-filter

A lightweight service that validates interop executing messages for op-geth or mc-reth transaction filtering.

Any reorg will trigger the failsafe which disables all interop transactions. If `--reorg-recovery-enabled` is set, reorg-triggered failsafe is automatically resolved by rewinding logs DB state to finalized.

## Usage

### Build from source

```bash
just mc-interop-filter
./bin/mc-interop-filter --help
```

### Run from source

```bash
go run ./cmd --help
```

### Build docker image

```bash
docker buildx bake mc-interop-filter
```
