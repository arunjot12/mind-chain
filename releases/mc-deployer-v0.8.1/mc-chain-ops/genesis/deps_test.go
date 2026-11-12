package genesis

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-service/testutils/depguard"
)

// TestNoSuperchainImport keeps mc-chain-ops/genesis free of mc-core/superchain
// in its whole build closure. That package embeds the generated, gitignored
// superchain config bundle, which is absent from a clean module download — and
// the superchain-registry's ops tooling consumes this package as a Go module
// dependency, so any path to the bundle breaks its build. genesis imports
// mc-core/params directly, so it sits one import away from the bundle boundary.
func TestNoSuperchainImport(t *testing.T) {
	depguard.RequireNoTransitiveImport(t, ".",
		"github.com/mindchain/mind-chain/mc-core/superchain")
}
