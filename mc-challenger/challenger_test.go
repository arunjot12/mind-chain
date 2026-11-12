package op_challenger

import (
	"context"
	"testing"

	"github.com/mindchain/mind-chain/mc-challenger/config"
	"github.com/mindchain/mind-chain/mc-challenger/metrics"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/testlog"
	"github.com/stretchr/testify/require"
)

func TestMainShouldReturnErrorWhenConfigInvalid(t *testing.T) {
	cfg := &config.Config{}
	app, err := Main(context.Background(), testlog.Logger(t, log.LevelInfo), cfg, metrics.NoopMetrics)
	require.ErrorIs(t, err, cfg.Check())
	require.Nil(t, app)
}
