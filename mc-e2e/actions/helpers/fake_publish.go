package helpers

import (
	"context"

	"github.com/mindchain/mind-chain/mc-service/apis"
	"github.com/mindchain/mind-chain/mc-service/log"
	opsigner "github.com/mindchain/mind-chain/mc-service/signer"
)

// FakePublishAPI is used to instantiate an opstack API backend without full block publishing
type FakePublishAPI struct {
	Log log.Logger
}

var _ apis.PublishAPI = (*FakePublishAPI)(nil)

func (f *FakePublishAPI) PublishBlock(ctx context.Context, signed *opsigner.SignedExecutionPayloadEnvelope) error {
	f.Log.Info("Publish block", "signed", signed)
	return nil
}
