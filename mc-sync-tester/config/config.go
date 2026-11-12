package config

import (
	"errors"

	stconf "github.com/mindchain/mind-chain/mc-sync-tester/synctester/backend/config"

	"github.com/mindchain/mind-chain/mc-service/log/logcli"
	opmetrics "github.com/mindchain/mind-chain/mc-service/metrics"
	"github.com/mindchain/mind-chain/mc-service/oppprof"
	oprpc "github.com/mindchain/mind-chain/mc-service/rpc"
)

const (
	DefaultConfigYaml = "config.yaml"
)

type Config struct {
	Version string

	LogConfig     logcli.CLIConfig
	MetricsConfig opmetrics.CLIConfig
	PprofConfig   oppprof.CLIConfig
	RPC           oprpc.CLIConfig

	SyncTesters stconf.Loader
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
		LogConfig:     logcli.DefaultCLIConfig(),
		MetricsConfig: opmetrics.DefaultCLIConfig(),
		PprofConfig:   oppprof.DefaultCLIConfig(),
		RPC:           oprpc.DefaultCLIConfig(),
		SyncTesters:   &stconf.YamlLoader{Path: DefaultConfigYaml},
	}
}
