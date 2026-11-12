package stack

import (
	"github.com/mindchain/mind-chain/mc-service/apis"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

// L1CLNode is a L1 ethereum consensus-layer node, aka Beacon node.
// This node may not be a full beacon node, and instead run a mock L1 consensus node.
type L1CLNode interface {
	Common
	ChainID() eth.ChainID

	BeaconHTTPAddr() string
	BeaconClient() apis.BeaconClient
}
