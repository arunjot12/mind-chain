package services

import (
	"context"
	"fmt"
	"os"

	rollupNode "github.com/mindchain/mind-chain/mc-node/node/runcfg"
	"github.com/mindchain/mind-chain/mc-node/p2p"
	"github.com/mindchain/mind-chain/mc-service/endpoint"
)

// ELKind selects which execution-layer backend mc-e2e runs as the L2 EL.
type ELKind string

const (
	ELKindOpGeth ELKind = "op-geth"
	ELKindOpReth ELKind = "mc-reth"
)

// DefaultELKind returns the L2 EL backend selected by the OP_E2E_L2_EL_KIND
// environment variable, defaulting to mc-reth. A non-empty but unrecognized
// value panics rather than silently misconfiguring the suite.
func DefaultELKind() ELKind {
	switch v := os.Getenv("OP_E2E_L2_EL_KIND"); v {
	case "":
		return ELKindOpReth
	case string(ELKindOpGeth):
		return ELKindOpGeth
	case string(ELKindOpReth):
		return ELKindOpReth
	default:
		panic(fmt.Sprintf("unrecognized OP_E2E_L2_EL_KIND %q (want %q or %q)", v, ELKindOpGeth, ELKindOpReth))
	}
}

// EthInstance is either an in process Geth or external process exposing its
// endpoints over the network
type EthInstance interface {
	UserRPC() endpoint.RPC
	AuthRPC() endpoint.RPC
	Close() error
}

type RollupNode interface {
	UserRPC() endpoint.RPC
	Stop(ctx context.Context) error
	Stopped() bool
	RuntimeConfig() rollupNode.ReadonlyRuntimeConfig
	P2P() p2p.Node
}
