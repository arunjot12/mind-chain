package backend

import (
	"github.com/mindchain/mind-chain/mc-service/eth"
	"github.com/mindchain/mind-chain/mc-service/event"
	opmetrics "github.com/mindchain/mind-chain/mc-service/metrics"
	"github.com/mindchain/mind-chain/mc-service/sources/caching"
	"github.com/mindchain/mind-chain/mc-supervisor/supervisor/backend/db/logs"
	"github.com/mindchain/mind-chain/mc-supervisor/supervisor/types"
)

type Metrics interface {
	CacheAdd(chainID eth.ChainID, label string, cacheSize int, evicted bool)
	CacheGet(chainID eth.ChainID, label string, hit bool)

	RecordCrossUnsafe(chainID eth.ChainID, ref types.BlockSeal)
	RecordCrossSafe(chainID eth.ChainID, ref types.BlockSeal)
	RecordLocalSafe(chainID eth.ChainID, ref types.BlockSeal)
	RecordLocalUnsafe(chainID eth.ChainID, ref types.BlockSeal)

	RecordDBEntryCount(chainID eth.ChainID, kind string, count int64)
	RecordDBSearchEntriesRead(chainID eth.ChainID, count int64)

	RecordAccessListVerifyFailure(chainID eth.ChainID)

	opmetrics.RPCMetricer
	event.Metrics
}

// chainMetrics is an adapter between the metrics API expected by clients that assume there's only a single chain
// and the actual metrics implementation which requires a chain ID to identify the source chain.
type chainMetrics struct {
	chainID  eth.ChainID
	delegate Metrics
}

func newChainMetrics(chainID eth.ChainID, delegate Metrics) *chainMetrics {
	return &chainMetrics{
		chainID:  chainID,
		delegate: delegate,
	}
}

func (c *chainMetrics) RecordCrossUnsafe(seal types.BlockSeal) {
	c.delegate.RecordCrossUnsafe(c.chainID, seal)
}

func (c *chainMetrics) RecordCrossSafe(seal types.BlockSeal) {
	c.delegate.RecordCrossSafe(c.chainID, seal)
}

func (c *chainMetrics) RecordLocalSafe(seal types.BlockSeal) {
	c.delegate.RecordLocalSafe(c.chainID, seal)
}

func (c *chainMetrics) RecordLocalUnsafe(seal types.BlockSeal) {
	c.delegate.RecordLocalUnsafe(c.chainID, seal)
}

func (c *chainMetrics) CacheAdd(label string, cacheSize int, evicted bool) {
	c.delegate.CacheAdd(c.chainID, label, cacheSize, evicted)
}

func (c *chainMetrics) CacheGet(label string, hit bool) {
	c.delegate.CacheGet(c.chainID, label, hit)
}

func (c *chainMetrics) RecordDBEntryCount(kind string, count int64) {
	c.delegate.RecordDBEntryCount(c.chainID, kind, count)
}

func (c *chainMetrics) RecordDBSearchEntriesRead(count int64) {
	c.delegate.RecordDBSearchEntriesRead(c.chainID, count)
}

func (c *chainMetrics) RecordAccessListVerifyFailure() {
	c.delegate.RecordAccessListVerifyFailure(c.chainID)
}

var _ caching.Metrics = (*chainMetrics)(nil)
var _ logs.Metrics = (*chainMetrics)(nil)
