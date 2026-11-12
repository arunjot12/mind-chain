package presets

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
	"github.com/mindchain/mind-chain/mc-node/rollup/sync"
)

func WithExecutionLayerSyncOnVerifiers() stack.CommonOption {
	return stack.MakeCommon(
		sysgo.WithGlobalL2CLOption(sysgo.L2CLOptionFn(
			func(_ devtest.P, id stack.L2CLNodeID, cfg *sysgo.L2CLConfig) {
				cfg.VerifierSyncMode = sync.ELSync
			})))
}

func WithConsensusLayerSync() stack.CommonOption {
	return stack.MakeCommon(
		sysgo.WithGlobalL2CLOption(sysgo.L2CLOptionFn(
			func(_ devtest.P, id stack.L2CLNodeID, cfg *sysgo.L2CLConfig) {
				cfg.SequencerSyncMode = sync.CLSync
				cfg.VerifierSyncMode = sync.CLSync
			})))
}

func WithSafeDBEnabled() stack.CommonOption {
	return stack.MakeCommon(
		sysgo.WithGlobalL2CLOption(sysgo.L2CLOptionFn(
			func(p devtest.P, id stack.L2CLNodeID, cfg *sysgo.L2CLConfig) {
				cfg.SafeDBPath = p.TempDir()
			})))
}
