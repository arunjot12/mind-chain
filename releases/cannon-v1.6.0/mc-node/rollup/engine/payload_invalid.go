package engine

import (
	"context"

	"github.com/mindchain/mind-chain/mc-service/eth"
)

type PayloadInvalidEvent struct {
	Envelope *eth.ExecutionPayloadEnvelope
	Err      error
}

func (ev PayloadInvalidEvent) String() string {
	return "payload-invalid"
}

func (eq *EngineController) onPayloadInvalid(ctx context.Context, ev PayloadInvalidEvent) {
	eq.log.Warn("Payload was invalid", "block", ev.Envelope.ExecutionPayload.ID(),
		"err", ev.Err, "timestamp", uint64(ev.Envelope.ExecutionPayload.Timestamp))
}
