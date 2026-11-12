package bindings

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

// SuperchainETHBridge binds the SuperchainETHBridge predeploy.
type SuperchainETHBridge struct {
	// Write functions. SendETH is payable: pass the amount with txplan.WithValue.
	SendETH   func(to common.Address, chainID eth.ChainID) TypedCall[eth.Bytes32]                             `sol:"sendETH"`
	RefundETH func(destination eth.ChainID, nonce *big.Int, to common.Address, amount eth.ETH) TypedCall[any] `sol:"refundETH"`

	// Read-only functions
	Refunded func(messageHash [32]byte) TypedCall[bool] `sol:"refunded"`
}
