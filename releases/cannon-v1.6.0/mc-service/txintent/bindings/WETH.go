package bindings

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type WETH struct {
	BalanceOf func(addr common.Address) TypedCall[eth.ETH]                 `sol:"balanceOf"`
	Transfer  func(dest common.Address, amount eth.ETH) TypedCall[bool]    `sol:"transfer"`
	Approve   func(spender common.Address, amount eth.ETH) TypedCall[bool] `sol:"approve"`
	Deposit   func() TypedCall[any]                                        `sol:"deposit"`
}
