package frontend

import (
	"errors"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/mindchain/mind-chain/mc-service/jsonrpc"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/seqtypes"
)

// toJsonError turns the error into a JSON error with error-code,
// to preserve the code when the error is wrapped
func toJsonError(err error) error {
	if err == nil {
		return nil
	}
	var x rpc.Error
	if errors.As(err, &x) {
		return x
	}
	return &jsonrpc.Error{
		Code:    seqtypes.ErrUnknownKind.Code,
		Message: err.Error(),
	}
}
