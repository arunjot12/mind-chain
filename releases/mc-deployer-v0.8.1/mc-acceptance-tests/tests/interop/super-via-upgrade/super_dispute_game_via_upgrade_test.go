package super_via_upgrade

import (
	"testing"

	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

// TestSuperRootGamesInstalledViaOPCMUpgrade verifies that opcm.upgrade installs
// the permissionless super games and makes SUPER_CANNON_KONA the respected type.
func TestSuperRootGamesInstalledViaOPCMUpgrade(gt *testing.T) {
	t := devtest.SerialT(gt)
	sys := presets.NewSingleChainInterop(t)

	sys.StandardBridge(sys.L2ChainA).VerifyRespectedGameType(gameTypes.SuperCannonKonaGameType)
	sys.DisputeGameFactory().VerifyGameImplPresent(gameTypes.SuperCannonKonaGameType)
	sys.DisputeGameFactory().VerifyGameImplPresent(gameTypes.SuperPermissionedGameType)
	sys.DisputeGameFactory().VerifyGameImplAbsent(gameTypes.PermissionedGameType)
}
