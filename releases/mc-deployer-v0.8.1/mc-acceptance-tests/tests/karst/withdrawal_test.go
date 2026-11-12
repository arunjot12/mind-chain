package karst

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-acceptance-tests/tests/base/withdrawal"
	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
	opforks "github.com/mindchain/mind-chain/mc-core/forks"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
)

// TestWithdrawal_Karst creates a withdrawal from the L2StandardBridge and
// observes the full withdrawal flow, including finalization on L1.
func TestWithdrawal_Karst(gt *testing.T) {
	withdrawal.TestWithdrawal(gt, gameTypes.CannonKonaGameType,
		presets.WithDeployerOptions(sysgo.WithKarstAtGenesis),
	)
}

// TestWithdrawal_KarstUpgrade is the same withdrawal flow but on a network that
// started pre-Karst and activated Karst mid-chain via a scheduled upgrade.
func TestWithdrawal_KarstUpgrade(gt *testing.T) {
	offset := uint64(10) // arbitrary offset to have a few blocks before Karst
	withdrawal.TestWithdrawalAfterUpgrade(gt, gameTypes.CannonKonaGameType, opforks.Karst,
		presets.WithDeployerOptions(sysgo.WithKarstAtOffset(&offset)),
	)
}
