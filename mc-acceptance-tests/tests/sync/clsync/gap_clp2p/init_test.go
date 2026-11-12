package gap_clp2p

import (
	bss "github.com/mindchain/mind-chain/mc-batcher/batcher"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
)

func gapCLP2POpts() []presets.Option {
	return []presets.Option{
		presets.WithBatcherOption(func(_ sysgo.ComponentTarget, cfg *bss.CLIConfig) {
			// For stopping derivation, not to advance safe heads
			cfg.Stopped = true
		}),
	}
}

func newGapCLP2PSystem(t devtest.T) *presets.SingleChainMultiNode {
	return presets.NewSingleChainMultiNodeWithoutP2PWithoutCheck(t, gapCLP2POpts()...)
}
