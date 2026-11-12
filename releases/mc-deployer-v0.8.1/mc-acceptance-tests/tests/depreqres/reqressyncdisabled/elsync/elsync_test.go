package elsync

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-acceptance-tests/tests/depreqres/common"
	"github.com/mindchain/mind-chain/mc-node/rollup/sync"
)

func TestUnsafeChainNotStalling_ELSync(gt *testing.T) {
	common.UnsafeChainNotStalling_Disconnect(gt, sync.ELSync, 20, common.ReqRespSyncDisabledOpts(sync.ELSync)...)
}

func TestUnsafeChainNotStalling_ELSync_RestartOpNode(gt *testing.T) {
	common.UnsafeChainNotStalling_RestartOpNode(gt, sync.ELSync, 20, common.ReqRespSyncDisabledOpts(sync.ELSync)...)
}
