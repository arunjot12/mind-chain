package reorg

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
)

func reorgOpts() []presets.Option {
	return []presets.Option{
		presets.WithGlobalL2CLOption(sysgo.L2CLOptionFn(func(_ devtest.T, _ sysgo.ComponentTarget, cfg *sysgo.L2CLConfig) {
			cfg.NoDiscovery = true
		})),
	}
}

func newReorgSystem(t devtest.T) *presets.SingleChainMultiNodeWithTestSeq {
	return presets.NewSingleChainMultiNodeWithTestSeq(t, reorgOpts()...)
}
