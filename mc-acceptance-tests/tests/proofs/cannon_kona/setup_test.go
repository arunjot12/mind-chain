package cannon_kona

import (
	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
)

func cannonOpts() []presets.Option {
	return []presets.Option{
		presets.WithGameTypeAdded(gameTypes.CannonKonaGameType),
		presets.WithDeployerOptions(sysgo.WithJovianAtGenesis),
		presets.WithGlobalL2CLOption(sysgo.L2CLOptionFn(func(p devtest.T, _ sysgo.ComponentTarget, cfg *sysgo.L2CLConfig) {
			cfg.SafeDBPath = p.TempDir()
		})),
	}
}

func newSystem(t devtest.T) *presets.Minimal {
	return presets.NewMinimal(t, cannonOpts()...)
}
