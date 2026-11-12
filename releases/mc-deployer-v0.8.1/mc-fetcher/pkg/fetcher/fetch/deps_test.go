package fetch

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-service/testutils/depguard"
)

// TestNoSuperchainImport keeps mc-fetcher's fetch package free of
// mc-core/superchain in its whole build closure. That package embeds the
// generated, gitignored superchain config bundle, which is absent from a clean
// module download — and the superchain-registry's ops tooling consumes this
// package as a Go module dependency, so any path to the bundle breaks its build.
func TestNoSuperchainImport(t *testing.T) {
	depguard.RequireNoTransitiveImport(t, ".",
		"github.com/mindchain/mind-chain/mc-core/superchain")
}
