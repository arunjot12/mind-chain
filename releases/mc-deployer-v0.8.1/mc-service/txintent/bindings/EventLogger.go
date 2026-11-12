package bindings

import (
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type EventLogger struct {
	EmitLog func(topics []eth.Bytes32, data []byte) TypedCall[any] `sol:"emitLog"`
}
