package stack

import (
	"github.com/mindchain/mind-chain/mc-service/apis"
	"github.com/mindchain/mind-chain/mc-service/client"
)

// RollupBoostNode is a shim service between an L2 consensus-layer node and an L2 ethereum execution-layer node
type RollupBoostNode interface {
	L2EthClient() apis.L2EthClient
	L2EngineClient() apis.EngineClient
	FlashblocksClient() *client.WSClient

	ELNode
}
