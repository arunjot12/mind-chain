package synctester

import (
	"context"
	"fmt"

	opservice "github.com/mindchain/mind-chain/mc-service"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/log/logcli"

	"github.com/mindchain/mind-chain/mc-service/cliapp"
	"github.com/mindchain/mind-chain/mc-sync-tester/config"
	"github.com/mindchain/mind-chain/mc-sync-tester/flags"
	"github.com/urfave/cli/v2"
)

type MainFn func(ctx context.Context, cfg *config.Config, logger log.Logger) (cliapp.Lifecycle, error)

// Main is the entrypoint into the service.
// This method returns a cliapp.LifecycleAction, to create an mc-service CLI-lifecycle-managed service with.
func Main(version string, fn MainFn) cliapp.LifecycleAction {
	return func(cliCtx *cli.Context, closeApp context.CancelCauseFunc) (cliapp.Lifecycle, error) {
		if err := flags.CheckRequired(cliCtx); err != nil {
			return nil, err
		}
		cfg := flags.ConfigFromCLI(cliCtx, version)
		if err := cfg.Check(); err != nil {
			return nil, fmt.Errorf("invalid CLI flags: %w", err)
		}

		l := logcli.NewLogger(logcli.AppOut(cliCtx), cfg.LogConfig)
		logcli.SetGlobalLogHandler(l.Handler())
		opservice.ValidateEnvVars(flags.EnvVarPrefix, flags.Flags, l)

		l.Info("Initializing mc-sync-tester")
		return fn(cliCtx.Context, cfg, l)
	}
}
