package monitor

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/mindchain/mind-chain/mc-dispute-mon/config"
	"github.com/mindchain/mind-chain/mc-service/testlog"
	"github.com/stretchr/testify/require"
)

func TestMainShouldReturnErrorWhenConfigInvalid(t *testing.T) {
	cfg := &config.Config{}
	app, err := Main(context.Background(), testlog.Logger(t, log.LvlInfo), cfg)
	require.ErrorIs(t, err, cfg.Check())
	require.Nil(t, app)
}
