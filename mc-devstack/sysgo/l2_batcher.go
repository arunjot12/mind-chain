package sysgo

import (
	bss "github.com/mindchain/mind-chain/mc-batcher/batcher"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type L2Batcher struct {
	name    string
	chainID eth.ChainID
	service *bss.BatcherService
	rpc     string
	l1RPC   string
	l2CLRPC string
	l2ELRPC string
}

func (b *L2Batcher) UserRPC() string {
	return b.rpc
}

type BatcherOption func(target ComponentTarget, cfg *bss.CLIConfig)
