package noopsigner

import (
	"github.com/mindchain/mind-chain/mc-service/signer"
	"github.com/mindchain/mind-chain/mc-test-sequencer/sequencer/backend/work"
)

type NoopSignedBlock struct {
	work.Block
}

func (s *NoopSignedBlock) VerifySignature(_ signer.Authenticator) error {
	return nil
}
