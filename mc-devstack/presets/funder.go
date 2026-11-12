package presets

import (
	"github.com/mindchain/mind-chain/mc-chain-ops/devkeys"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
)

// newFunderEOA builds a prefunded funder for a chain. It replaces the former
// mc-faucet service + dsl.Funder pair: the funder's key is the genesis-prefunded
// funder account (sysgo.FunderKey), and it mints/funds child EOAs from the given
// test wallet so that funded accounts stay isolated per test.
func newFunderEOA(t devtest.T, keys devkeys.Keys, el dsl.ELNode, wallet *dsl.HDWallet) *dsl.FunderEOA {
	fkey, err := sysgo.FunderKey(keys)
	t.Require().NoError(err, "must derive funder key")
	return dsl.NewFunderEOA(dsl.NewEOA(dsl.NewKey(t, fkey), el), wallet)
}
