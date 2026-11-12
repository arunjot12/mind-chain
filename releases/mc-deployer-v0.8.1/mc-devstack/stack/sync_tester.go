package stack

import (
	"github.com/mindchain/mind-chain/mc-service/apis"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type SyncTester interface {
	Common
	ChainID() eth.ChainID
	API() apis.SyncTester

	APIWithSession(sessionID string) apis.SyncTester
}
