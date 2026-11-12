package main

import (
	"context"
	"os"

	"github.com/mindchain/mind-chain/mc-challenger/metrics"
	"github.com/urfave/cli/v2"

	challenger "github.com/mindchain/mind-chain/mc-challenger"
	"github.com/mindchain/mind-chain/mc-challenger/config"
	"github.com/mindchain/mind-chain/mc-challenger/flags"
	"github.com/mindchain/mind-chain/mc-challenger/version"
	opservice "github.com/mindchain/mind-chain/mc-service"
	"github.com/mindchain/mind-chain/mc-service/cliapp"
	"github.com/mindchain/mind-chain/mc-service/ctxinterrupt"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/log/logcli"
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
	if err := run(ctx, args, func(ctx context.Context, l log.Logger, config *config.Config) (cliapp.Lifecycle, error) {
		return challenger.Main(ctx, l, config, metrics.NewMetrics())
	}); err != nil {
		log.Crit("Application failed", "err", err)
	}
}

type ConfiguredLifecycle func(ctx context.Context, log log.Logger, config *config.Config) (cliapp.Lifecycle, error)

func run(ctx context.Context, args []string, action ConfiguredLifecycle) error {
	logcli.SetupDefaults()

	app := cli.NewApp()
	app.Version = VersionWithMeta
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Name = "mc-challenger"
	app.Usage = "Challenge outputs"
	app.Description = "Ensures that on chain outputs are correct."
	app.Commands = []*cli.Command{
		ListGamesCommand,
		ListClaimsCommand,
		ListCreditsCommand,
		GameProposalOutputsCommand,
		CreateGameCommand,
		MoveCommand,
		ResolveCommand,
		ResolveClaimCommand,
		RunTraceCommand,
	}
	app.Action = cliapp.LifecycleCmd(func(ctx *cli.Context, close context.CancelCauseFunc) (cliapp.Lifecycle, error) {
		logger, err := setupLogging(ctx)
		if err != nil {
			return nil, err
		}
		logger.Info("Starting mc-challenger", "version", VersionWithMeta)

		cfg, err := flags.NewConfigFromCLI(ctx, logger)
		if err != nil {
			return nil, err
		}
		return action(ctx.Context, logger, cfg)
	})
	return app.RunContext(ctx, args)
}

func setupLogging(ctx *cli.Context) (log.Logger, error) {
	logCfg := logcli.ReadCLIConfig(ctx)
	logger := logcli.NewLogger(logcli.AppOut(ctx), logCfg)
	logcli.SetGlobalLogHandler(logger.Handler())
	return logger, nil
}
