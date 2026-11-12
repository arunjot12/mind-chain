package sysgo

import (
	"github.com/mindchain/mind-chain/mc-challenger/config"
	"github.com/mindchain/mind-chain/mc-service/cliapp"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type L2Challenger struct {
	name     string
	chainIDs []eth.ChainID
	service  cliapp.Lifecycle
	config   *config.Config
}

func (p *L2Challenger) Config() *config.Config {
	return p.config
}
