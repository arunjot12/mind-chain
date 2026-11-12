package disputegame

import (
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethclient"
	opgenesis "github.com/mindchain/mind-chain/mc-core/genesis"
	"github.com/mindchain/mind-chain/mc-core/interop/depset"
	"github.com/mindchain/mind-chain/mc-e2e/interop"
	"github.com/mindchain/mind-chain/mc-e2e/system/e2esys"
	"github.com/mindchain/mind-chain/mc-node/rollup"
	"github.com/mindchain/mind-chain/mc-service/endpoint"
	"github.com/mindchain/mind-chain/mc-service/sources"
)

type SuperDisputeSystem struct {
	sys  interop.SuperSystem
	opts *e2esys.SystemConfigOpts
}

func (s *SuperDisputeSystem) SupernodeClient() *sources.SuperNodeClient {
	return s.sys.SupernodeClient()
}

func NewSuperDisputeSystem(sys interop.SuperSystem, opts *e2esys.SystemConfigOpts) *SuperDisputeSystem {
	return &SuperDisputeSystem{sys, opts}
}

func splitName(name string) (string, string) {
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		panic("Invalid super system name: " + name)
	}
	return parts[0], parts[1]
}

func (s *SuperDisputeSystem) L1BeaconEndpoint() endpoint.RestHTTP {
	beacon := s.sys.L1Beacon()
	return endpoint.RestHTTPURL(beacon.BeaconAddr())
}

func (s *SuperDisputeSystem) NodeEndpoint(name string) endpoint.RPC {
	if name == "l1" {
		return s.sys.L1().UserRPC()
	}
	network, node := splitName(name)
	return s.sys.L2GethEndpoint(network, node)
}

func (s *SuperDisputeSystem) L2NodeEndpoints() []endpoint.RPC {
	networks := s.sys.L2IDs()
	endpoints := make([]endpoint.RPC, len(networks))
	for i, network := range networks {
		endpoints[i] = s.sys.L2GethEndpoint(network, "sequencer")
	}
	return endpoints
}

func (s *SuperDisputeSystem) SupernodeEndpoint() endpoint.RPC {
	return s.sys.SupernodeEndpoint()
}

func (s *SuperDisputeSystem) NodeClient(name string) *ethclient.Client {
	if name == "l1" {
		return s.sys.L1GethClient()
	}
	network, node := splitName(name)
	return s.sys.L2GethClient(network, node)
}

func (s *SuperDisputeSystem) RollupEndpoint(name string) endpoint.RPC {
	network, node := splitName(name)
	return s.sys.L2RollupEndpoint(network, node)
}

func (s *SuperDisputeSystem) RollupClient(name string) *sources.RollupClient {
	network, node := splitName(name)
	return s.sys.L2RollupClient(network, node)
}

func (s *SuperDisputeSystem) DisputeGameFactoryAddr() common.Address {
	return s.sys.DisputeGameFactoryAddr()
}

func (s *SuperDisputeSystem) RollupCfgs() []*rollup.Config {
	networks := s.sys.L2IDs()
	cfgs := make([]*rollup.Config, len(networks))
	for i, network := range networks {
		cfgs[i] = s.sys.RollupConfig(network)
	}
	return cfgs
}

func (s *SuperDisputeSystem) DependencySet() *depset.StaticConfigDependencySet {
	return s.sys.DependencySet()
}

func (s *SuperDisputeSystem) L1Genesis() *core.Genesis {
	return s.sys.L1Genesis()
}

func (s *SuperDisputeSystem) L2Geneses() []*opgenesis.Genesis {
	networks := s.sys.L2IDs()
	cfgs := make([]*opgenesis.Genesis, len(networks))
	for i, network := range networks {
		cfgs[i] = s.sys.L2Genesis(network)
	}
	return cfgs
}

func (s *SuperDisputeSystem) AdvanceTime(duration time.Duration) {
	s.sys.AdvanceL1Time(duration)
}

func (s *SuperDisputeSystem) IsSupersystem() bool {
	return true
}

var _ DisputeSystem = (*SuperDisputeSystem)(nil)
