package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/log"

	monitor "github.com/mindchain/mind-chain/mc-dispute-mon"
	"github.com/mindchain/mind-chain/mc-dispute-mon/config"
	"github.com/mindchain/mind-chain/mc-dispute-mon/flags"
	"github.com/mindchain/mind-chain/mc-dispute-mon/version"
	opservice "github.com/mindchain/mind-chain/mc-service"
	"github.com/mindchain/mind-chain/mc-service/cliapp"
	"github.com/mindchain/mind-chain/mc-service/ctxinterrupt"
	oplog "github.com/mindchain/mind-chain/mc-service/log"
)

var (
	GitCommit = ""
	GitDate   = ""
)

// VersionWithMeta holds the textual version string including the metadata.
var VersionWithMeta = opservice.FormatVersion(version.Version, GitCommit, GitDate, version.Meta)

func main() {
	args := os.Args
	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	if err := run(ctx, args, monitor.Main); err != nil {
		log.Crit("Application failed", "err", err)
	}
}

type ConfiguredLifecycle func(ctx context.Context, log log.Logger, config *config.Config) (cliapp.Lifecycle, error)

func run(ctx context.Context, args []string, action ConfiguredLifecycle) error {
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Version = VersionWithMeta
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Name = "mc-dispute-mon"
	app.Usage = "Monitor dispute games"
	app.Description = "Monitors output proposals and dispute games."
	app.Action = cliapp.LifecycleCmd(func(ctx *cli.Context, close context.CancelCauseFunc) (cliapp.Lifecycle, error) {
		logger, err := setupLogging(ctx)
		if err != nil {
			return nil, err
		}
		logger.Info("Starting mc-dispute-mon", "version", VersionWithMeta)

		cfg, err := flags.NewConfigFromCLI(ctx)
		if err != nil {
			return nil, err
		}
		return action(ctx.Context, logger, cfg)
	})
	return app.RunContext(ctx, args)
}

func setupLogging(ctx *cli.Context) (log.Logger, error) {
	logCfg := oplog.ReadCLIConfig(ctx)
	logger := oplog.NewLogger(oplog.AppOut(ctx), logCfg)
	oplog.SetGlobalLogHandler(logger.Handler())
	return logger, nil
}
