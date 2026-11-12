package seqwindow

import (
	"testing"

	bss "github.com/mindchain/mind-chain/mc-batcher/batcher"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
	"github.com/mindchain/mind-chain/mc-node/rollup/derive"
)

func TestMain(m *testing.M) {
	presets.DoMain(m,
		presets.WithSimpleInterop(),
		// Short enough that we can run the test,
		// long enough that the batcher can still submit something before we make things expire.
		presets.WithSequencingWindow(10, 30),
		stack.MakeCommon(sysgo.WithBatcherOption(func(id stack.L2BatcherID, cfg *bss.CLIConfig) {
			// Span-batches during recovery don't appear to align well with the starting-point.
			// It can be off by ~6 L2 blocks, possibly due to off-by-one
			// in L1 block sync considerations in batcher stop or start.
			// So we end up having to encode block by block, so the full batch data does not get dropped.
			cfg.BatchType = derive.SingularBatchType
		})))
}
