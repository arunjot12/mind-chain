package base

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

func TestPublicRpcAdvance(gt *testing.T) {
	t := devtest.ParallelT(gt)
	sys := presets.NewMinimal(t)

	sys.L2Chain.PublicRPC().Advanced(eth.Unsafe, 5)
}
