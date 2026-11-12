package fakepos

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/mindchain/mind-chain/mc-e2e/e2eutils/geth"
)

type Config struct {
	ChainConfig       *params.ChainConfig
	Backend           Blockchain
	EngineAPI         geth.EngineAPI
	Beacon            Beacon
	FinalizedDistance uint64
	SafeDistance      uint64
	BlockTime         uint64
}
