package frontend

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/backend/work"
)

type EthTxFrontend struct {
	Sequencer work.Sequencer
	Logger    log.Logger
}

func (etf *EthTxFrontend) SendRawTransaction(ctx context.Context, tx hexutil.Bytes) error {
	etf.Logger.Debug("EthTxFrontend SendRawTransaaction request", "tx", tx)

	return toJsonError(etf.Sequencer.BuildJob().IncludeTx(ctx, tx))
}
