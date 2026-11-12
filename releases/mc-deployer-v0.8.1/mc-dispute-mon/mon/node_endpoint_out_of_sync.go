package mon

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/mindchain/mind-chain/mc-dispute-mon/mon/types"
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
