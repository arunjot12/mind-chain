package super_opnode

import (
	"testing"

	sfp "github.com/ethereum-optimism/optimism/op-acceptance-tests/tests/superfaultproofs"
	"github.com/ethereum-optimism/optimism/op-devstack/devtest"
	"github.com/ethereum-optimism/optimism/op-devstack/presets"
)

// TestSingleChainSuperFaultProofsFromOpNode runs the single-chain super-fault-proof smoke
// test with super roots served by mc-node instead of op-supernode. It exercises the FPP and
// the challenger trace provider against the mc-node superroot_atTimestamp source.
func TestSingleChainSuperFaultProofsFromOpNode(gt *testing.T) {
	t := devtest.SerialT(gt)
	sys := presets.NewSingleChainInteropNoSupernode(t)
	sfp.RunSingleChainSuperFaultProofSmokeTest(t, sys)
}
