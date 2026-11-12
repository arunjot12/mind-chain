package chaincfg

import (
	"fmt"
	"sort"
	"strings"

	registry "github.com/mindchain/mind-chain/mc-core/superchain"
	"github.com/mindchain/mind-chain/mc-node/rollup"
	"github.com/mindchain/mind-chain/mc-node/superchain"
)

// L2ChainIDToNetworkDisplayName maps chain IDs to built-in chain names. It is built from the
// chain index alone, so no chain config is parsed. Mind Chain ships no built-in chains, so
// it is empty.
var L2ChainIDToNetworkDisplayName = func() map[string]string {
	out := make(map[string]string)
	for chainID, netCfg := range registry.Chains {
		out[fmt.Sprintf("%d", chainID)] = netCfg.Name
	}
	return out
}()

// AvailableNetworks returns the selection of network configurations that is available by default.
func AvailableNetworks() []string {
	var networks []string
	for _, cfg := range registry.Chains {
		networks = append(networks, cfg.Name+"-"+cfg.Network)
	}
	sort.Strings(networks)
	return networks
}

func handleLegacyName(name string) string {
	switch name {
	case "mainnet":
		return "op-mainnet"
	case "sepolia":
		return "op-sepolia"
	default:
		return name
	}
}

// ChainByName returns a chain, from known available configurations, by name.
// ChainByName returns nil when the chain name is unknown.
func ChainByName(name string) *registry.ChainConfig {
	// Handle legacy name aliases
	name = handleLegacyName(name)
	for _, chainCfg := range registry.Chains {
		if !strings.EqualFold(chainCfg.Name+"-"+chainCfg.Network, name) {
			continue
		}

		cfg, err := chainCfg.Config()
		if err != nil {
			panic(fmt.Errorf("failed to load chain config: %w", err))
		}
		return cfg
	}
	return nil
}

func GetRollupConfig(name string) (*rollup.Config, error) {
	chainCfg := ChainByName(name)
	if chainCfg == nil {
		return nil, fmt.Errorf("invalid network: %q", name)
	}
	rollupCfg, err := superchain.LoadOPStackRollupConfig(chainCfg.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to load rollup config: %w", err)
	}
	return rollupCfg, nil
}
