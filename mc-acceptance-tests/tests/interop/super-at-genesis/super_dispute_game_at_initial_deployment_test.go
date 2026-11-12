package super_at_genesis

import (
	"testing"

	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

// TestSuperPermissionedDisputeGameInstalledAtInitialDeployment verifies that
// a chain deployed with the SuperRootGamesMigration dev feature flag has
// SuperPermissionedDisputeGame in the permissioned slot at initial deploy,
// without requiring a post-deploy OPCMv2 migration.
func TestSuperPermissionedDisputeGameInstalledAtInitialDeployment(gt *testing.T) {
	t := devtest.SerialT(gt)
	sys := presets.NewSingleChainInteropSuperRootAtGenesis(t)

	bridge := sys.StandardBridge(sys.L2ChainA)
	bridge.VerifyRespectedGameType(gameTypes.SuperPermissionedGameType)
	t.Require().Zero(bridge.GameResolutionDelay(), "SuperPermissioned games resolve during initialization")
	sys.DisputeGameFactory().VerifyGameImplPresent(gameTypes.SuperPermissionedGameType)
	sys.DisputeGameFactory().VerifyGameImplAbsent(gameTypes.PermissionedGameType)
}
