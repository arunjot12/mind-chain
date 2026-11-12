package proofs_singlechain

import (
	"os"
	"testing"

	"github.com/mindchain/mind-chain/mc-acceptance-tests/tests/sdm/sdmtest"
	sfp "github.com/mindchain/mind-chain/mc-acceptance-tests/tests/superfaultproofs"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

func TestInteropSingleChainFaultProofs(gt *testing.T) {
	t := devtest.SerialT(gt)
	sys := presets.NewSingleChainInterop(t)
	sfp.RunSingleChainSuperFaultProofSmokeTest(t, sys, proofRunners()...)
}

func TestInteropSingleChainFaultProofsWithSDM(gt *testing.T) {
	t := devtest.SerialT(gt)
	sys := sdmtest.NewFixtureSingleChainFaultProofSystem(t)
	sfp.RunSingleChainSuperFaultProofSDMSmokeTest(t, sys)
}

func proofRunners() []sfp.ProofRunner {
	runners := []sfp.ProofRunner{sfp.NewKonaProofRunner()}
	if os.Getenv("RUST_BINARY_PATH_KONA_SP1_SUPER_RANGE_EXECUTOR") != "" || os.Getenv("RUST_JIT_BUILD") != "" {
		runners = append(runners, sfp.NewSP1NativeProofRunner())
	}
	return runners
}
