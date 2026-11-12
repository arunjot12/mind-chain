package preinterop_singlechain

import (
	"testing"

	sfp "github.com/mindchain/mind-chain/mc-acceptance-tests/tests/superfaultproofs"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

func TestPreinteropSingleChainFaultProofs(gt *testing.T) {
	t := devtest.ParallelT(gt)
	sys := presets.NewSingleChainInteropIsthmusSuper(t)
	sfp.RunSingleChainSuperFaultProofSmokeTest(t, sys)
}
