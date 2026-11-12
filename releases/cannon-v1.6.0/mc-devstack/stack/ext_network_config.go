package stack

import "github.com/mindchain/mind-chain/mc-service/eth"

type ExtNetworkConfig struct {
	L2NetworkName      string
	L1ChainID          eth.ChainID
	L2ELEndpoint       string
	L1CLBeaconEndpoint string
	L1ELEndpoint       string
}
