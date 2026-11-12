package interop

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

func newSupernodeInteropWithTimeTravel(t devtest.T, delaySeconds uint64) *presets.TwoL2SupernodeInterop {
	return presets.NewTwoL2SupernodeInterop(t, delaySeconds,
		presets.WithTimeTravelEnabled(),
	)
}
