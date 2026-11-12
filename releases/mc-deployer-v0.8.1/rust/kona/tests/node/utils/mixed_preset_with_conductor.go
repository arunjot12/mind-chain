package node_utils

import (
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

type MinimalWithConductors = presets.MinimalWithConductors

func NewMixedOpKonaWithConductors(t devtest.T) *MinimalWithConductors {
	return presets.NewMinimalWithConductors(t)
}

func NewMixedOpKonaWithConductorsForConfig(t devtest.T, _ L2NodeConfig, opts ...presets.Option) *MinimalWithConductors {
	return presets.NewMinimalWithConductors(t, opts...)
}
