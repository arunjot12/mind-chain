package standard

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/superchain"

	"github.com/mindchain/mind-chain/mc-chain-ops/genesis"
	"github.com/mindchain/mind-chain/mc-node/rollup"
	op_service "github.com/mindchain/mind-chain/mc-service"

	"github.com/ethereum/go-ethereum/common"
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
	MIPSVersion                     uint64 = 8
	DisputeGameType                 uint32 = 1 // PERMISSIONED game type
	DisputeMaxGameDepth             uint64 = 73
	DisputeSplitDepth               uint64 = 30
	DisputeClockExtension           uint64 = 10800
	DisputeMaxClockDuration         uint64 = 302400
	Eip1559DenominatorCanyon        uint64 = 250
	Eip1559Denominator              uint64 = 50
	Eip1559Elasticity               uint64 = 6

	ContractsV160Tag        = "op-contracts/v1.6.0"
	ContractsV180Tag        = "op-contracts/v1.8.0-rc.4"
	ContractsV170Beta1L2Tag = "op-contracts/v1.7.0-beta.1+l2-contracts"
	ContractsV200Tag        = "op-contracts/v2.0.0"
	ContractsV300Tag        = "op-contracts/v3.0.0"
	ContractsV400Tag        = "op-contracts/v4.0.0-rc.7"
	CurrentTag              = ContractsV400Tag
)

var DisputeAbsolutePrestate = common.HexToHash("0x038512e02c4c3f7bdaec27d00edf55b7155e0905301e1a88083e4e0a6764d54c")

var VaultMinWithdrawalAmount = mustHexBigFromHex("0x8ac7230489e80000")

var GovernanceTokenOwner = common.HexToAddress("0xDeaDDEaDDeAdDeAdDEAdDEaddeAddEAdDEAdDEad")

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

// DefaultHardforkScheduleForTag is used to determine which hardforks should be activated by default given a
// contract tag. For example, passing in v1.6.0 will return all hardforks up to and including Granite. This allows
// OP Deployer to set sane defaults for hardforks. This is not an ideal solution, but it will have to work until we get
// to MCP L2.
func DefaultHardforkScheduleForTag(tag string) *genesis.UpgradeScheduleDeployConfig {
	sched := &genesis.UpgradeScheduleDeployConfig{
		L2GenesisRegolithTimeOffset: op_service.U64UtilPtr(0),
		L2GenesisCanyonTimeOffset:   op_service.U64UtilPtr(0),
		L2GenesisDeltaTimeOffset:    op_service.U64UtilPtr(0),
		L2GenesisEcotoneTimeOffset:  op_service.U64UtilPtr(0),
		L2GenesisFjordTimeOffset:    op_service.U64UtilPtr(0),
		L2GenesisGraniteTimeOffset:  op_service.U64UtilPtr(0),
	}

	switch tag {
	case ContractsV160Tag, ContractsV170Beta1L2Tag:
		return sched
	case ContractsV180Tag, ContractsV200Tag, ContractsV300Tag:
		sched.ActivateForkAtGenesis(rollup.Holocene)
	case ContractsV400Tag:
		sched.ActivateForkAtGenesis(rollup.Holocene)
		sched.ActivateForkAtGenesis(rollup.Isthmus)
	default:
		sched.ActivateForkAtGenesis(rollup.Holocene)
		sched.ActivateForkAtGenesis(rollup.Isthmus)
	}

	return sched
}

func mustHexBigFromHex(hex string) *hexutil.Big {
	num := hexutil.MustDecodeBig(hex)
	hexBig := hexutil.Big(*num)
	return &hexBig
}

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

func OPCMImplAddressFor(chainID uint64, tag string) (common.Address, error) {
	return common.Address{}, fmt.Errorf("OPCM for tag %s: %w", tag, noStandardSuperchain(chainID))
}

func L1ProxyAdminOwner(chainID uint64) (common.Address, error) {
	return common.Address{}, noStandardSuperchain(chainID)
}

func L2ProxyAdminOwner(chainID uint64) (common.Address, error) {
	return common.Address{}, noStandardSuperchain(chainID)
}

func ProtocolVersionsOwner(chainID uint64) (common.Address, error) {
	return common.Address{}, noStandardSuperchain(chainID)
}
