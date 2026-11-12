package sysgo

import (
	"github.com/mindchain/mind-chain/mc-chain-ops/devkeys"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
)

func WithMnemonicKeys(mnemonic string) stack.Option[*Orchestrator] {
	return stack.BeforeDeploy(func(orch *Orchestrator) {
		require := orch.P().Require()
		hd, err := devkeys.NewMnemonicDevKeys(mnemonic)
		require.NoError(err)
		orch.keys = hd
	})
}
