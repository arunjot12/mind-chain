package main

import (
	"context"
	"io"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/log"

	"github.com/mindchain/mind-chain/mc-faucet/config"
	"github.com/mindchain/mind-chain/mc-faucet/faucet"
	"github.com/mindchain/mind-chain/mc-faucet/flags"
	"github.com/mindchain/mind-chain/mc-faucet/metrics"
	opservice "github.com/mindchain/mind-chain/mc-service"
	"github.com/mindchain/mind-chain/mc-service/cliapp"
	"github.com/mindchain/mind-chain/mc-service/ctxinterrupt"
	oplog "github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/metrics/doc"
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

func run(ctx context.Context, w io.Writer, ew io.Writer, args []string, fn faucet.MainFn) error {
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Writer = w
	app.ErrWriter = ew
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Version = opservice.FormatVersion(Version, GitCommit, GitDate, "")
	app.Name = "mc-faucet"
	app.Usage = "mc-faucet hosts configurable faucets to send test-ETH with."
	app.Description = "Faucet service for devnets.\n" +
		" Try the faucet RPC on /chain/{CHAIN_ID_HERE} or /faucet/{FAUCET_NAME_HERE}."
	app.Action = cliapp.LifecycleCmd(faucet.Main(app.Version, fn))
	app.Commands = []*cli.Command{
		{
			Name:        "doc",
			Subcommands: doc.NewSubcommands(metrics.NewMetrics("default")),
		},
	}
	return app.RunContext(ctx, args)
}

func fromConfig(ctx context.Context, cfg *config.Config, logger log.Logger) (cliapp.Lifecycle, error) {
	return faucet.FromConfig(ctx, cfg, logger)
}
