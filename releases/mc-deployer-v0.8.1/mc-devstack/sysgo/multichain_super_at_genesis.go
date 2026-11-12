package sysgo

import (
	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
	"github.com/mindchain/mind-chain/mc-core/devfeatures"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
)

func withSuperRootGamesAtGenesisDeployerFeatures(cfg PresetConfig) PresetConfig {
	cfg.DeployerOptions = append([]DeployerOption{
		WithDevFeatureEnabled(devfeatures.OptimismPortalInteropFlag),
		WithDevFeatureEnabled(devfeatures.SuperRootGamesMigrationFlag),
	}, cfg.DeployerOptions...)
	return cfg
}

// NewSingleChainSuperRootAtGenesisRuntimeWithConfig builds a single-chain
// supernode runtime with SUPER_PERMISSIONED installed in the
// permissioned slot at initial deploy.
func NewSingleChainSuperRootAtGenesisRuntimeWithConfig(t devtest.T, cfg PresetConfig) *MultiChainRuntime {
	cfg = withSuperRootGamesAtGenesisDeployerFeatures(cfg)
	runtime := newSingleChainSupernodeRuntimeWithConfig(t, true, cfg)
	attachTestSequencerToRuntime(t, runtime, "dev")
	attachSuperChallengerAndProposer(t, runtime, cfg, gameTypes.SuperPermissionedGameType)
	return runtime
}
