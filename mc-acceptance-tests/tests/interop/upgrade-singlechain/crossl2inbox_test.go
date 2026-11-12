package upgrade

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mindchain/mind-chain/mc-chain-ops/devkeys"
	"github.com/mindchain/mind-chain/mc-chain-ops/genesis"
	"github.com/mindchain/mind-chain/mc-core/forks"
	"github.com/mindchain/mind-chain/mc-core/predeploys"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-e2e/e2eutils/intentbuilder"
)

func TestPostInbox(gt *testing.T) {
	t := devtest.ParallelT(gt)
	offset := uint64(30)
	sys := presets.NewSingleChainInterop(t, presets.WithDeployerOptions(
		func(p devtest.T, keys devkeys.Keys, builder intentbuilder.Builder) {
			for _, l2Cfg := range builder.L2s() {
				l2Cfg.WithForkAtOffset(forks.Lagoon, &offset)
			}
		},
	))
	devtest.RunParallel(t, sys.L2Networks(), func(t devtest.T, net *dsl.L2Network) {
		require := t.Require()
		el := net.PrimaryEL()

		activationBlock := net.AwaitActivation(t, forks.Lagoon)
		require.NotZero(activationBlock, "must not activate interop at genesis")

		pre := activationBlock.Number - 1

		verifyNoCrossL2InboxAtBlock := func(blockNum uint64) {
			net.PublicRPC().WaitForBlockNumber(blockNum)
			implAddrBytes, err := el.EthClient().GetStorageAt(t.Ctx(), predeploys.CrossL2InboxAddr,
				genesis.ImplementationSlot, hexutil.Uint64(blockNum).String())
			require.NoError(err)
			implAddr := common.BytesToAddress(implAddrBytes[:])
			require.Equal(common.Address{}, implAddr, "Should not have CrossL2Inbox implementation")
		}

		verifyNoCrossL2InboxAtBlock(pre)
		verifyNoCrossL2InboxAtBlock(activationBlock.Number)
		verifyNoCrossL2InboxAtBlock(activationBlock.Number + 1)
	})
}
