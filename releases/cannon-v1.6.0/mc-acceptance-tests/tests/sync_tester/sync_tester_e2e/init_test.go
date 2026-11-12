package sync_tester_e2e

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-devstack/compat"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

func TestMain(m *testing.M) {
	presets.DoMain(m, presets.WithSimpleWithSyncTester(),
		presets.WithCompatibleTypes(compat.SysGo),
	)
}
