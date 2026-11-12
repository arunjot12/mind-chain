package main

import (
	"context"
	"os"

	"github.com/mindchain/mind-chain/mc-service/ctxinterrupt"

	opservice "github.com/mindchain/mind-chain/mc-service"
	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/log"
	"github.com/mindchain/mind-chain/mc-proposer/flags"
	"github.com/mindchain/mind-chain/mc-proposer/metrics"
	"github.com/mindchain/mind-chain/mc-proposer/proposer"
	"github.com/mindchain/mind-chain/mc-service/cliapp"
	oplog "github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/metrics/doc"
)

var (
	Version   = "v0.0.0"
	GitCommit = ""
	GitDate   = ""
)

func main() {
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Version = opservice.FormatVersion(Version, GitCommit, GitDate, "")
	app.Name = "mc-proposer"
	app.Usage = "L2 Output Submitter"
	app.Description = "Service for generating and proposing L2 Outputs"
	app.Action = cliapp.LifecycleCmd(proposer.Main(Version))
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
