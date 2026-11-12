package super_opnode

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-devstack/dsl/proofs"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-service/eth"
	safety "github.com/mindchain/mind-chain/mc-service/eth/safety"
)

// TestChallengerPlaysGameFromOpNode verifies mc-challenger plays a super-cannon-kona game
// when its super roots are sourced from a single mc-node (superroot_atTimestamp) rather than
// mc-supernode. An attacker posts a bad super root; the challenger, backed by mc-node, must
// counter it down to the split depth.
func TestChallengerPlaysGameFromOpNode(gt *testing.T) {
	t := devtest.ParallelT(gt)
	sys := presets.NewSingleChainInteropNoSupernode(t)
	dsl.CheckAll(t, sys.L2CLA.AdvancedFn(safety.CrossSafe, 1, 30))

	badClaim := common.HexToHash("0xdeadbeef00000000000000000000000000000000000000000000000000000000")
	attacker := sys.FunderL1.NewFundedEOA(eth.Ether(15))
	dgf := sys.DisputeGameFactory()

	game := dgf.StartSuperCannonKonaGame(attacker, proofs.WithSuperRootFrom(eth.Bytes32(badClaim)))

	claim := game.RootClaim()                   // The bad claim from the attacker.
	counterClaim := claim.WaitForCounterClaim() // The counter-claim from the challenger.
	for counterClaim.Depth() <= game.SplitDepth() {
		claim = counterClaim.Attack(attacker, badClaim)
		counterClaim = claim.WaitForCounterClaim()
	}
}
