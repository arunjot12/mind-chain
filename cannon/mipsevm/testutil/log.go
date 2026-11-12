package testutil

import (
	"os"

	"github.com/mindchain/mind-chain/mc-service/log"
)

func CreateLogger() log.Logger {
	return log.NewLogger(log.LogfmtHandlerWithLevel(os.Stdout, log.LevelInfo))
}
