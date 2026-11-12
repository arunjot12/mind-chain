package bindings

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mindchain/mind-chain/mc-service/eth"
	supTypes "github.com/mindchain/mind-chain/mc-supervisor/supervisor/types"
)

type ABIIdentifier struct {
	Origin      common.Address
	BlockNumber *big.Int
	LogIndex    *big.Int
	Timestamp   *big.Int
	ChainId     *big.Int
}

type CrossL2Inbox struct {
	ValidateMessage func(identifier supTypes.Identifier, msgHash eth.Bytes32) TypedCall[any] `sol:"validateMessage"`
}
