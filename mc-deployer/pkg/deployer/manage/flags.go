package manage

import (
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer"
	"github.com/mindchain/mind-chain/mc-deployer/pkg/deployer/upgrade"
	"github.com/mindchain/mind-chain/mc-service/log/logcli"
	"github.com/urfave/cli/v2"
)

var Commands = cli.Commands{
	&cli.Command{
		Name:  "add-game-type-v2",
		Usage: "allows to add new game types to the chain using the OPContractsManager V2",
		Flags: append([]cli.Flag{
			deployer.L1RPCURLFlag,
			upgrade.ConfigFlag,
			upgrade.OverrideArtifactsURLFlag,
			upgrade.OutfileFlag,
			deployer.CacheDirFlag,
		}, logcli.CLIFlags(deployer.EnvVarPrefix)...),
		Action: AddGameTypeOPCMV2CLI,
	},
}
