package standardcommitter

import (
	"context"

	"github.com/mindchain/mind-chain/mc-service/client"
	"github.com/mindchain/mind-chain/mc-service/endpoint"
	"github.com/mindchain/mind-chain/mc-service/sources"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/backend/work"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/seqtypes"
)

type Config struct {
	// RPC to commit block to using op-stack RPC
	RPC endpoint.MustRPC `yaml:"rpc"`
}

func (c *Config) Start(ctx context.Context, id seqtypes.CommitterID, opts *work.ServiceOpts) (work.Committer, error) {
	rpcCl, err := client.NewRPC(ctx, opts.Log, c.RPC.Value.RPC(), client.WithLazyDial())
	if err != nil {
		return nil, err
	}
	cl := sources.NewOPStackClient(rpcCl)
	return &Committer{
		id:      id,
		log:     opts.Log,
		api:     cl,
		onClose: rpcCl.Close,
	}, nil
}
