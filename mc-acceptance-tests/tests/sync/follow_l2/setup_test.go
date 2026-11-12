package follow_l2

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
)

func newSingleChainTwoVerifiersFollowL2(t devtest.T) *presets.SingleChainTwoVerifiers {
	return presets.NewSingleChainTwoVerifiersWithoutCheck(t,
		presets.WithGlobalL2CLOption(sysgo.L2CLOptionFn(
			func(_ devtest.T, _ sysgo.ComponentTarget, cfg *sysgo.L2CLConfig) {
				cfg.NoDiscovery = true
			})),
	)
}
