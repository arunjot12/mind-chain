package flags

import (
	"flag"

	"github.com/mindchain/mind-chain/mc-service/log/logcli"
)

type TestConfig struct {
	LogConfig logcli.CLIConfig
}

func ReadTestConfig() TestConfig {
	flag.Parse()

	cfg := logcli.ReadTestCLIConfig()

	return TestConfig{
		LogConfig: cfg,
	}
}
