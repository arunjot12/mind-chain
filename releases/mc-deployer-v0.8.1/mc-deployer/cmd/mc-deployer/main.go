package main

import (
	"fmt"
	"os"

	"github.com/mindchain/mind-chain/mc-deployer/pkg/cli"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/artifacts"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/version"
)

func main() {
	// Ensure cleanup happens on exit
	defer func() {
		_ = artifacts.CleanupAll()
	}()

	app := cli.NewApp(version.VersionWithMeta)
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr
	err := app.Run(os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Application failed: %v\n", err)
		os.Exit(1)
	}
}
