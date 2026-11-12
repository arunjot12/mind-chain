package flags

import (
	"flag"

	"github.com/mindchain/mind-chain/mc-service/log"
)

type TestConfig struct {
	LogConfig log.CLIConfig
}

func ReadTestConfig() TestConfig {
	flag.Parse()

	cfg := log.ReadTestCLIConfig()

	return TestConfig{
		LogConfig: cfg,
	}
}
