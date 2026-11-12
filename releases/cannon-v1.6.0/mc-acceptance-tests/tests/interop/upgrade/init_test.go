package upgrade

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

func TestMain(m *testing.M) {
	presets.DoMain(m,
		presets.WithSimpleInterop(),
		presets.WithSuggestedInteropActivationOffset(60),
		presets.WithInteropNotAtGenesis())
}
