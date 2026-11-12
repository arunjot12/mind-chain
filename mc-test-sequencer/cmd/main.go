package main

import (
	"context"
	"io"
	"os"

	"github.com/urfave/cli/v2"

	opservice "github.com/mindchain/mind-chain/mc-service"
	"github.com/mindchain/mind-chain/mc-service/cliapp"
	"github.com/mindchain/mind-chain/mc-service/ctxinterrupt"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/log/logcli"
	"github.com/mindchain/mind-chain/mc-service/metrics/doc"
	"github.com/mindchain/mind-chain/mc-test-sequencer/config"
	"github.com/mindchain/mind-chain/mc-test-sequencer/flags"
	"github.com/mindchain/mind-chain/mc-test-sequencer/metrics"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer"
)

var (
	Version   = "v0.0.0"
	GitCommit = ""
	GitDate   = ""
)

func main() {
	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := run(ctx, os.Stdout, os.Stderr, os.Args, fromConfig)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}

func run(ctx context.Context, w io.Writer, ew io.Writer, args []string, fn sequencer.MainFn) error {
	logcli.SetupDefaults()

	app := cli.NewApp()
	app.Writer = w
	app.ErrWriter = ew
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Version = opservice.FormatVersion(Version, GitCommit, GitDate, "")
	app.Name = "mc-test-sequencer"
	app.Usage = "mc-test-sequencer sequences blocks"
	app.Description = "mc-test-sequencer sequences blocks"
	app.Action = cliapp.LifecycleCmd(sequencer.Main(app.Version, fn))
	app.Commands = []*cli.Command{
		{
			Name:        "doc",
			Subcommands: doc.NewSubcommands(metrics.NewMetrics("default")),
		},
	}
	return app.RunContext(ctx, args)
}

func fromConfig(ctx context.Context, cfg *config.Config, logger log.Logger) (cliapp.Lifecycle, error) {
	return sequencer.FromConfig(ctx, cfg, logger)
}
