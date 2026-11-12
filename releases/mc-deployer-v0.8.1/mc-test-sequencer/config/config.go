package config

import (
	"errors"

	oplog "github.com/mindchain/mind-chain/mc-service/log"
	opmetrics "github.com/mindchain/mind-chain/mc-service/metrics"
	"github.com/mindchain/mind-chain/mc-service/oppprof"
	oprpc "github.com/mindchain/mind-chain/mc-service/rpc"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/backend/work"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/backend/work/config"
)

const (
	DefaultConfigYaml = "config.yaml"
)

type Config struct {
	Version string

	LogConfig     oplog.CLIConfig
	MetricsConfig opmetrics.CLIConfig
	PprofConfig   oppprof.CLIConfig
	RPC           oprpc.CLIConfig

	JWTSecretPath string

	Ensemble work.Loader

	MockRun bool
}

func (c *Config) Check() error {
	var result error
	result = errors.Join(result, c.MetricsConfig.Check())
	result = errors.Join(result, c.PprofConfig.Check())
	result = errors.Join(result, c.RPC.Check())
	return result
}

func DefaultCLIConfig() *Config {
	return &Config{
		Version:       "dev",
		LogConfig:     oplog.DefaultCLIConfig(),
		MetricsConfig: opmetrics.DefaultCLIConfig(),
		PprofConfig:   oppprof.DefaultCLIConfig(),
		RPC:           oprpc.DefaultCLIConfig(),
		Ensemble:      &config.YamlLoader{Path: DefaultConfigYaml},
		MockRun:       false,
	}
}
