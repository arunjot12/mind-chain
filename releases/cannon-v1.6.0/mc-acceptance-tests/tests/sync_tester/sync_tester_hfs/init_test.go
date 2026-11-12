package sync_tester_hfs

import (
	"testing"

	bss "github.com/mindchain/mind-chain/mc-batcher/batcher"
	"github.com/mindchain/mind-chain/mc-devstack/compat"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
	"github.com/mindchain/mind-chain/mc-node/rollup"
	"github.com/mindchain/mind-chain/mc-node/rollup/derive"
)

func TestMain(m *testing.M) {
	presets.DoMain(m, presets.WithSimpleWithSyncTester(),
		presets.WithCompatibleTypes(compat.SysGo),
		presets.WithHardforkSequentialActivation(rollup.Bedrock, rollup.Jovian, 15),
		stack.MakeCommon(sysgo.WithBatcherOption(func(id stack.L2BatcherID, cfg *bss.CLIConfig) {
			// For supporting pre-delta batches
			cfg.BatchType = derive.SingularBatchType
			// For supporting pre-Fjord batches
			cfg.CompressionAlgo = derive.Zlib
		})))
}
