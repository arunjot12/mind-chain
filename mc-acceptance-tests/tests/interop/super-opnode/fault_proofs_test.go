package super_opnode

import (
	"testing"

	sfp "github.com/mindchain/mind-chain/mc-acceptance-tests/tests/superfaultproofs"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

// TestSingleChainSuperFaultProofsFromOpNode runs the single-chain super-fault-proof smoke
// test with super roots served by mc-node instead of mc-supernode. It exercises the FPP and
// the challenger trace provider against the mc-node superroot_atTimestamp source.
func TestSingleChainSuperFaultProofsFromOpNode(gt *testing.T) {
	t := devtest.SerialT(gt)
	sys := presets.NewSingleChainInteropNoSupernode(t)
	sfp.RunSingleChainSuperFaultProofSmokeTest(t, sys)
}
