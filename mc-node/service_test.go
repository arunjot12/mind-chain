package opnode

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-node/flags"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func syncConfigCliApp() *cli.App {
	syncConfigFlags := []cli.Flag{
		flags.SequencerEnabledFlag,
		flags.L2EngineSyncEnabled,
		flags.SyncModeFlag,
		flags.SyncModeOffsetELSafeFlag,
		flags.L2FollowSource,
		flags.L2EngineKind,
		flags.SkipSyncStartCheck,
	}
	return &cli.App{
		Flags: syncConfigFlags,
		Action: func(c *cli.Context) error {
			_, err := NewSyncConfig(c, log.New())
			return err
		},
	}
}

func run(args []string) error {
	return syncConfigCliApp().Run(append([]string{"test"}, args...))
}

func TestNewSyncConfigDefault(t *testing.T) {
	require.NoError(t, run(nil))
}
