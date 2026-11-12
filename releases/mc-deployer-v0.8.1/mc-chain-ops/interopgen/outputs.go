package interopgen

import (
	"github.com/ethereum/go-ethereum/core"

	"github.com/mindchain/mind-chain/mc-chain-ops/interopgen/config"
	"github.com/mindchain/mind-chain/mc-node/rollup"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type L1Output struct {
	Genesis *core.Genesis
}

type L2Output struct {
	Genesis   *core.Genesis
	RollupCfg *rollup.Config
}

type WorldOutput struct {
	L1  *L1Output
	L2s map[string]*L2Output
}

func (wo *WorldOutput) RollupConfigSet() config.StaticRollupConfigSet {
	rcfgs := make(map[eth.ChainID]*rollup.Config)
	for _, rcfg := range wo.L2s {
		rcfgs[eth.ChainIDFromBig(rcfg.RollupCfg.L2ChainID)] = rcfg.RollupCfg
	}
	return config.StaticRollupConfigSetFromRollupConfigMap(rcfgs, config.StaticTimestamp(wo.L1.Genesis.Timestamp))
}
