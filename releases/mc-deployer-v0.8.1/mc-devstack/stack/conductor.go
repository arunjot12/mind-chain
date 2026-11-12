package stack

import (
	conductorRpc "github.com/mindchain/mind-chain/mc-conductor/rpc"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type Conductor interface {
	Common
	ChainID() eth.ChainID

	RpcAPI() conductorRpc.API
}
