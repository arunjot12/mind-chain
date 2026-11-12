package frontend

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/mindchain/mind-chain/mc-service/eth"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/seqtypes"
)

func TestToJsonError(t *testing.T) {
	require.NoError(t, toJsonError(nil))

	err := toJsonError(fmt.Errorf("context: %w", seqtypes.ErrUnknownJob))
	require.Same(t, seqtypes.ErrUnknownJob, err)

	err = toJsonError(fmt.Errorf("context: %w", eth.InputError{Inner: errors.New("bad"), Code: eth.InvalidParams}))
	var rpcErr rpc.Error
	require.ErrorAs(t, err, &rpcErr)
	require.Equal(t, int(eth.InvalidParams), rpcErr.ErrorCode())

	err = toJsonError(errors.New("plain"))
	require.ErrorAs(t, err, &rpcErr)
	require.Equal(t, seqtypes.ErrUnknownKind.Code, rpcErr.ErrorCode())
	require.Equal(t, "plain", rpcErr.Error())
}
