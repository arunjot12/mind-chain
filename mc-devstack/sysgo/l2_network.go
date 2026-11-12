package sysgo

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"

	"github.com/mindchain/mind-chain/mc-chain-ops/devkeys"
	opgenesis "github.com/mindchain/mind-chain/mc-core/genesis"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
	"github.com/mindchain/mind-chain/mc-node/rollup"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type L2Network struct {
	name       string
	chainID    eth.ChainID
	l1ChainID  eth.ChainID
	genesis    *opgenesis.Genesis
	rollupCfg  *rollup.Config
	deployment *L2Deployment
	opcmImpl   common.Address
	mipsImpl   common.Address
	keys       devkeys.Keys
}

func (c *L2Network) Name() string {
	return c.name
}

func (c *L2Network) ChainID() eth.ChainID {
	return c.chainID
}

func (c *L2Network) L1ChainID() eth.ChainID {
	return c.l1ChainID
}

func (c *L2Network) ChainConfig() *params.ChainConfig {
	return c.genesis.Config.GethChainConfig()
}

func (c *L2Network) RollupConfig() *rollup.Config {
	return c.rollupCfg
}

func (c *L2Network) Deployment() stack.L2Deployment {
	return c.deployment
}

func (c *L2Network) Keys() devkeys.Keys {
	return c.keys
}
