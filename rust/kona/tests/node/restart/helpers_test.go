package node_restart

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	node_utils "github.com/mindchain/mind-chain/rust/kona/tests/node/utils"
)

func newRestartPreset(t devtest.T) *node_utils.MixedOpKonaPreset {
	if sharedRestartRuntime != nil {
		return node_utils.NewMixedOpKonaFromRuntime(t, sharedRestartRuntime)
	}

	// Restart tests currently target a minimal kona-only topology.
	return node_utils.NewMixedOpKonaForConfig(t, node_utils.L2NodeConfig{
		KonaSequencerNodesWithReth: 1,
		KonaNodesWithReth:          1,
	})
}
