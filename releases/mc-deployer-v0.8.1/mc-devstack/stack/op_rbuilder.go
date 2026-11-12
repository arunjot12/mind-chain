package stack

import (
	"github.com/mindchain/mind-chain/mc-service/apis"
	"github.com/mindchain/mind-chain/mc-service/client"
)

// OPRBuilderNode is a L2 ethereum execution-layer node
type OPRBuilderNode interface {
	L2EthClient() apis.L2EthClient
	L2EngineClient() apis.EngineClient
	FlashblocksClient() *client.WSClient
	UpdateRuleSet(rulesYaml string) error

	ELNode
}
