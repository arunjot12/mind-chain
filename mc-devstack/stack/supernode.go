package stack

import (
	"github.com/mindchain/mind-chain/mc-service/apis"
	"github.com/mindchain/mind-chain/mc-service/client"
	"github.com/mindchain/mind-chain/mc-supernode/supernode/activity/interop"
)

type Supernode interface {
	Common
	ClientRPC() client.RPC
	QueryAPI() apis.SupernodeQueryAPI
	UserRPC() string
}

// SupernodeTestControl is the integration-test surface on a running
// supernode. See mc-supernode/supernode/activity/interop for the methods
// available on the InteropActivity pointer.
type SupernodeTestControl interface {
	// InteropActivity returns the current interop activity, or nil if the
	// supernode is stopped or interop is not configured. Do not cache the
	// pointer across RestartWithFreshDataDir.
	InteropActivity() *interop.Interop

	// RestartWithFreshDataDir stops the supernode, deletes its on-disk
	// data directory, and starts a fresh supernode against the same chain
	// containers, virtual nodes, and externally-visible RPC address.
	RestartWithFreshDataDir() error

	// Stop halts the supernode while preserving its data directory and RPC
	// address; Start brings it back up. Used by sync tests that need to halt
	// the verifier, mutate external state, and resume.
	Stop()
	Start()
}
