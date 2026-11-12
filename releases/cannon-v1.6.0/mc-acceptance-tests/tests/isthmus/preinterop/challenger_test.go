package preinterop

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-devstack/dsl/proofs"
	"github.com/mindchain/mind-chain/mc-supervisor/supervisor/types"

	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

func TestChallengerPlaysGame(gt *testing.T) {
	gt.Skip("TODO(#16166): Re-enable once the supervisor endpoint supports super roots before interop")
	t := devtest.ParallelT(gt)
	sys := presets.NewSimpleInterop(t)
	dsl.CheckAll(t,
		sys.L2CLA.AdvancedFn(types.CrossSafe, 1, 30),
		sys.L2CLB.AdvancedFn(types.CrossSafe, 1, 30),
	)

	badClaim := common.HexToHash("0xdeadbeef00000000000000000000000000000000000000000000000000000000")
	attacker := sys.FunderL1.NewFundedEOA(eth.Ether(15))
	dgf := sys.DisputeGameFactory()

	game := dgf.StartSuperCannonGame(attacker, proofs.WithRootClaim(badClaim))

	claim := game.RootClaim()                   // This is the bad claim from attacker
	counterClaim := claim.WaitForCounterClaim() // This is the counter-claim from the challenger
	for counterClaim.Depth() <= game.SplitDepth() {
		claim = counterClaim.Attack(attacker, badClaim)
		// Wait for the challenger to counter the attacker's claim, then attack again
		counterClaim = claim.WaitForCounterClaim()
	}
}
