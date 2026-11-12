package noopcommitter

import (
	"context"

	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/backend/work"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/seqtypes"
)

type Config struct {
}

func (c *Config) Start(ctx context.Context, id seqtypes.CommitterID, opts *work.ServiceOpts) (work.Committer, error) {
	return &Committer{
		id:  id,
		log: opts.Log,
	}, nil
}
