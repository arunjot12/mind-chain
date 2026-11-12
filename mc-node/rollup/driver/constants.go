package driver

import "github.com/mindchain/mind-chain/mc-node/rollup/sequencing"

// aliases to not disrupt mc-conductor code
var (
	ErrSequencerAlreadyStarted = sequencing.ErrSequencerAlreadyStarted
	ErrSequencerAlreadyStopped = sequencing.ErrSequencerAlreadyStopped
)
