package safeheaddb_elsync

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-devstack/compat"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

func TestMain(m *testing.M) {
	presets.DoMain(m, presets.WithSingleChainMultiNode(),
		presets.WithExecutionLayerSyncOnVerifiers(),
		presets.WithSafeDBEnabled(),
		// Destructive test that requiring an in-memory only geth database
		presets.WithCompatibleTypes(compat.SysGo),
	)
}
