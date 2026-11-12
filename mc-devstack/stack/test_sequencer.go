package stack

import (
	"github.com/mindchain/mind-chain/mc-service/apis"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

// TestSequencer
type TestSequencer interface {
	Common

	AdminAPI() apis.TestSequencerAdminAPI
	BuildAPI() apis.TestSequencerBuildAPI
	ControlAPI(chainID eth.ChainID) apis.TestSequencerControlAPI
}
