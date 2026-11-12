package contracts

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mindchain/mind-chain/devnet-sdk/contracts/registry/client"
	"github.com/mindchain/mind-chain/devnet-sdk/contracts/registry/empty"
	"github.com/mindchain/mind-chain/devnet-sdk/interfaces"
)

// NewClientRegistry creates a new Registry that uses the provided client
func NewClientRegistry(c *ethclient.Client) interfaces.ContractsRegistry {
	return &client.ClientRegistry{Client: c}
}

func NewEmptyRegistry() interfaces.ContractsRegistry {
	return &empty.EmptyRegistry{}
}
