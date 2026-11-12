package engine

import (
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type PayloadInvalidEvent struct {
	Envelope *eth.ExecutionPayloadEnvelope
	Err      error
}

func (ev PayloadInvalidEvent) String() string {
	return "payload-invalid"
}
