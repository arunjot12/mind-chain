package opcm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mindchain/mind-chain/mc-chain-ops/script"
	"github.com/mindchain/mind-chain/mc-chain-ops/script/addresses"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/broadcaster"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/standard"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/testutil"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/env"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/testlog"
	"github.com/stretchr/testify/require"
)

func TestDeployDisputeGame(t *testing.T) {
	t.Parallel()

	_, artifacts := testutil.LocalArtifacts(t)

	host, err := env.DefaultScriptHost(
		broadcaster.NoopBroadcaster(),
		testlog.Logger(t, log.LevelInfo),
		common.Address{'D'},
		artifacts,
	)
	require.NoError(t, err)

	vmAddr := deployDisputeGameScriptVM(t, host)

	input := DeployDisputeGameInput{
		Release:                  "dev",
		VmAddress:                vmAddr,
		GameKind:                 "PermissionedDisputeGame",
		GameType:                 1,
		AbsolutePrestate:         common.Hash{'A'},
		MaxGameDepth:             big.NewInt(int64(standard.DisputeMaxGameDepth)),
		SplitDepth:               big.NewInt(int64(standard.DisputeSplitDepth)),
		ClockExtension:           standard.DisputeClockExtension,
		MaxClockDuration:         standard.DisputeMaxClockDuration,
		DelayedWethProxy:         common.Address{'D'},
		AnchorStateRegistryProxy: common.Address{'A'},
		L2ChainId:                big.NewInt(69),
		Proposer:                 common.Address{'P'},
		Challenger:               common.Address{'C'},
	}

	script, err := NewDeployDisputeGameScript(host)
	require.NoError(t, err)

	output, err := script.Run(input)
	require.NoError(t, err)

	require.NotEmpty(t, output.DisputeGameImpl)
	require.NotEmpty(t, host.GetCode(output.DisputeGameImpl))
}

func deployDisputeGameScriptVM(t *testing.T, host *script.Host) common.Address {
	preimageOracleArtifact, err := host.Artifacts().ReadArtifact("PreimageOracle.sol", "PreimageOracle")
	require.NoError(t, err)

	encodedPreimageOracleConstructor, err := preimageOracleArtifact.ABI.Pack("", big.NewInt(0), big.NewInt(0))
	require.NoError(t, err)

	preimageOracleAddress, err := host.Create(addresses.ScriptDeployer, append(preimageOracleArtifact.Bytecode.Object, encodedPreimageOracleConstructor...))
	require.NoError(t, err)

	bigStepperArtifact, err := host.Artifacts().ReadArtifact("MIPS64.sol", "MIPS64")
	require.NoError(t, err)

	encodedBigStepperConstructor, err := bigStepperArtifact.ABI.Pack("", preimageOracleAddress, new(big.Int).SetUint64(standard.MIPSVersion))
	require.NoError(t, err)

	bigStepperAddress, err := host.Create(addresses.ScriptDeployer, append(bigStepperArtifact.Bytecode.Object, encodedBigStepperConstructor...))
	require.NoError(t, err)

	host.Label(preimageOracleAddress, "PreimageOracle")
	host.Label(bigStepperAddress, "BigStepper")

	return bigStepperAddress
}
