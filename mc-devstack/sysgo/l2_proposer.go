package sysgo

import (
	ps "github.com/mindchain/mind-chain/mc-proposer/proposer"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type L2Proposer struct {
	name    string
	chainID eth.ChainID
	service *ps.ProposerService
	userRPC string
}

type ProposerOption func(target ComponentTarget, cfg *ps.CLIConfig)
