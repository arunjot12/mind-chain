package stack

import "github.com/mindchain/mind-chain/mc-service/eth"

// L2Proposer is a L2 output proposer, posting claims of L2 state to L1.
type L2Proposer interface {
	Common
	ChainID() eth.ChainID
}
