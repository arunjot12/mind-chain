package logcli

import (
	"os"

	"github.com/mindchain/mind-chain/mc-service/log"
)

func SetupDefaults() {
	SetGlobalLogHandler(log.LogfmtHandlerWithLevel(os.Stdout, log.LevelInfo))
}
