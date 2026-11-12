package dsl

import (
	"github.com/mindchain/mind-chain/mc-devstack/stack"
)

type TestSequencer struct {
	commonImpl

	inner stack.TestSequencer
}

func NewTestSequencer(inner stack.TestSequencer) *TestSequencer {
	return &TestSequencer{
		commonImpl: commonFromT(inner.T()),
		inner:      inner,
	}
}

func (s *TestSequencer) String() string {
	return s.inner.ID().String()
}

func (s *TestSequencer) Escape() stack.TestSequencer {
	return s.inner
}
