package p2p

import (
	"github.com/mindchain/mind-chain/mc-service/apis"
	"github.com/mindchain/mind-chain/mc-service/sources"
)

// backwards compatibility
var NamespaceRPC = sources.P2PNamespaceRPC

type API = apis.P2PClient
type PeerInfo = apis.PeerInfo
type PeerDump = apis.PeerDump
type PeerStats = apis.PeerStats
