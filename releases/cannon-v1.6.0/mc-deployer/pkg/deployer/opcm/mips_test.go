package opcm

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/broadcaster"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/standard"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/testutil"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/env"
	"github.com/mindchain/mind-chain/mc-service/testlog"
	"github.com/stretchr/testify/require"
)

func TestDeployMIPS(t *testing.T) {
	t.Parallel()

	_, artifacts := testutil.LocalArtifacts(t)

	host, err := env.DefaultScriptHost(
		broadcaster.NoopBroadcaster(),
		testlog.Logger(t, log.LevelInfo),
		common.Address{'D'},
		artifacts,
	)
	require.NoError(t, err)

	input := DeployMIPSInput{
		MipsVersion:    uint64(standard.MIPSVersion),
		PreimageOracle: common.Address{0xab},
	}

	output, err := DeployMIPS(host, input)
	require.NoError(t, err)

	require.NotEmpty(t, output.MipsSingleton)
}
