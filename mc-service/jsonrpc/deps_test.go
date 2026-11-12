package jsonrpc

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-service/testutils/depguard"
)

// TestLeafPackage keeps mc-service/jsonrpc free of in-repo imports, so RPC clients, servers and
// metrics can share its types without an import cycle, and of go-ethereum imports, so the types
// stay independent of any RPC implementation.
func TestLeafPackage(t *testing.T) {
	const pkg = "github.com/mindchain/mind-chain/mc-service/jsonrpc"
	depguard.RequireNoTransitiveImportUnder(t, pkg, "github.com/mindchain/mind-chain")
	depguard.RequireNoTransitiveImportUnder(t, pkg, "github.com/ethereum/go-ethereum")
}
