package zk

import (
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
	"github.com/mindchain/mind-chain/mc-core/devfeatures"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-devstack/dsl/proofs"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

func newNonInteropZKSystem(t devtest.T) (*presets.Minimal, *proofs.DisputeGameFactory) {
	sysgo.SkipOnKonaNode(t, "requires mc-node's single-chain superroot_atTimestamp RPC")
	sys := presets.NewMinimal(t,
		presets.WithDeployerOptions(
			sysgo.WithJovianAtGenesis,
			sysgo.WithDevFeatureEnabled(devfeatures.ZKDisputeGameFlag),
		),
		presets.WithGameTypeAdded(gameTypes.ZKDisputeGameType),
		presets.WithRespectedGameTypeOverride(gameTypes.ZKDisputeGameType),
		presets.RequireRespectedGameType(gameTypes.ZKDisputeGameType),
	)
	t.Require().Len(sys.L2Networks(), 1)
	t.Require().Nil(sys.L2Chain.Escape().RollupConfig().LagoonTime, "interop must remain disabled")

	// Minimal runs no supernode; the proposer and challenger source super roots from mc-node.
	factory := proofs.NewDisputeGameFactory(t, sys.L1Network, sys.L1EL.EthClient(),
		sys.L2Chain.DisputeGameFactoryProxyAddr(), sys.L2CL, sys.L2EL, dsl.NewOpNodeSuperRoots(sys.L2CL), nil, nil)
	return sys, factory
}

func TestNonInteropZKProposerCreatesValidProposal(gt *testing.T) {
	t := devtest.ParallelT(gt)
	sys, factory := newNonInteropZKSystem(t)
	_, anchorSequence := sys.AnchorStateRegistry().AnchorRoot()

	game := factory.WaitForZKGameAtIndex(0)
	t.Require().Equal(uint32(math.MaxUint32), game.ParentIndex())
	t.Require().Greater(game.L2SequenceNumber(), anchorSequence)
	t.Require().Equal(proofs.ZKProposalUnchallenged, game.ProposalStatus())
	response := dsl.NewOpNodeSuperRoots(sys.L2CL).SuperRootAtTimestamp(game.L2SequenceNumber())
	t.Require().Equal([]eth.ChainID{sys.L2Chain.ChainID()}, response.ChainIDs)
	t.Require().NotNil(response.Data)
	t.Require().Equal(common.Hash(eth.SuperRoot(response.Data.Super)), game.RootClaimValue())
}

func TestNonInteropZKChallengerChallengesInvalidProposal(gt *testing.T) {
	t := devtest.ParallelT(gt)
	sys, factory := newNonInteropZKSystem(t)
	factory.WaitForZKGameAtIndex(0)

	_, anchorSequence := sys.AnchorStateRegistry().AnchorRoot()
	timestamp, outputRoots := factory.WaitForSafeSuperRootAfter(anchorSequence)
	t.Require().Len(outputRoots, 1)
	outputRoots[0][0] ^= 0xff
	invalidProposer := sys.FunderL1.NewFundedEOA(eth.OneEther)
	game := factory.StartZKGame(invalidProposer,
		proofs.WithL2SequenceNumber(timestamp),
		proofs.WithSuperRootFrom(outputRoots...),
	)

	game.WaitForProposalStatus(proofs.ZKProposalChallenged)
	t.Require().Equal(zkChallengerAddress(t, sys.L2Chain.ChainID()), game.ClaimData().Challenger)
}
