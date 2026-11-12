package withdrawal

import (
	"testing"

	faultTypes "github.com/mindchain/mind-chain/mc-challenger/game/fault/types"
	"github.com/mindchain/mind-chain/mc-devstack/compat"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

func TestMain(m *testing.M) {
	presets.DoMain(m,
		// TODO(infra#401): Re-enable the test when the sysext missing toolset is implemented
		presets.WithCompatibleTypes(compat.SysGo),
		presets.WithMinimal(),
		presets.WithProposerGameType(faultTypes.FastGameType),
		// Fast game for test
		presets.WithFastGame(),
		// Deployer must be L1PAO to make mc-deployer happy
		presets.WithDeployerMatchL1PAO(),
		// Guardian must be L1PAO to make AnchorStateRegistry's setRespectedGameType method work
		presets.WithGuardianMatchL1PAO(),
		// Fast finalization for fast withdrawal
		presets.WithFinalizationPeriodSeconds(1),
		// Satisfy OptimismPortal2 PROOF_MATURITY_DELAY_SECONDS check, avoid OptimismPortal_ProofNotOldEnough() revert
		presets.WithProofMaturityDelaySeconds(2),
		// Satisfy AnchorStateRegistry DISPUTE_GAME_FINALITY_DELAY_SECONDS check, avoid OptimismPortal_InvalidRootClaim() revert
		presets.WithDisputeGameFinalityDelaySeconds(2),
	)
}
