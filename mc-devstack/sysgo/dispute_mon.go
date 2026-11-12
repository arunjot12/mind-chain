package sysgo

import (
	"context"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	disputeMon "github.com/mindchain/mind-chain/mc-dispute-mon"
	disputeMonConfig "github.com/mindchain/mind-chain/mc-dispute-mon/config"
	disputeMonService "github.com/mindchain/mind-chain/mc-dispute-mon/mon"
	"github.com/mindchain/mind-chain/mc-service/clock"
	opmetrics "github.com/mindchain/mind-chain/mc-service/metrics"
)

type DisputeMonConfig struct {
	L1RPC              string
	GameFactoryAddress common.Address
	RollupRPCs         []string
	SupernodeRPCs      []string
	HonestActors       []common.Address
	Clock              clock.Clock
}

type DisputeMonRuntime struct {
	metricsURL string
}

func (r *DisputeMonRuntime) MetricsURL() string {
	return r.metricsURL
}

func StartDisputeMon(t devtest.T, input DisputeMonConfig) *DisputeMonRuntime {

	require := t.Require()
	require.NotEmpty(input.L1RPC, "L1 RPC is required")
	require.NotEqual(common.Address{}, input.GameFactoryAddress, "dispute game factory address is required")
	require.NotEmpty(append(append([]string{}, input.RollupRPCs...), input.SupernodeRPCs...), "at least one rollup or supernode RPC is required")

	cfg := disputeMonConfig.NewCombinedConfig(
		input.GameFactoryAddress,
		input.L1RPC,
		input.RollupRPCs,
		input.SupernodeRPCs,
	)
	cfg.HonestActors = append([]common.Address(nil), input.HonestActors...)
	cfg.MonitorInterval = time.Second
	cfg.MetricsConfig = opmetrics.CLIConfig{
		Enabled:    true,
		ListenAddr: "127.0.0.1",
		ListenPort: 0,
	}

	var serviceOptions []disputeMonService.ServiceOption
	if input.Clock != nil {
		serviceOptions = append(serviceOptions, disputeMonService.WithClock(input.Clock))
	}

	logger := t.Logger().New("component", "mc-dispute-mon")
	service, err := disputeMon.Main(t.Ctx(), logger, &cfg, serviceOptions...)
	require.NoError(err, "create dispute monitor service")
	metricsAddr := service.MetricsAddr()
	require.NotNil(metricsAddr, "dispute monitor metrics address")
	metricsURL := "http://" + metricsAddr.String()
	require.NoError(service.Start(t.Ctx()), "start dispute monitor service")
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		require.NoError(service.Stop(ctx), "stop dispute monitor service")
	})
	waitTCPReady(t, metricsURL, 10*time.Second)
	return &DisputeMonRuntime{metricsURL: metricsURL}
}
