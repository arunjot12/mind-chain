package interop

import (
	"os"

	"github.com/ethereum/go-ethereum/log"
	oplog "github.com/mindchain/mind-chain/mc-service/log"
)

func init() {
	oplog.SetGlobalLogHandler(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelDebug, true))
}
