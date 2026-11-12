package mon

import (
	"github.com/mindchain/mind-chain/mc-dispute-mon/mon/types"
	"github.com/mindchain/mind-chain/mc-service/log"
)

type NodeEndpointOutOfSyncMetrics interface {
	RecordNodeEndpointOutOfSyncCount(count int)
}

type NodeEndpointOutOfSyncMonitor struct {
	logger  log.Logger
	metrics NodeEndpointOutOfSyncMetrics
}

func NewNodeEndpointOutOfSyncMonitor(logger log.Logger, metrics NodeEndpointOutOfSyncMetrics) *NodeEndpointOutOfSyncMonitor {
	return &NodeEndpointOutOfSyncMonitor{
		logger:  logger,
		metrics: metrics,
	}
}

func (m *NodeEndpointOutOfSyncMonitor) CheckNodeEndpointOutOfSync(games []*types.CommonGameData) {
	totalOutOfSync := 0

	for _, game := range games {
		totalOutOfSync += game.NodeEndpointOutOfSyncCount
	}

	m.metrics.RecordNodeEndpointOutOfSyncCount(totalOutOfSync)
}
