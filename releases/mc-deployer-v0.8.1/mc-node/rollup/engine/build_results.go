package engine

import (
	"time"

	"github.com/mindchain/mind-chain/mc-service/eth"
)

// BuildStartResult contains the result of starting a block-building job.
type BuildStartResult struct {
	Info         eth.PayloadInfo
	BuildStarted time.Time
	Parent       eth.L2BlockRef
}

// SealResult contains the result of sealing a block-building job.
type SealResult struct {
	Envelope *eth.ExecutionPayloadEnvelope
	Ref      eth.L2BlockRef
}
