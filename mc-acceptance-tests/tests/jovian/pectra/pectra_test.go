package pectra

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-acceptance-tests/tests/jovian/bpo2/joviantest"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
)

func setupPectra(t devtest.T) *presets.Minimal {
	return presets.NewMinimal(t,
		presets.WithDeployerOptions(sysgo.WithJovianAtGenesis),
	)
}

func TestDAFootprint(gt *testing.T) {
	joviantest.RunDAFootprint(gt, setupPectra)
}

func TestMinBaseFee(gt *testing.T) {
	joviantest.RunMinBaseFee(gt, setupPectra)
}

func TestOperatorFee(gt *testing.T) {
	joviantest.RunOperatorFee(gt, setupPectra)
}
