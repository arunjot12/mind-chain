package sync

import (
	"testing"

	"github.com/ethereum/go-ethereum/log"

	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-devstack/stack"
	"github.com/mindchain/mind-chain/mc-service/log/logfilter"
)

func TestMain(m *testing.M) {
	presets.DoMain(m, presets.WithMultiSupervisorInterop(),
		presets.WithLogFilter(logfilter.DefaultMute(
			stack.KindSelector(stack.SupervisorKind).And(logfilter.Level(log.LevelInfo)).Show(),
			logfilter.Level(log.LevelError).Show(),
		)))
}
