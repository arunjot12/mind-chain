package bpo2

import (
	"testing"

	"github.com/ethereum/go-ethereum/params/forks"
	"github.com/mindchain/mind-chain/mc-acceptance-tests/tests/fusaka"
	"github.com/mindchain/mind-chain/mc-acceptance-tests/tests/jovian/bpo2/joviantest"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
)

func setupBPO2(t devtest.T) *presets.Minimal {
	return presets.NewMinimal(t,
		fusaka.L1GethOption(),
		presets.WithDeployerOptions(
			sysgo.WithJovianAtGenesis,
			sysgo.WithDefaultBPOBlobSchedule,
			sysgo.WithForkAtL1Genesis(forks.BPO2),
		),
	)
}

func TestDAFootprint(gt *testing.T) {
	joviantest.RunDAFootprint(gt, setupBPO2)
}

func TestMinBaseFee(gt *testing.T) {
	joviantest.RunMinBaseFee(gt, setupBPO2)
}

func TestOperatorFee(gt *testing.T) {
	joviantest.RunOperatorFee(gt, setupBPO2)
}
