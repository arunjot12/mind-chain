package presets

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
)

type MinimalWithConductors struct {
	*Minimal

	Conductors dsl.ConductorSet
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
