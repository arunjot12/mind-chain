package stack

import (
	"github.com/mindchain/mind-chain/mc-challenger/config"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type L2Challenger interface {
	Common
	ChainID() eth.ChainID
	Config() *config.Config
}
