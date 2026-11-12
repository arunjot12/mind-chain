package withdrawal

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-acceptance-tests/tests/base/withdrawal"
	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
)

func TestWithdrawal_Permissioned(gt *testing.T) {
	withdrawal.TestWithdrawal(gt, gameTypes.PermissionedGameType)
}
