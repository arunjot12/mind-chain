package p2p

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/mindchain/mind-chain/mc-service/client"
	"github.com/mindchain/mind-chain/mc-service/sources"
)

type Client struct {
	c sources.P2PClient
}

// Legacy for supporting backwards compatibility
func NewClient(c *rpc.Client) *Client {
	return &Client{c: *sources.NewP2PClient(client.NewBaseRPCClient(c))}
}
