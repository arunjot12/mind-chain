package operatorfee

// NOTE: These utility functions have been converted from devnet-sdk to mc-devstack types
// but are currently unused by tests. They would need implementation updates if used.

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

func EnsureSufficientBalance(wallet *dsl.EOA, to common.Address, value eth.ETH) (err error) {
	return fmt.Errorf("not implemented for mc-devstack - utility function not used by tests")
}

func SendValueTx(wallet *dsl.EOA, to common.Address, value eth.ETH) (tx *gethTypes.Transaction, receipt *gethTypes.Receipt, err error) {
	return nil, nil, fmt.Errorf("not implemented for mc-devstack - utility function not used by tests")
}

func ReturnRemainingFunds(wallet *dsl.EOA, to common.Address) (receipt *gethTypes.Receipt, err error) {
	return nil, fmt.Errorf("not implemented for mc-devstack - utility function not used by tests")
}

func NewTestWallet(ctx context.Context, el dsl.ELNode) (*dsl.EOA, error) {
	return nil, fmt.Errorf("not implemented for mc-devstack - utility function not used by tests")
}
