package metrics

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/mindchain/mind-chain/mc-service/log"
	opmetrics "github.com/mindchain/mind-chain/mc-service/metrics"
	txmetrics "github.com/mindchain/mind-chain/mc-service/txmgr/metrics"
)

type noopMetrics struct {
	opmetrics.NoopRefMetrics
	txmetrics.NoopTxMetrics
}

var NoopMetrics Metricer = new(noopMetrics)

func (*noopMetrics) RecordInfo(version string) {}
func (*noopMetrics) RecordUp()                 {}

func (*noopMetrics) RecordDripExecuted(name string) {}

func (*noopMetrics) StartBalanceMetrics(log.Logger, *ethclient.Client, common.Address) io.Closer {
	return nil
}
