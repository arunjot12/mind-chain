package client

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mindchain/mind-chain/devnet-sdk/contracts/bindings"
	"github.com/mindchain/mind-chain/devnet-sdk/interfaces"
	"github.com/mindchain/mind-chain/devnet-sdk/types"
)

type L2ToL2CrossDomainMessengerBinding struct {
	contractAddress types.Address
	client          *ethclient.Client
	binding         *bindings.L2ToL2CrossDomainMessenger
	abi             *abi.ABI
}

var _ interfaces.L2ToL2CrossDomainMessenger = (*L2ToL2CrossDomainMessengerBinding)(nil)

func (b *L2ToL2CrossDomainMessengerBinding) ABI() *abi.ABI {
	return b.abi
}
