package presets

import (
	"github.com/mindchain/mind-chain/mc-chain-ops/devkeys"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
	"github.com/mindchain/mind-chain/mc-e2e/e2eutils/intentbuilder"
	"github.com/mindchain/mind-chain/mc-node/rollup"
)

// WithJovianAtGenesis configures all L2s to activate the Jovian fork at genesis in sysgo mode.
func WithJovianAtGenesis() stack.CommonOption {
	return stack.MakeCommon(sysgo.WithDeployerOptions(
		func(p devtest.P, _ devkeys.Keys, builder intentbuilder.Builder) {
			for _, l2Cfg := range builder.L2s() {
				l2Cfg.WithForkAtGenesis(rollup.Jovian)
			}
		},
	))
}
