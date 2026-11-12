package msg

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-supervisor/supervisor/types"
)

// TestUnscheduledInterop runs against an interop system (i.e. op-nodes are managed by mc-supervisor),
// before interop is scheduled with an actual hardfork time.
// And then confirms we can finalize the chains.
func TestUnscheduledInterop(gt *testing.T) {
	gt.Skip("Skipping Interop Acceptance Test")
	t := devtest.SerialT(gt)
	sys := presets.NewSimpleInterop(t)
	t.Logger().Info("Checking that chain A and B can sync, even though interop is not scheduled")
	dsl.CheckAll(t,
		sys.L2CLA.AdvancedFn(types.Finalized, 5, 100),
		sys.L2CLB.AdvancedFn(types.Finalized, 5, 100),
	)
	// Note: supervisor sync status won't be ready till after interop activation.
}
