package health

import (
	"context"

	"github.com/mindchain/mind-chain/mc-service/eth"
)

// SupervisorHealthAPI defines the interface for the supervisor's health check.
type SupervisorHealthAPI interface {
	SyncStatus(ctx context.Context) (eth.SupervisorSyncStatus, error)
}
