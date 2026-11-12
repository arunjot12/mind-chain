package msg

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-core/interop/messages"
	"github.com/mindchain/mind-chain/mc-core/predeploys"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-service/eth"
	"github.com/mindchain/mind-chain/mc-service/txintent/bindings"
)

func TestDepositCannotValidateMessage(gt *testing.T) {
	t := devtest.ParallelT(gt)
	sys := presets.NewSimpleInterop(t)
	require := t.Require()

	sys.L1Network.WaitForOnline()
	alice := sys.FunderL1.NewFundedEOA(eth.OneEther)
	depositor := alice.AsEL(sys.L2ELA).ViaDepositTx(alice, sys.L2ELA, sys.L2ChainA)

	inbox := bindings.NewBindings[bindings.CrossL2Inbox]()
	calldata, err := inbox.ValidateMessage(messages.Identifier{}, eth.Bytes32{}).EncodeInputLambda()
	require.NoError(err, "failed to encode validateMessage calldata")

	depositor.DepositTxExpectRevert(predeploys.CrossL2InboxAddr, calldata, "CrossL2Inbox_NoExecutingDeposits()")
}
