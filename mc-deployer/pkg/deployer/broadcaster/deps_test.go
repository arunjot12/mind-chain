package broadcaster

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-service/testutils/depguard"
)

// TestNoScriptHostImport keeps broadcasters independent of the Go script host, so they
// serve every script engine.
func TestNoScriptHostImport(t *testing.T) {
	depguard.RequireNoTransitiveImport(t, ".",
		"github.com/mindchain/mind-chain/mc-chain-ops/script")
}
