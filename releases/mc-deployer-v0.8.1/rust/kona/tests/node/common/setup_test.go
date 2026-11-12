package node

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	node_utils "github.com/mindchain/mind-chain/rust/kona/tests/node/utils"
)

func newCommonPreset(t devtest.T) *node_utils.MixedOpKonaPreset {
	t.Helper()
	return node_utils.NewMixedOpKona(t)
}
