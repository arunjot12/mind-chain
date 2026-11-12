package standard

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/mindchain/mind-chain/mc-chain-ops/genesis"
	"github.com/mindchain/mind-chain/mc-core/forks"
	"github.com/mindchain/mind-chain/mc-core/superchain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	GasLimit                        uint64 = 60_000_000
	BasefeeScalar                   uint32 = 1368
	BlobBaseFeeScalar               uint32 = 801949
	WithdrawalDelaySeconds          uint64 = 302400
	MinProposalSizeBytes            uint64 = 126000
	ChallengePeriodSeconds          uint64 = 86400
	ProofMaturityDelaySeconds       uint64 = 604800
	DisputeGameFinalityDelaySeconds uint64 = 302400
	// Bounds baked into the OptimismPortal, AnchorStateRegistry and DelayedWETH implementations
	// for the per-chain withdrawal delays above.
	MinProofMaturityDelaySeconds       uint64 = 86400
	MaxProofMaturityDelaySeconds       uint64 = 604800
	MinDisputeGameFinalityDelaySeconds uint64 = 43200
	MaxDisputeGameFinalityDelaySeconds uint64 = 302400
	MinWithdrawalDelaySeconds          uint64 = 43200
	MaxWithdrawalDelaySeconds          uint64 = 604800
	MIPSVersion                        uint64 = 8
	// DisputeGameType is the SUPER_PERMISSIONED game type. DeployOPChain requires the initial game
	// type to match the OPCM's family, and SUPER_ROOT_GAMES_MIGRATION is enabled by default, so the
	// permissioned selector for a standard deploy is the super root one.
	// TODO(#21662): revisit with the broader SuperRootGamesMigration cleanup.
	DisputeGameType          uint32 = 5
	DisputeMaxGameDepth      uint64 = 73
	DisputeSplitDepth        uint64 = 30
	DisputeClockExtension    uint64 = 10800
	DisputeMaxClockDuration  uint64 = 302400
	Eip1559DenominatorCanyon uint64 = 250
	Eip1559Denominator       uint64 = 50
	Eip1559Elasticity        uint64 = 6

	ContractsV160Tag        = "op-contracts/v1.6.0"
	ContractsV180Tag        = "op-contracts/v1.8.0-rc.4"
	ContractsV170Beta1L2Tag = "op-contracts/v1.7.0-beta.1+l2-contracts"
	ContractsV200Tag        = "op-contracts/v2.0.0"
	ContractsV300Tag        = "op-contracts/v3.0.0"
	ContractsV400Tag        = "op-contracts/v4.0.0-rc.7"
	ContractsV410Tag        = "op-contracts/v4.1.0"
	ContractsV500Tag        = "op-contracts/v5.0.0"
	ContractsV600Tag        = "op-contracts/v6.0.0-rc.2"
	ContractsV700Tag        = "op-contracts/v7.0.0-rc.4"
	ContractsV800Tag        = "op-contracts/v8.0.0"
	CurrentTag              = ContractsV800Tag
)

var DisputeAbsolutePrestate = common.HexToHash("0x038512e02c4c3f7bdaec27d00edf55b7155e0905301e1a88083e4e0a6764d54c")

var VaultMinWithdrawalAmount = mustHexBigFromHex("0x8ac7230489e80000")

var GovernanceTokenOwner = common.HexToAddress("0xDeaDDEaDDeAdDeAdDEAdDEaddeAddEAdDEAdDEad")

// PlaceholderAddress is a non-zero sentinel address. It's used  as the default deployer when no private key is provided.
// The same-sender check keys off this value to skip when no real deployer is set.
var PlaceholderAddress = common.Address{0x01}

// ErrNoStandardSuperchain is returned by every lookup of the upstream "standard"
// Superchain values (Optimism's shared contracts, roles and contract-manager
// releases on Ethereum mainnet and Sepolia). Mind Chain does not ship the upstream
// superchain-registry, so standard intents are unsupported: use a custom intent
// with Mind Chain's own roles and contract addresses.
var ErrNoStandardSuperchain = errors.New("standard Superchain values are not available in Mind Chain; use a custom intent")

func noStandardSuperchain(chainID uint64) error {
	return fmt.Errorf("%w (L1 chain ID %d)", ErrNoStandardSuperchain, chainID)
}

func GuardianAddressFor(chainID uint64) (common.Address, error) {
	return common.Address{}, noStandardSuperchain(chainID)
}

func ChallengerAddressFor(chainID uint64) (common.Address, error) {
	return common.Address{}, noStandardSuperchain(chainID)
}

const (
	// Source: succinctlabs/sp1-contracts@2ac5ecbbe473421a963d67e55f182e9a36576f7c,
	// contracts/deployments/1.json, V6_1_0_SP1_VERIFIER_PLONK.
	mainnetSP1VerifierV610 = "0xc3c6dDDAc8829b233Dc6536Ec024775a57b0AF2A"
	// Succinct deploys the same verifier bytecode and address deterministically on both networks.
	// Source: succinctlabs/sp1-contracts@2ac5ecbbe473421a963d67e55f182e9a36576f7c,
	// contracts/deployments/11155111.json, V6_1_0_SP1_VERIFIER_PLONK.
	sepoliaSP1VerifierV610 = "0xc3c6dDDAc8829b233Dc6536Ec024775a57b0AF2A"
)

