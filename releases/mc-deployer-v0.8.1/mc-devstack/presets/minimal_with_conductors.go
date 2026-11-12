package presets

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type MinimalWithConductors struct {
	*Minimal

	ConductorSets map[eth.ChainID]dsl.ConductorSet
}

// NewMinimalWithConductors creates a fresh MinimalWithConductors target for the current
// test.
//
// The target is created from the runtime plus any additional preset options.
func NewMinimalWithConductors(t devtest.T, opts ...Option) *MinimalWithConductors {
	presetCfg, presetOpts := collectSupportedPresetConfig(t, "NewMinimalWithConductors", opts, minimalWithConductorsPresetSupportedOptionKinds)
	out := minimalWithConductorsFromRuntime(t, sysgo.NewMinimalWithConductorsRuntimeWithConfig(t, presetCfg))
	presetOpts.applyPreset(out)
	return out
}
