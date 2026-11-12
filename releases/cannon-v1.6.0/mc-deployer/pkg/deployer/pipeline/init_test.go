package pipeline

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/mindchain/mind-chain/mc-chain-ops/addresses"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/broadcaster"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/testutil"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/env"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/artifacts"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/standard"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/state"
	"github.com/mindchain/mind-chain/mc-service/testlog"
	"github.com/mindchain/mind-chain/mc-service/testutils/devnet"
	"github.com/stretchr/testify/require"
)

func TestInitLiveStrategy_OPCMReuseLogicSepolia(t *testing.T) {
	t.Parallel()

	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	require.NotEmpty(t, rpcURL, "SEPOLIA_RPC_URL must be set")

	lgr := testlog.Logger(t, slog.LevelInfo)
	retryProxy := devnet.NewRetryProxy(lgr, rpcURL)
	require.NoError(t, retryProxy.Start())
	t.Cleanup(func() {
		require.NoError(t, retryProxy.Stop())
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rpcClient, err := rpc.Dial(retryProxy.Endpoint())
	require.NoError(t, err)
	client := ethclient.NewClient(rpcClient)

	l1ChainID := uint64(11155111)
	t.Run("untagged L1 locator", func(t *testing.T) {
		st := &state.State{
			Version: 1,
		}
		require.NoError(t, InitLiveStrategy(
			ctx,
			&Env{
				L1Client: client,
				Logger:   lgr,
			},
			&state.Intent{
				L1ChainID:          l1ChainID,
				L1ContractsLocator: artifacts.MustNewLocatorFromURL("file:///not-a-path"),
				L2ContractsLocator: artifacts.MustNewLocatorFromURL("file:///not-a-path"),
			},
			st,
		))

		// Defining a file locator will always deploy a new superchain and OPCM
		require.Nil(t, st.SuperchainDeployment)
		require.Nil(t, st.ImplementationsDeployment)
	})

	t.Run("embedded L1 locator with standard intent types and standard roles", func(t *testing.T) {
		runTest := func(configType state.IntentType) {
			_, afacts := testutil.LocalArtifacts(t)
			host, err := env.DefaultForkedScriptHost(
				ctx,
				broadcaster.NoopBroadcaster(),
				testlog.Logger(t, log.LevelInfo),
				common.Address{'D'},
				afacts,
				rpcClient,
			)
			require.NoError(t, err)

			stdSuperchainRoles, err := state.GetStandardSuperchainRoles(l1ChainID)
			require.NoError(t, err)

			opcmAddr, err := standard.OPCMImplAddressFor(l1ChainID, standard.CurrentTag)
			require.NoError(t, err)

			intent := &state.Intent{
				ConfigType:         configType,
				L1ChainID:          l1ChainID,
				L1ContractsLocator: artifacts.EmbeddedLocator,
				L2ContractsLocator: artifacts.EmbeddedLocator,
				OPCMAddress:        &opcmAddr,
			}
			st := &state.State{
				Version: 1,
			}
			require.NoError(t, InitLiveStrategy(
				ctx,
				&Env{
					L1Client:     client,
					Logger:       lgr,
					L1ScriptHost: host,
				},
				intent,
				st,
			))

			// Defining a file locator will always deploy a new superchain and OPCM
			superCfg, err := standard.SuperchainFor(l1ChainID)
			require.NoError(t, err)
			proxyAdmin, err := standard.SuperchainProxyAdminAddrFor(l1ChainID)
			require.NoError(t, err)

			expDeployment := &addresses.SuperchainContracts{
				SuperchainProxyAdminImpl: proxyAdmin,
				ProtocolVersionsProxy:    superCfg.ProtocolVersionsAddr,
				ProtocolVersionsImpl:     common.HexToAddress("0x37E15e4d6DFFa9e5E320Ee1eC036922E563CB76C"),
				SuperchainConfigProxy:    superCfg.SuperchainConfigAddr,
				SuperchainConfigImpl:     common.HexToAddress("0xCe28685EB204186b557133766eCA00334EB441E4"),
			}

			// Tagged locator will reuse the existing superchain and OPCM
			require.NotNil(t, st.SuperchainDeployment)
			require.NotNil(t, st.ImplementationsDeployment)
			require.NotNil(t, st.SuperchainRoles)
			require.Equal(t, *expDeployment, *st.SuperchainDeployment)
			require.Equal(t, opcmAddr, st.ImplementationsDeployment.OpcmImpl)
			require.Equal(t, *stdSuperchainRoles, *st.SuperchainRoles)
		}

		runTest(state.IntentTypeStandard)
		runTest(state.IntentTypeStandardOverrides)
	})

	t.Run("tagged L1 locator with standard intent types and modified roles", func(t *testing.T) {
		runTest := func(configType state.IntentType) {
			intent := &state.Intent{
				ConfigType:         configType,
				L1ChainID:          l1ChainID,
				L1ContractsLocator: artifacts.DefaultL1ContractsLocator,
				L2ContractsLocator: artifacts.DefaultL2ContractsLocator,
				SuperchainRoles: &addresses.SuperchainRoles{
					SuperchainGuardian: common.Address{0: 99},
				},
			}
			st := &state.State{
				Version: 1,
			}
			require.NoError(t, InitLiveStrategy(
				ctx,
				&Env{
					L1Client: client,
					Logger:   lgr,
				},
				intent,
				st,
			))

			// Modified roles will cause a new superchain and OPCM to be deployed
			require.Nil(t, st.SuperchainDeployment)
			require.Nil(t, st.ImplementationsDeployment)
		}

		runTest(state.IntentTypeStandard)
		runTest(state.IntentTypeStandardOverrides)
	})

	t.Run("tagged locator with custom intent type", func(t *testing.T) {
		intent := &state.Intent{
			ConfigType:         state.IntentTypeCustom,
			L1ChainID:          l1ChainID,
			L1ContractsLocator: artifacts.DefaultL1ContractsLocator,
			L2ContractsLocator: artifacts.DefaultL2ContractsLocator,
			SuperchainRoles: &addresses.SuperchainRoles{
				SuperchainGuardian: common.Address{0: 99},
			},
		}
		st := &state.State{
			Version: 1,
		}
		require.NoError(t, InitLiveStrategy(
			ctx,
			&Env{
				L1Client: client,
				Logger:   lgr,
			},
			intent,
			st,
		))

		// Custom intent types always deploy a new superchain and OPCM
		require.Nil(t, st.SuperchainDeployment)
		require.Nil(t, st.ImplementationsDeployment)
	})
}
