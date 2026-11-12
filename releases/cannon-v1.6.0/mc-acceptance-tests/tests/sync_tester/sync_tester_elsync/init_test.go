package sync_tester_elsync

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-devstack/compat"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

func TestMain(m *testing.M) {
	presets.DoMain(m,
		presets.WithExecutionLayerSyncOnVerifiers(),
		presets.WithSimpleWithSyncTester(),
		presets.WithELSyncTarget(35),
		presets.WithCompatibleTypes(compat.SysGo),
	)
}
