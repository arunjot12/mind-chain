package presets

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

func WithSyncTesterELInitialState(fcu eth.FCUState) stack.CommonOption {
	return stack.MakeCommon(
		sysgo.WithGlobalSyncTesterELOption(sysgo.SyncTesterELOptionFn(
			func(_ devtest.P, id stack.L2ELNodeID, cfg *sysgo.SyncTesterELConfig) {
				cfg.FCUState = fcu
			})))
}

func WithELSyncTarget(elSyncTarget uint64) stack.CommonOption {
	return stack.MakeCommon(
		sysgo.WithGlobalSyncTesterELOption(sysgo.SyncTesterELOptionFn(
			func(_ devtest.P, id stack.L2ELNodeID, cfg *sysgo.SyncTesterELConfig) {
				cfg.ELSyncActive = true
				cfg.ELSyncTarget = elSyncTarget
			})))
}
