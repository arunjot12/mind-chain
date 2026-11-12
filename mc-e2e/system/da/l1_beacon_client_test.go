package da

import (
	"context"
	"testing"

	op_e2e "github.com/mindchain/mind-chain/mc-e2e"
	"github.com/mindchain/mind-chain/mc-service/log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mindchain/mind-chain/mc-e2e/e2eutils/blobstore"
	"github.com/mindchain/mind-chain/mc-e2e/e2eutils/fakebeacon"
	"github.com/mindchain/mind-chain/mc-service/client"
	"github.com/mindchain/mind-chain/mc-service/sources"
	"github.com/mindchain/mind-chain/mc-service/testlog"
	"github.com/stretchr/testify/require"
)

func Test404NotFound(t *testing.T) {
	op_e2e.InitParallel(t)

	l := testlog.Logger(t, log.LevelInfo)

	blobStore := blobstore.New()
	beaconApi := fakebeacon.NewBeacon(l, blobStore, uint64(0), uint64(12))
	t.Cleanup(func() {
		_ = beaconApi.Close()
	})
	require.NoError(t, beaconApi.Start("127.0.0.1:0"))

	beaconCfg := sources.L1BeaconClientConfig{FetchAllSidecars: false}
	cl := sources.NewL1BeaconClient(sources.NewBeaconHTTPClient(client.NewBasicHTTPClient(beaconApi.BeaconAddr(), l)), beaconCfg)

	hashes := []common.Hash{{}}
	_, err := cl.GetBlobsByHash(context.Background(), uint64(0), hashes)
	require.ErrorIs(t, err, ethereum.NotFound)
}
