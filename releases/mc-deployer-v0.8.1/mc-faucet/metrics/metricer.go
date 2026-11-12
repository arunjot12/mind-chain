package metrics

import (
	ftypes "github.com/mindchain/mind-chain/mc-faucet/faucet/backend/types"
	"github.com/mindchain/mind-chain/mc-service/eth"
	"github.com/mindchain/mind-chain/mc-service/txmgr/metrics"
)

type Metricer interface {
	RecordInfo(version string)
	RecordUp()

	RecordFundAction(faucet ftypes.FaucetID, chainID eth.ChainID, amount eth.ETH) (onDone func(err error))

	metrics.TxMetricer
}
