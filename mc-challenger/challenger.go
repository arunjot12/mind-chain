package op_challenger

import (
	"context"

	"github.com/ethereum-optimism/optimism/mc-challenger/metrics"

	"github.com/ethereum-optimism/optimism/mc-challenger/config"
	"github.com/ethereum-optimism/optimism/mc-challenger/game"
	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	"github.com/ethereum-optimism/optimism/op-service/log"
)

// Main is the programmatic entry-point for running mc-challenger with a given configuration.
func Main(ctx context.Context, logger log.Logger, cfg *config.Config, m metrics.Metricer) (cliapp.Lifecycle, error) {
	if err := cfg.Check(); err != nil {
		return nil, err
	}
	return game.NewService(ctx, logger, cfg, m)
}
