package sync

import (
	"testing"

	"github.com/mindchain/mind-chain/mc-devstack/presets"
)

func TestMain(m *testing.M) {
	presets.DoMain(m, presets.WithSimpleInterop())
}
