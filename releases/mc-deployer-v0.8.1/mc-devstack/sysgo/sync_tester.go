package sysgo

import (
	"github.com/mindchain/mind-chain/mc-service/eth"
	"github.com/mindchain/mind-chain/mc-sync-tester/synctester"
)

// Caveat: id is binded by a single EL(chainID), but service can support multiple ELs
type SyncTesterService struct {
	service *synctester.Service
}

func (n *SyncTesterService) DefaultEndpoint(chainID eth.ChainID) (string, string, bool) {
	if n == nil || n.service == nil {
		return "", "", false
	}
	for syncTesterID, mappedChainID := range n.service.SyncTesters() {
		if mappedChainID != chainID {
			continue
		}
		return syncTesterID.String(), n.service.SyncTesterRPC(chainID, false), true
	}
	return "", "", false
}

func (n *SyncTesterService) RPC() string {
	if n == nil || n.service == nil {
		return ""
	}
	return n.service.RPC()
}

func (n *SyncTesterService) SyncTesterRPCPath(chainID eth.ChainID, withSessionID bool) string {
	if n == nil || n.service == nil {
		return ""
	}
	return n.service.SyncTesterRPCPath(chainID, withSessionID)
}
