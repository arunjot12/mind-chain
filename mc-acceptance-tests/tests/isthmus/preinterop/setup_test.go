package preinterop

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

func newSimpleInteropPreinterop(t devtest.T) *presets.SimpleInterop {
	return presets.NewSimpleInteropIsthmusSuper(t)
}
