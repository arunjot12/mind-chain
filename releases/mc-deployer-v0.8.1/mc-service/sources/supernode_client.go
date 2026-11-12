package sources

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mindchain/mind-chain/mc-service/client"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type SuperNodeClient struct {
	rpc client.RPC
}

func NewSuperNodeClient(rpc client.RPC) *SuperNodeClient {
	return &SuperNodeClient{
		rpc: rpc,
	}
}

func (c *SuperNodeClient) SuperRootAtTimestamp(ctx context.Context, timestamp uint64) (result eth.SuperRootAtTimestampResponse, err error) {
	err = c.rpc.CallContext(ctx, &result, "superroot_atTimestamp", hexutil.Uint64(timestamp))
	return
}

func (c *SuperNodeClient) SyncStatus(ctx context.Context) (result eth.SuperNodeSyncStatusResponse, err error) {
	err = c.rpc.CallContext(ctx, &result, "supernode_syncStatus")
	return
}

func (cl *SuperNodeClient) Close() {
	cl.rpc.Close()
}
