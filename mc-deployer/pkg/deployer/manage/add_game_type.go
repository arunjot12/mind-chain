package manage

import (
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/upgrade"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/upgrade/embedded"
	"github.com/urfave/cli/v2"
)

// AddGameTypeOPCMV2CLI is the CLI command for adding a new game type to the chain using the OPContractsManager V2.
// This command is an alias for the upgrade command with the default upgrader.
func AddGameTypeOPCMV2CLI(cliCtx *cli.Context) error {
	return upgrade.UpgradeCLI(embedded.DefaultUpgrader)(cliCtx)
}
