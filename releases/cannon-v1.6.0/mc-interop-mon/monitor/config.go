package monitor

import (
	"errors"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/mindchain/mind-chain/mc-interop-mon/flags"
	oplog "github.com/mindchain/mind-chain/mc-service/log"
	opmetrics "github.com/mindchain/mind-chain/mc-service/metrics"
	"github.com/mindchain/mind-chain/mc-service/oppprof"
	oprpc "github.com/mindchain/mind-chain/mc-service/rpc"
)

type CLIConfig struct {
	L2Rpcs              []string
	PollInterval        time.Duration
	SupervisorEndpoints []string
	TriggerFailsafe     bool
	RPCConfig           oprpc.CLIConfig
	LogConfig           oplog.CLIConfig
	MetricsConfig       opmetrics.CLIConfig
	PprofConfig         oppprof.CLIConfig
}

func (c *CLIConfig) Check() error {
	if err := c.RPCConfig.Check(); err != nil {
		return err
	}
	if err := c.MetricsConfig.Check(); err != nil {
		return err
	}
	if err := c.PprofConfig.Check(); err != nil {
		return err
	}

	if len(c.L2Rpcs) == 0 {
		return errors.New("l2 rpcs are required")
	}

	return nil
}

func NewConfig(ctx *cli.Context) *CLIConfig {
	return &CLIConfig{
		// Required Flags
		L2Rpcs: ctx.StringSlice(flags.L2RpcsFlag.Name),

		// Optional Flags
		SupervisorEndpoints: ctx.StringSlice(flags.SupervisorEndpointsFlag.Name),
		TriggerFailsafe:     ctx.Bool(flags.TriggerFailsafeFlag.Name),
		RPCConfig:           oprpc.ReadCLIConfig(ctx),
		LogConfig:           oplog.ReadCLIConfig(ctx),
		MetricsConfig:       opmetrics.ReadCLIConfig(ctx),
		PprofConfig:         oppprof.ReadCLIConfig(ctx),
	}
}
