package depreqres

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-acceptance-tests/tests/depreqres/common"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-node/rollup/sync"
	"github.com/mindchain/mind-chain/mc-service/eth"

	safety "github.com/mindchain/mind-chain/mc-service/eth/safety"
)

func TestUnsafeChainNotStalling_DisabledReqRespSync(gt *testing.T) {
	t := devtest.ParallelT(gt)
	sys := presets.NewSingleChainMultiNodeWithoutCheck(t, common.ReqRespSyncDisabledOpts(sync.ELSync)...)
	l := t.Logger()

	l.Info("Confirm that the CL nodes are progressing the unsafe chain")
	delta := uint64(3)
	dsl.CheckAll(t,
		sys.L2CL.AdvancedFn(safety.LocalUnsafe, delta, 30),
		sys.L2CLB.AdvancedFn(safety.LocalUnsafe, delta, 30),
	)

	l.Info("Disconnect L2CL from L2CLB, and vice versa")
	sys.L2CLB.DisconnectPeer(sys.L2CL)
	sys.L2CL.DisconnectPeer(sys.L2CLB)

	l.Info("Wait for the CL nodes to be disconnected")
	sys.L2CL.WaitForPeerDisconnected(sys.L2CLB)
	sys.L2CLB.WaitForPeerDisconnected(sys.L2CL)

	l.Info("Wait for L2CLB unsafe head to stall after disconnect")
	sys.L2CLB.WaitForStall(safety.LocalUnsafe)

	l.Info("Confirm that the unsafe chain for L2CL advances while L2CLB remains stalled")
	sys.L2CL.AdvancedUnsafe(delta, 30)
	ssA := sys.L2CL.SyncStatus()

	l.Info("Re-connect L2CL to L2CLB")
	sys.L2CLB.ConnectPeer(sys.L2CL)
	sys.L2CL.ConnectPeer(sys.L2CLB)

	l.Info("Confirm that the unsafe chain for L2CLB can advance")
	sys.L2CLB.Reached(safety.LocalUnsafe, ssA.UnsafeL2.Number, 30)
	sys.L2ELB.Reached(eth.Unsafe, ssA.UnsafeL2.Number, 30)
}
