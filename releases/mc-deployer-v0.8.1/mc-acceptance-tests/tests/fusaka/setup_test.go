package fusaka

import (
	"github.com/ethereum/go-ethereum/params/forks"
	"github.com/mindchain/mind-chain/mc-batcher/batcher"
	"github.com/mindchain/mind-chain/mc-batcher/flags"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
)

func newMinimalFusaka(t devtest.T) *presets.Minimal {
	return presets.NewMinimal(t,
		L1GethOption(),
		presets.WithDeployerOptions(
			sysgo.WithDefaultBPOBlobSchedule,
			// Make the BPO fork happen after Osaka so we can easily use geth's eip4844.CalcBlobFee
			// to calculate the blob base fee using the Osaka parameters.
			sysgo.WithForkAtL1Genesis(forks.Osaka),
			sysgo.WithForkAtL1Offset(forks.BPO1, 1),
		),
		presets.WithBatcherOption(func(_ sysgo.ComponentTarget, cfg *batcher.CLIConfig) {
			cfg.DataAvailabilityType = flags.BlobsType
			cfg.TxMgrConfig.CellProofTime = 0 // Force cell proofs to be used
		}),
	)
}
