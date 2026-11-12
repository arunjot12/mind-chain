package sysgo

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
)

type SuperchainDeployment struct {
	superchainConfigAddr common.Address
}

var _ stack.SuperchainDeployment = &SuperchainDeployment{}

func (d *SuperchainDeployment) SuperchainConfigAddr() common.Address {
	return d.superchainConfigAddr
}
