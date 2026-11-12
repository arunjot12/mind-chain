package main

import (
	"fmt"
	"os"

	"github.com/mindchain/mind-chain/mc-deployer/pkg/cli"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/version"

	opservice "github.com/mindchain/mind-chain/mc-service"
)

var (
	GitCommit = ""
	GitDate   = ""
)

// VersionWithMeta holds the textual version string including the metadata.
var VersionWithMeta = opservice.FormatVersion(version.Version, GitCommit, GitDate, version.Meta)

func main() {
	app := cli.NewApp(VersionWithMeta)
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr
	err := app.Run(os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Application failed: %v\n", err)
		os.Exit(1)
	}
}
