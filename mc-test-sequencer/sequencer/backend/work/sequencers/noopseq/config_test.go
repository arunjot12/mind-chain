package noopseq

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/testlog"
	"github.com/mindchain/mind-chain/mc-test-sequencer/metrics"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/backend/work"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/seqtypes"
)

func TestConfig(t *testing.T) {
	logger := testlog.Logger(t, log.LevelInfo)
	cfg := &Config{}
	id := seqtypes.SequencerID("test")
	ensemble := &work.Ensemble{}
	opts := &work.ServiceOpts{
		StartOpts: &work.StartOpts{
			Log:     logger,
			Metrics: &metrics.NoopMetrics{},
		},
		Services: ensemble,
	}
	sequencer, err := cfg.Start(context.Background(), id, opts)
	require.NoError(t, err)
	require.Equal(t, id, sequencer.ID())
}
