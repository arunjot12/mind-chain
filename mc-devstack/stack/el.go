package stack

import (
	"time"

	"github.com/mindchain/mind-chain/mc-service/apis"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type ELNode interface {
	Common
	ChainID() eth.ChainID
	EthClient() apis.EthClient
	UserRPC() string
	TransactionTimeout() time.Duration
}
