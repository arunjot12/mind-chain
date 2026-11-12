package apis

import (
	"context"

	"github.com/mindchain/mind-chain/mc-service/eth"
)

// SupernodeQueryAPI is the minimal API surface the devstack DSL needs from mc-supernode.
// It is intentionally small and can be expanded as needed.
type SupernodeQueryAPI interface {
	SuperRootAtTimestamp(ctx context.Context, timestamp uint64) (eth.SuperRootAtTimestampResponse, error)
	SyncStatus(ctx context.Context) (eth.SuperNodeSyncStatusResponse, error)
}
