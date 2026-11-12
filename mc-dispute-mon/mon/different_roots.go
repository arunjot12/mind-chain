package mon

import (
	"github.com/mindchain/mind-chain/mc-dispute-mon/mon/types"
	"github.com/mindchain/mind-chain/mc-service/log"
)

type DifferentRootMetrics interface {
	RecordDifferentRootGames(count int)
}

type DifferentRootMonitor struct {
	logger  log.Logger
	metrics DifferentRootMetrics
}

func NewDifferentRootMonitor(logger log.Logger, metrics DifferentRootMetrics) *DifferentRootMonitor {
	return &DifferentRootMonitor{
		logger:  logger,
		metrics: metrics,
	}
}

func (m *DifferentRootMonitor) CheckDifferentRoots(games []*types.CommonGameData) {
	count := 0
	for _, game := range games {
		if game.NodeEndpointDifferentRoots {
			count++
			m.logger.Debug("Different roots detected",
				"game", game.Proxy,
				"l2SequenceNumber", game.L2SequenceNumber,
				"rootClaim", game.RootClaim)
		}
	}

	m.metrics.RecordDifferentRootGames(count)

	if count > 0 {
		m.logger.Info("Different roots summary", "gamesWithDifferentRoots", count, "totalGames", len(games))
	}
}
