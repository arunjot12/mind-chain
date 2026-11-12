package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/mindchain/mind-chain/mc-deployer/pkg/cli"
)

// RunCLI provides a testable interface for running the mc-deployer CLI
func RunCLI(ctx context.Context, w io.Writer, ew io.Writer, args []string) error {
	app := cli.NewApp("v0.0.0-test")
	app.Writer = w
	app.ErrWriter = ew
	err := app.RunContext(ctx, args)
	if err != nil {
		_, _ = fmt.Fprintf(ew, "Application failed: %v\n", err)
	}
	return err
}
