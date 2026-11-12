package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/mindchain/mind-chain/mc-batcher/batcher"
	"github.com/mindchain/mind-chain/mc-batcher/flags"
	"github.com/mindchain/mind-chain/mc-batcher/metrics"
	opservice "github.com/mindchain/mind-chain/mc-service"
	"github.com/mindchain/mind-chain/mc-service/cliapp"
	"github.com/mindchain/mind-chain/mc-service/ctxinterrupt"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/log/logcli"
	"github.com/mindchain/mind-chain/mc-service/metrics/doc"
)

var (
	Version   = "v0.0.0"
	GitCommit = ""
	GitDate   = ""
)

func main() {
	logcli.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Version = opservice.FormatVersion(Version, GitCommit, GitDate, "")
	app.Name = "mc-batcher"
	app.Usage = "Batch Submitter Service"
	app.Description = "Service for generating and submitting L2 tx batches to L1"
	app.Action = cliapp.LifecycleCmd(batcher.Main(Version))
	app.Commands = []*cli.Command{
		{
			Name:        "doc",
			Subcommands: doc.NewSubcommands(metrics.NewMetrics("default")),
		},
	}

	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}
