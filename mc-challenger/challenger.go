package op_challenger

import (
	"context"

	"github.com/mindchain/mind-chain/mc-challenger/metrics"

	"github.com/mindchain/mind-chain/mc-challenger/config"
	"github.com/mindchain/mind-chain/mc-challenger/game"
	"github.com/mindchain/mind-chain/mc-service/cliapp"
	"github.com/mindchain/mind-chain/mc-service/log"
)

// Main is the programmatic entry-point for running mc-challenger with a given configuration.
func Main(ctx context.Context, logger log.Logger, cfg *config.Config, m metrics.Metricer) (cliapp.Lifecycle, error) {
	if err := cfg.Check(); err != nil {
		return nil, err
	}
	return game.NewService(ctx, logger, cfg, m)
}
