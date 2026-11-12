package apis

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mindchain/mind-chain/mc-service/eth"
)

// BeaconClient is a thin wrapper over the Beacon APIs.
type BeaconClient interface {
	ConfigSpec(ctx context.Context) (eth.APIConfigResponse, error)
	BeaconGenesis(ctx context.Context) (eth.APIGenesisResponse, error)
	BeaconBlobs(ctx context.Context, slot uint64, hashes []common.Hash) (eth.APIBeaconBlobsResponse, error)
}
