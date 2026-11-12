package main

import (
	"errors"
	"os"

	opservice "github.com/mindchain/mind-chain/mc-service"
	"github.com/urfave/cli/v2"

	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/log/logcli"
	wheel "github.com/mindchain/mind-chain/mc-wheel"
)

var (
	Version   = ""
	GitCommit = ""
	GitDate   = ""
)

func main() {
	app := cli.NewApp()
	app.Version = opservice.FormatVersion(Version, GitCommit, GitDate, "")
	app.Name = "mc-wheel"
	app.Usage = "Mind Chain Wheel is a CLI tool for the execution engine"
	app.Description = "Mind Chain Wheel is a CLI tool to direct the engine one way or the other with Engine API routines."
	app.Flags = []cli.Flag{wheel.GlobalGethLogLvlFlag}
	app.Before = func(c *cli.Context) error {
		lvl := c.Generic(wheel.GlobalGethLogLvlFlag.Name).(*logcli.LevelFlagValue).Level()
		logcli.SetGlobalLogHandler(log.NewTerminalHandlerWithLevel(os.Stdout, lvl, true))
		return nil
	}
	app.Action = func(c *cli.Context) error {
		return errors.New("see 'engine' subcommand and --help")
	}
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr
	app.Commands = []*cli.Command{
		wheel.EngineCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}
