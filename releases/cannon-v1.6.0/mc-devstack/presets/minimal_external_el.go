package presets

import (
	"github.com/ethereum/go-ethereum/log"

	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/dsl"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
	"github.com/mindchain/mind-chain/mc-devstack/sysgo"
)

type MinimalExternalEL struct {
	Log          log.Logger
	T            devtest.T
	ControlPlane stack.ControlPlane

	L1Network *dsl.L1Network
	L1EL      *dsl.L1ELNode

	L2Chain      *dsl.L2Network
	L2CL         *dsl.L2CLNode
	L2EL         *dsl.L2ELNode
	L2ELReadOnly *dsl.L2ELNode

	SyncTester *dsl.SyncTester
}

func (m *MinimalExternalEL) L2Networks() []*dsl.L2Network {
	return []*dsl.L2Network{
		m.L2Chain,
	}
}

func WithExternalELWithSuperchainRegistry(networkPreset stack.ExtNetworkConfig) stack.CommonOption {
	return stack.MakeCommon(sysgo.ExternalELSystemWithEndpointAndSuperchainRegistry(&sysgo.DefaultMinimalExternalELSystemIDs{}, networkPreset))
}
