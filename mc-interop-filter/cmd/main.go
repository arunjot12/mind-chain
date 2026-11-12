package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"

	opservice "github.com/mindchain/mind-chain/mc-service"
	"github.com/mindchain/mind-chain/mc-service/cliapp"
	"github.com/mindchain/mind-chain/mc-service/ctxinterrupt"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/log/logcli"
	"github.com/mindchain/mind-chain/mc-service/metrics/doc"

	"github.com/mindchain/mind-chain/mc-interop-filter/filter"
	"github.com/mindchain/mind-chain/mc-interop-filter/flags"
	"github.com/mindchain/mind-chain/mc-interop-filter/metrics"
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
	app.Name = "mc-interop-filter"
	app.Usage = "Interop transaction filter service"
	app.Description = "Validates interop executing messages for transaction filtering"
	app.Action = cliapp.LifecycleCmd(filter.Main(app.Version))
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
