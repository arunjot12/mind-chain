package noopsigner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mindchain/mind-chain/mc-service/eth"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/testlog"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/seqtypes"
)

func TestNoopSigner(t *testing.T) {
	logger := testlog.Logger(t, log.LevelDebug)
	id := seqtypes.SignerID("foo")
	x := NewSigner(id, logger)

	_, err := x.Sign(context.Background(), nil)
	require.NoError(t, err)

	require.NoError(t, x.Close())
	require.Equal(t, "noop-signer-foo", x.String())
	require.Equal(t, id, x.ID())

	require.Equal(t, eth.ChainID{}, x.ChainID())
}
