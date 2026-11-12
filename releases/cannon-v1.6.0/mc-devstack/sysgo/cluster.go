package sysgo

import (
	"github.com/mindchain/mind-chain/mc-devstack/shim"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
	"github.com/mindchain/mind-chain/mc-supervisor/supervisor/backend/depset"
)

type Cluster struct {
	id     stack.ClusterID
	cfgset depset.FullConfigSetMerged
}

func (c *Cluster) hydrate(system stack.ExtensibleSystem) {
	sysCluster := shim.NewCluster(shim.ClusterConfig{
		CommonConfig:  shim.NewCommonConfig(system.T()),
		ID:            c.id,
		DependencySet: c.cfgset.DependencySet,
	})
	system.AddCluster(sysCluster)
}

func (c *Cluster) DepSet() *depset.StaticConfigDependencySet {
	if c.cfgset.DependencySet == nil {
		return nil
	}
	return c.cfgset.DependencySet.(*depset.StaticConfigDependencySet)
}
