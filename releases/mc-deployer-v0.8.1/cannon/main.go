package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mindchain/mind-chain/mc-service/ctxinterrupt"
	"github.com/urfave/cli/v2"

	"github.com/mindchain/mind-chain/cannon/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = os.Args[0]
	app.Usage = "MIPS Fault Proof tool"
	app.Description = "MIPS Fault Proof tool"
	app.Commands = []*cli.Command{
		cmd.LoadELFCommand,
		cmd.WitnessCommand,
		cmd.RunCommand,
	}
	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		if errors.Is(err, ctx.Err()) {
			_, _ = fmt.Fprintf(os.Stderr, "command interrupted")
			os.Exit(130)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
	}
}
