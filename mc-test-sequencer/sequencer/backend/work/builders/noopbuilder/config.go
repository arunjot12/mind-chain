package noopbuilder

import (
	"context"

	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/backend/work"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/seqtypes"
)

type Config struct {
}

func (c *Config) Start(ctx context.Context, id seqtypes.BuilderID, opts *work.ServiceOpts) (work.Builder, error) {
	return &Builder{
		id:       id,
		registry: opts.Jobs,
	}, nil
}
