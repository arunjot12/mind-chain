package shim

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"

	"github.com/mindchain/mind-chain/mc-chain-ops/devkeys"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
	"github.com/mindchain/mind-chain/mc-service/testreq"
)

type keyringImpl struct {
	keys    devkeys.Keys
	require *testreq.Assertions
}

var _ stack.Keys = (*keyringImpl)(nil)

func NewKeyring(keys devkeys.Keys, req *testreq.Assertions) stack.Keys {
	return &keyringImpl{
		keys:    keys,
		require: req,
	}
}

func (k *keyringImpl) Secret(key devkeys.Key) *ecdsa.PrivateKey {
	pk, err := k.keys.Secret(key)
	k.require.NoError(err)
	return pk
}

func (k *keyringImpl) Address(key devkeys.Key) common.Address {
	addr, err := k.keys.Address(key)
	k.require.NoError(err)
	return addr
}
