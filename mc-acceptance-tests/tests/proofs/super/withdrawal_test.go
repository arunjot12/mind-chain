package super

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-acceptance-tests/tests/base/withdrawal"
	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
)

func TestSuperPermissionedWithdrawal(gt *testing.T) {
	withdrawal.TestWithdrawal(gt, gameTypes.SuperPermissionedGameType)
}

func TestSuperCannonKonaWithdrawal(gt *testing.T) {
	withdrawal.TestWithdrawal(gt, gameTypes.SuperCannonKonaGameType)
}

func TestZKDisputeGameWithdrawal(gt *testing.T) {
	withdrawal.TestWithdrawal(gt, gameTypes.ZKDisputeGameType)
}
