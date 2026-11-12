package reorgs

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

// TestMain creates the test-setups against the shared backend
func TestMain(m *testing.M) {
	// Other setups may be added here, hydrated from the same orchestrator
	presets.DoMain(m, presets.WithSimpleInterop())
}
