package noopcommitter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/testlog"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/seqtypes"
)

func TestNoopCommitter(t *testing.T) {
	logger := testlog.Logger(t, log.LevelDebug)
	id := seqtypes.CommitterID("foo")
	x := NewCommitter(id, logger)

	err := x.Commit(context.Background(), nil)
	require.NoError(t, err)

	require.NoError(t, x.Close())
	require.Equal(t, "noop-committer-foo", x.String())
	require.Equal(t, id, x.ID())
}
