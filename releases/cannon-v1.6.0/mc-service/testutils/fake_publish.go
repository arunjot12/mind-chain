package testutils

import (
	"context"

	"github.com/ethereum/go-ethereum/log"

	"github.com/mindchain/mind-chain/mc-service/apis"
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
