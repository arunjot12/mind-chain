package log

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-service/testutils/depguard"
)

// TestLeafPackage keeps mc-service/log free of in-repo imports, so that every
// package in the monorepo can log through it without an import cycle.
func TestLeafPackage(t *testing.T) {
	depguard.RequireNoTransitiveImportUnder(t,
		"github.com/mindchain/mind-chain/mc-service/log",
		"github.com/mindchain/mind-chain")
}
