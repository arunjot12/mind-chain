package node

import (
	"sync"
	"testing"

	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"

	safety "github.com/mindchain/mind-chain/mc-service/eth/safety"
	node_utils "github.com/mindchain/mind-chain/rust/kona/tests/node/utils"
	"github.com/stretchr/testify/require"
)

func TestEngine(gt *testing.T) {
	t := devtest.ParallelT(gt)

	out := newCommonPreset(t)

	// Get the nodes from the network.
	nodes := out.L2CLKonaNodes()

	wg := sync.WaitGroup{}
	for _, node := range nodes {
		wg.Add(1)
		go func(node *dsl.L2CLNode) {
			defer wg.Done()

			done := make(chan struct{})
			go func() {
				defer close(done)
				// Wait for 40 unsafe blocks to be produced.
				node.Advanced(safety.LocalUnsafe, 40, 100)
			}()

			queueLens := node_utils.GetDevWS(t, node, "engine_queue_size", done)
			for _, queueLen := range queueLens {
				require.LessOrEqual(t, queueLen, uint64(1), "engine queue length should be 1 or less")
			}
		}(&node)
	}

	wg.Wait()

}
