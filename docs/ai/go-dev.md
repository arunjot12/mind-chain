# Go Service Development

Guidance for AI agents working with Go code in the Optimism monorepo. See [dev-workflow.md](dev-workflow.md) for shared code simplicity principles, tool versions, and PR workflow.

## Build System

Each Go service has its own justfile — run `just --list` in any service directory to see available targets.

```bash
# Build a single service (pattern: just ./<service>/<binary>)
just ./mc-node/mc-node

# Build all Go components
just build-go
```

### Built-in chain configs

Mind Chain does not ship the upstream superchain-registry. `mc-core/superchain` starts
with an empty built-in chain list, so `--network <name>` lookups report an unknown chain
and every service is configured from its own rollup config and genesis files
(`--rollup.config`, `--l2.genesis`, or a numeric chain ID where a flag accepts one).

### Running Tests

```bash
# Test a single service
cd <service> && just test

# Test specific packages
go test ./mc-node/rollup/derive/...

# Run the full test suite from the repo root
just go-tests
```

### Generating Mocks

Some service justfiles have a `generate-mocks` target:

```bash
cd <service> && just generate-mocks
```

## Conventions

Apply the [shared code simplicity principles](dev-workflow.md#code-simplicity-and-design) with
these Go-specific conventions:

- **Prefer concrete dependencies.** Introduce an interface when a real boundary requires it, rather
  than wrapping a single implementation only to forward its methods. Keep `nil` where absence has
  domain meaning; supply required dependencies at construction.
- **Let packages provide naming context.** Prefer `derivation.Trusted` over
  `derivation.TrustedDerivation`.
- **OP type encodings**: a change to an OP type's wire format, hash rule, or codec accept-set has to agree with `rust/op-alloy` (the types mc-reth and kona consume), not just op-geth — see [Cross-implementation parity](opgeth-decoupling.md#cross-implementation-parity).
- **Pointers to values**: use `ptr.New(v)` from `github.com/mindchain/mind-chain/mc-service/ptr` to take the address of a literal or expression — common for optional `*uint64` config fields like fork-activation times (`cfg.SomeTime = ptr.New(uint64(123))`). Don't define a local `ptr`/`ptrTo` helper; the shared one avoids per-package duplicates, and a local `func ptr` collides with importing the `ptr` package in the same package.

## Linting

The repo uses a **custom golangci-lint build** with additional analyzer plugins. The standard `golangci-lint` binary will not catch all issues — always lint through `just`.

```bash
# Lint (also verifies compilation and module tidiness)
just lint-go

# Lint with auto-fix
just lint-go-fix
```

The linter configuration is in `.golangci.yaml` — read it when you need specifics on which linters are enabled and how they're scoped.

## Before Every Commit

Run these checks before committing Go changes. Fix all issues — CI enforces zero warnings.

1. **Lint** — this also verifies the code compiles and modules are tidy:
   ```bash
   just lint-go
   ```

2. **Test** — run tests for changed packages:
   ```bash
   cd <service> && just test
   ```

3. **op-geth ratchet** — if the change uses a go-ethereum API it didn't use before, check that
   the API exists upstream, not only in op-geth (see
   [opgeth-decoupling.md §19](opgeth-decoupling.md#19-ci-ratchet--continuously-rehearsing-the-cutover)):
   ```bash
   just opgeth-ratchet
   ```
