package opnode

import (
	"context"

	"github.com/mindchain/mind-chain/mc-e2e/e2eutils/services"
	"github.com/mindchain/mind-chain/mc-node/config"
	"github.com/mindchain/mind-chain/mc-node/metrics"
	rollupNode "github.com/mindchain/mind-chain/mc-node/node"
	"github.com/mindchain/mind-chain/mc-node/node/runcfg"
	"github.com/mindchain/mind-chain/mc-node/p2p"
	"github.com/mindchain/mind-chain/mc-service/cliapp"
	"github.com/mindchain/mind-chain/mc-service/clock"
	"github.com/mindchain/mind-chain/mc-service/endpoint"
	"github.com/mindchain/mind-chain/mc-service/log"
)

type Opnode struct {
	node *rollupNode.OpNode
}

func (o *Opnode) UserRPC() endpoint.RPC {
	return endpoint.HttpURL(o.node.HTTPEndpoint())
}

func (o *Opnode) UserRPCPort() (int, error) {
	return o.node.HTTPPort()
}

func (o *Opnode) Stop(ctx context.Context) error {
	return o.node.Stop(ctx)
}

func (o *Opnode) Stopped() bool {
	return o.node.Stopped()
}

func (o *Opnode) RuntimeConfig() runcfg.ReadonlyRuntimeConfig {
	return o.node.RuntimeConfig()
}

func (o *Opnode) P2P() p2p.Node {
	return o.node.P2P()
}

var _ services.RollupNode = (*Opnode)(nil)

func NewOpnode(l log.Logger, c *config.Config, clk clock.Clock, errFn func(error)) (*Opnode, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	var cycle cliapp.Lifecycle
	c.Cancel = func(errCause error) {
		l.Warn("node requested early shutdown!", "err", errCause)
		go func() {
			postCtx, postCancel := context.WithCancel(context.Background())
			postCancel() // don't allow the stopping to continue for longer than needed
			if err := cycle.Stop(postCtx); err != nil {
				// Report the original cause that triggered shutdown, not the
				// stop error. The stop context is pre-cancelled, so Stop may
				// return "context canceled" from StopSequencer — that is a
				// consequence of force-quit, not an additional failure.
				// errFn must be goroutine-safe (use t.Error, not require/FailNow).
				errFn(errCause)
			}
			l.Warn("closed mc-node!")
		}()
	}
	node, err := rollupNode.New(context.Background(), c, l, "", metrics.NewMetrics("", nil), clk)
	if err != nil {
		return nil, err
	}
	cycle = node
	err = node.Start(context.Background())
	if err != nil {
		return nil, err
	}
	return &Opnode{node: node}, nil
}
