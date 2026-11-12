package stack

import (
	conductorRpc "github.com/mindchain/mind-chain/mc-conductor/rpc"
	"github.com/mindchain/mind-chain/mc-service/client"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type Conductor interface {
	Common
	ChainID() eth.ChainID

	RpcAPI() conductorRpc.API

	// ProxyRPC exposes the conductor's RPC endpoint for raw calls. With RPC
	// proxying enabled, mc-conductor forwards execution (eth_*), rollup
	// (optimism_*), and admin (admin_*) requests to its sequencer while it is
	// the leader, so batchers and proposers can follow the active sequencer.
	ProxyRPC() client.RPC

	// ConsensusEndpoint is the Raft consensus address of this conductor, used
	// when (re-)adding it to a cluster.
	ConsensusEndpoint() string
}
