package presets

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-devstack/shim"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
	"github.com/mindchain/mind-chain/mc-devstack/stack/match"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type MinimalWithSyncTester struct {
	Minimal

	SyncTester *dsl.SyncTester
}

func WithMinimalWithSyncTester(fcu eth.FCUState) stack.CommonOption {
	return stack.MakeCommon(sysgo.DefaultMinimalSystemWithSyncTester(&sysgo.DefaultMinimalSystemWithSyncTesterIDs{}, fcu))
}

func NewMinimalWithSyncTester(t devtest.T) *MinimalWithSyncTester {
	system := shim.NewSystem(t)
	orch := Orchestrator()
	orch.Hydrate(system)
	minimal := minimalFromSystem(t, system, orch)
	l2 := system.L2Network(match.Assume(t, match.L2ChainA))
	syncTester := l2.SyncTester(match.Assume(t, match.FirstSyncTester))
	return &MinimalWithSyncTester{
		Minimal:    *minimal,
		SyncTester: dsl.NewSyncTester(syncTester),
	}
}