// sp1VerifierHashJSON maps each L1 chain ID to the VERIFIER_HASH() of the verifier from
// SP1VerifierFor, i.e. the PLONK circuit it accepts proofs for. Proofs carry the first four bytes
// of that value as their selector, so the linked sp1-sdk must prove for the same circuit;
// op-zk-proposer's release-pin test reads the same file and holds the SDK to it.
// Source: succinctlabs/sp1-contracts@2ac5ecbbe473421a963d67e55f182e9a36576f7c,
// contracts/src/v6.1.0/SP1VerifierPlonk.sol, VERIFIER_HASH().
// DO NOT MODIFY THIS FILE WITHOUT CLEARING IT WITH THE EVM SAFETY TEAM.
//
//go:embed sp1-verifier.json
var sp1VerifierHashJSON []byte

var sp1VerifierHashes = func() map[string]common.Hash {
	var hashes map[string]common.Hash
	if err := json.Unmarshal(sp1VerifierHashJSON, &hashes); err != nil {
		panic(fmt.Sprintf("sp1-verifier.json: %v", err))
	}
	return hashes
}()

// SP1VerifierFor returns the raw SP1 verifier approved for the current OPCM release on the given L1
// chain ID. Both `bootstrap implementations` and `apply` default to it when ZK dispute games are
// enabled and the operator did not pin a verifier explicitly. Change it together with
// SP1VerifierHashFor and the sp1-sdk pin in rust/Cargo.toml.
// DO NOT MODIFY THIS METHOD WITHOUT CLEARING IT WITH THE EVM SAFETY TEAM.
func SP1VerifierFor(chainID uint64) (common.Address, error) {
	switch chainID {
	case 1:
		return common.HexToAddress(mainnetSP1VerifierV610), nil
	case 11155111:
		return common.HexToAddress(sepoliaSP1VerifierV610), nil
	default:
		return common.Address{}, fmt.Errorf("unsupported chain ID: %d", chainID)
	}
}

// SP1VerifierHashFor returns the VERIFIER_HASH() the verifier from SP1VerifierFor implements on
// the given L1 chain ID, from sp1-verifier.json. TestApplyDefaultsSP1VerifierOnSepolia holds it to
// the chain, and the EVM Safety release runbook passes it to VerifyOPCM as
// EXPECTED_SP1_VERIFIER_HASH.
func SP1VerifierHashFor(chainID uint64) (common.Hash, error) {
	hash, ok := sp1VerifierHashes[strconv.FormatUint(chainID, 10)]
	if !ok {
		return common.Hash{}, fmt.Errorf("unsupported chain ID: %d", chainID)
	}
	return hash, nil
}

func SuperchainFor(chainID uint64) (superchain.Superchain, error) {
	switch chainID {
	case 1:
		return superchain.GetSuperchain("mainnet")
	case 11155111:
		return superchain.GetSuperchain("sepolia")
	default:
		return superchain.Superchain{}, fmt.Errorf("unsupported chain ID: %d", chainID)
	}
}

func OPCMImplAddressFor(chainID uint64, tag string) (common.Address, error) {
	return common.Address{}, fmt.Errorf("OPCM for tag %s: %w", tag, noStandardSuperchain(chainID))
}

// SuperchainProxyAdminAddrFor returns the address of the Superchain ProxyAdmin for the given chain ID.
// These have been verified to be the ProxyAdmin addresses on Mainnet and Sepolia.
// DO NOT MODIFY THIS METHOD WITHOUT CLEARING IT WITH THE EVM SAFETY TEAM.
func SuperchainProxyAdminAddrFor(chainID uint64) (common.Address, error) {
	switch chainID {
	case 1:
		return common.HexToAddress("0x543bA4AADBAb8f9025686Bd03993043599c6fB04"), nil
	case 11155111:
		return common.HexToAddress("0x189aBAAaa82DfC015A588A7dbaD6F13b1D3485Bc"), nil
	default:
		return common.Address{}, fmt.Errorf("unsupported chain ID: %d", chainID)
	}
}

func L1ProxyAdminOwner(chainID uint64) (common.Address, error) {
	return common.Address{}, noStandardSuperchain(chainID)
}

func L2ProxyAdminOwner(chainID uint64) (common.Address, error) {
	return common.Address{}, noStandardSuperchain(chainID)
}

// DefaultHardfork is the fork activated at genesis by default. It was previously
// read from OP Mainnet's schedule in the superchain-registry; it is pinned here
// because Mind Chain does not ship that registry. Raise it when adopting a newer fork.
const DefaultHardfork = forks.Karst

// DefaultHardforkSchedule is used to determine which hardforks should be activated by default.
// It activates DefaultHardfork, and every fork before it, at genesis.
func DefaultHardforkSchedule() *genesis.UpgradeScheduleDeployConfig {
	sched := &genesis.UpgradeScheduleDeployConfig{}
	sched.ActivateForkAtGenesis(DefaultHardfork)

	return sched
}

func mustHexBigFromHex(hex string) *hexutil.Big {
	num := hexutil.MustDecodeBig(hex)
	hexBig := hexutil.Big(*num)
	return &hexBig
}
