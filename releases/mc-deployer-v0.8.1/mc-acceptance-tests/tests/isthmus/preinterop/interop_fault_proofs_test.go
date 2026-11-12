package preinterop

import (
	"testing"

	sfp "github.com/mindchain/mind-chain/mc-acceptance-tests/tests/superfaultproofs"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

func TestPreinteropFaultProofs(gt *testing.T) {
	t := devtest.ParallelT(gt)
	sys := newSimpleInteropPreinterop(t)
	sfp.RunSuperFaultProofTest(t, sys)
}

func TestPreinteropFaultProofs_TraceExtensionActivation(gt *testing.T) {
	t := devtest.ParallelT(gt)
	sys := newSimpleInteropPreinterop(t)
	sfp.RunTraceExtensionActivationTest(t, sys)
}

func TestPreinteropFaultProofs_UnsafeProposal(gt *testing.T) {
	t := devtest.ParallelT(gt)
	sys := newSimpleInteropPreinterop(t)
	sfp.RunUnsafeProposalTest(t, sys)
}

func TestPreinteropFaultProofs_VariedBlockTimes(gt *testing.T) {
	t := devtest.SerialT(gt)
	sys := presets.NewSimpleInteropIsthmusSuper(
		t,
		presets.WithL2BlockTimes(map[eth.ChainID]uint64{
			sysgo.DefaultL2AID: 1,
			sysgo.DefaultL2BID: 2,
		}),
	)
	sfp.RunVariedBlockTimesTest(t, sys)
}

func TestPreinteropFaultProofs_VariedBlockTimes_FasterChainB(gt *testing.T) {
	t := devtest.SerialT(gt)
	sys := presets.NewSimpleInteropIsthmusSuper(
		t,
		presets.WithL2BlockTimes(map[eth.ChainID]uint64{
			sysgo.DefaultL2AID: 2,
			sysgo.DefaultL2BID: 1,
		}),
	)
	sfp.RunVariedBlockTimesTest(t, sys)
}
