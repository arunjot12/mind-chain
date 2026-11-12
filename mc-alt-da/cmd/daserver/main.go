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
)

var Version = "v0.0.0"

func main() {
	logcli.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(Flags)
	app.Version = opservice.FormatVersion(Version, "", "", "")
	app.Name = "da-server"
	app.Usage = "AltDA Storage Service"
	app.Description = "Service for storing AltDA inputs"
	app.Action = StartDAServer

	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}

}
