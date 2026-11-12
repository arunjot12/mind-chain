package monitor

import (
	"context"

	"github.com/mindchain/mind-chain/mc-dispute-mon/config"
	"github.com/mindchain/mind-chain/mc-dispute-mon/mon"
	"github.com/mindchain/mind-chain/mc-service/log"
)

func Main(ctx context.Context, logger log.Logger, cfg *config.Config, options ...mon.ServiceOption) (*mon.Service, error) {
	if err := cfg.Check(); err != nil {
		return nil, err
	}
	return mon.NewService(ctx, logger, cfg, options...)
}
