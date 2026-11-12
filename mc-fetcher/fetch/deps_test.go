package fetch

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-service/testutils/depguard"
)

// TestImportClosure keeps mc-fetcher's library packages light, because the
// superchain-registry's ops tooling consumes them as a Go module dependency:
//   - mc-core is off limits: mc-core/superchain embeds the generated, gitignored superchain
//     config bundle, which is absent from a clean module download, and the OP types in
//     mc-core still need op-geth.
//   - mc-deployer and the mc-chain-ops script host need op-geth, and a fetch needs no EVM of
//     its own.
func TestImportClosure(t *testing.T) {
	const repo = "github.com/mindchain/mind-chain"
	for _, pkgs := range []string{repo + "/mc-fetcher/fetch/...", repo + "/mc-fetcher/chaininfo/..."} {
		depguard.RequireNoTransitiveImportUnder(t, pkgs, repo+"/mc-core")
		depguard.RequireNoTransitiveImportUnder(t, pkgs, repo+"/mc-deployer")
		depguard.RequireNoTransitiveImport(t, pkgs,
			repo+"/mc-chain-ops/script",
			repo+"/mc-chain-ops/foundry",
		)
	}
}
