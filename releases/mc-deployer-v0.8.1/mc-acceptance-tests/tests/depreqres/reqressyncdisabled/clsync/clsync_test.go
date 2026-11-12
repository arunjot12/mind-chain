package clsync

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-acceptance-tests/tests/depreqres/common"
	"github.com/mindchain/mind-chain/mc-node/rollup/sync"
)

func TestUnsafeChainNotStalling_CLSync(gt *testing.T) {
	common.UnsafeChainNotStalling_Disconnect(gt, sync.CLSync, 20, common.ReqRespSyncDisabledOpts(sync.CLSync)...)
}

func TestUnsafeChainNotStalling_CLSync_RestartOpNode(gt *testing.T) {
	common.UnsafeChainNotStalling_RestartOpNode(gt, sync.CLSync, 20, common.ReqRespSyncDisabledOpts(sync.CLSync)...)
}
