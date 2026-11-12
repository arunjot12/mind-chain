package sources

import (
	"context"
	"fmt"

	"github.com/mindchain/mind-chain/mc-service/client"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

type FollowClient struct {
	rollupClient *RollupClient
}

type FollowStatus struct {
	SafeL2      eth.L2BlockRef
	LocalSafeL2 eth.L2BlockRef
	FinalizedL2 eth.L2BlockRef
	CurrentL1   eth.L1BlockRef
}

func NewFollowClient(client client.RPC) (*FollowClient, error) {
	rollupClient := NewRollupClient(client)
	return &FollowClient{rollupClient: rollupClient}, nil
}

func (s *FollowClient) GetFollowStatus(ctx context.Context) (*FollowStatus, error) {
	status, err := s.rollupClient.SyncStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch external syncStatus: %w", err)
	}
	return &FollowStatus{
		FinalizedL2: status.FinalizedL2,
		SafeL2:      status.SafeL2,
		LocalSafeL2: status.LocalSafeL2,
		CurrentL1:   status.CurrentL1,
	}, nil
}
