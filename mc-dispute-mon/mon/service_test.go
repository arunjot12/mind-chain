package mon

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mindchain/mind-chain/mc-challenger/game/fault/contracts"
	"github.com/mindchain/mind-chain/mc-dispute-mon/config"
	"github.com/mindchain/mind-chain/mc-dispute-mon/metrics"
	"github.com/mindchain/mind-chain/mc-dispute-mon/mon/extract"
	monTypes "github.com/mindchain/mind-chain/mc-dispute-mon/mon/types"
	"github.com/mindchain/mind-chain/mc-service/clock"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/sources"
	"github.com/mindchain/mind-chain/mc-service/sources/batching"
	"github.com/mindchain/mind-chain/mc-service/sources/batching/rpcblock"
	batchingTest "github.com/mindchain/mind-chain/mc-service/sources/batching/test"
	"github.com/mindchain/mind-chain/mc-service/testlog"
	"github.com/mindchain/mind-chain/packages/contracts-bedrock/snapshots"
	"github.com/stretchr/testify/require"
)

func TestInitMonitorWiresEveryTypedLane(t *testing.T) {
	ctx := t.Context()
	factoryAddress := common.Address{0xaa}
	logger := testlog.Logger(t, log.LevelCrit)
	stubRPC := &serviceRPCStub{AbiBasedRpc: batchingTest.NewAbiBasedRpc(t, factoryAddress, snapshots.LoadDisputeGameFactoryABI())}
	stubRPC.SetResponse(factoryAddress, "version", rpcblock.Latest, nil, []interface{}{"1.4.0"})
	caller := batching.NewMultiCaller(stubRPC, batching.DefaultBatchSize)
	factory, err := contracts.NewDisputeGameFactoryContract(ctx, metrics.NoopMetrics, factoryAddress, caller)
	require.NoError(t, err)
	l1Client, err := sources.NewL1Client(stubRPC, logger, metrics.NoopMetrics, sources.L1ClientSimpleConfig(true, sources.RPCKindAny, 100))
	require.NoError(t, err)

	service := &Service{
		logger:          logger,
		metrics:         metrics.NoopMetrics,
		honestActors:    monTypes.NewHonestActors(nil),
		factoryContract: factory,
		cl:              clock.NewDeterministicClock(time.Unix(1, 0)),
		l1Client:        l1Client,
		l1Caller:        caller,
	}
	service.game = extract.NewGameCallerCreator(service.metrics, caller)
	cfg := config.NewCombinedConfig(factoryAddress, "unused", nil, nil)

	require.Equal(t, []string{
		"*extract.L1HeadBlockNumEnricher",
		"*extract.OutputAgreementEnricher",
		"*extract.SuperAgreementEnricher",
		"*extract.AnchorStateRegistryEnricher",
	}, registeredTypeNames(service.commonEnrichers()))
	require.Equal(t, []string{
		"*extract.ClaimEnricher",
		"*extract.BondDataEnricher",
	}, registeredTypeNames(service.faultEnrichers(extract.NewBondDataEnricher())))

	service.initMonitor(ctx, &cfg)
	require.NotNil(t, service.monitor)
	require.NotNil(t, service.monitor.fetchHeadBlock)
	require.NotNil(t, service.monitor.extract)
	require.NotNil(t, service.monitor.forecast)
	require.NotNil(t, service.monitor.checkAnchorState)
	expectedCommon := []string{
		"(*UpdateTimeMonitor).CheckUpdateTimes-fm",
		"(*NodeEndpointErrorsMonitor).CheckNodeEndpointErrors-fm",
		"(*NodeEndpointErrorCountMonitor).CheckNodeEndpointErrorCount-fm",
		"(*NodeEndpointOutOfSyncMonitor).CheckNodeEndpointOutOfSync-fm",
		"(*MixedAvailability).CheckMixedAvailability-fm",
		"(*MixedSafetyMonitor).CheckMixedSafety-fm",
		"(*DifferentRootMonitor).CheckDifferentRoots-fm",
		"(*GameTypeMonitor).CheckGameTypes-fm",
	}
	require.Len(t, service.monitor.commonMonitors, len(expectedCommon))
	for i, expected := range expectedCommon {
		require.Contains(t, registeredFunctionName(service.monitor.commonMonitors[i]), expected)
	}
	expectedFault := []string{
		"(*ResolutionMonitor).CheckResolutions-fm",
		"(*ClaimMonitor).CheckClaims-fm",
		"(*L2ChallengesMonitor).CheckL2Challenges-fm",
	}
	require.Len(t, service.monitor.faultMonitors, len(expectedFault))
	for i, expected := range expectedFault {
		require.Contains(t, registeredFunctionName(service.monitor.faultMonitors[i]), expected)
	}
	expectedBond := []string{
		"bonds.(*Bonds).CheckBonds-fm",
		"(*WithdrawalMonitor).CheckWithdrawals-fm",
	}
	require.Len(t, service.monitor.bondMonitors, len(expectedBond))
	for i, expected := range expectedBond {
		require.Contains(t, registeredFunctionName(service.monitor.bondMonitors[i]), expected)
	}
	expectedZK := []string{"(*ZKLifecycleMonitor).CheckLifecycle-fm"}
	require.Len(t, service.monitor.zkMonitors, len(expectedZK))
	for i, expected := range expectedZK {
		require.Contains(t, registeredFunctionName(service.monitor.zkMonitors[i]), expected)
	}
}

func registeredFunctionName(fn any) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func registeredTypeNames[T any](values []T) []string {
	names := make([]string, len(values))
	for i, value := range values {
		names[i] = reflect.TypeOf(value).String()
	}
	return names
}

type serviceRPCStub struct {
	*batchingTest.AbiBasedRpc
}

func (*serviceRPCStub) Close() {}

func (*serviceRPCStub) Subscribe(context.Context, string, any, ...any) (ethereum.Subscription, error) {
	return nil, errors.New("subscriptions are not supported by the test RPC")
}
