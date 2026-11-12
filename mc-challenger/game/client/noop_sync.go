package client

import (
	"context"

	"github.com/mindchain/mind-chain/mc-challenger/game/types"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type NoopSyncStatusValidator struct{}

func (n *NoopSyncStatusValidator) ValidateNodeSynced(_ context.Context, _ eth.BlockID) error {
	return nil
}

var _ types.SyncValidator = (*NoopSyncStatusValidator)(nil)
