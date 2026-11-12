package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/mindchain/mind-chain/mc-conductor/conductor"
	"github.com/mindchain/mind-chain/mc-conductor/flags"
	opservice "github.com/mindchain/mind-chain/mc-service"
	"github.com/mindchain/mind-chain/mc-service/cliapp"
	"github.com/mindchain/mind-chain/mc-service/ctxinterrupt"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/log/logcli"
)

var (
	Version   = "v0.0.1"
	GitCommit = ""
	GitDate   = ""
)

func main() {
	logcli.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Version = opservice.FormatVersion(Version, GitCommit, GitDate, "")
	app.Name = "mc-conductor"
	app.Usage = "Mind Chain Sequencer Conductor Service"
	app.Description = "mc-conductor help sequencer to run in highly available mode"
	app.Action = cliapp.LifecycleCmd(OpConductorMain)
	app.Commands = []*cli.Command{}

	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}

func OpConductorMain(ctx *cli.Context, closeApp context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	logCfg := logcli.ReadCLIConfig(ctx)
	log := logcli.NewLogger(logcli.AppOut(ctx), logCfg)
	logcli.SetGlobalLogHandler(log.Handler())
	opservice.ValidateEnvVars(flags.EnvVarPrefix, flags.Flags, log)

	cfg, err := conductor.NewConfig(ctx, log)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	c, err := conductor.New(ctx.Context, cfg, log, Version)
	if err != nil {
		return nil, fmt.Errorf("failed to create conductor: %w", err)
	}

	return c, nil
}
