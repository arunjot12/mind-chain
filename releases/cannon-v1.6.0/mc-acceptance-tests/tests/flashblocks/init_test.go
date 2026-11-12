package flashblocks

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

// TestMain creates the test-setups against the shared backend
func TestMain(m *testing.M) {
	presets.DoMain(m, presets.WithSimpleFlashblocks())
}
