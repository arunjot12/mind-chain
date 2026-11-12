package apis

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type Faucet interface {
	ChainID(ctx context.Context) (eth.ChainID, error)
	RequestETH(ctx context.Context, addr common.Address, amount eth.ETH) error
	Balance(ctx context.Context) (eth.ETH, error)
}
