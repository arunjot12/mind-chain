package mon

import (
	"testing"
	"time"

	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
	"github.com/mindchain/mind-chain/mc-dispute-mon/metrics"
	"github.com/mindchain/mind-chain/mc-dispute-mon/mon/types"
	"github.com/mindchain/mind-chain/mc-service/clock"
	"github.com/mindchain/mind-chain/mc-service/log"
	"github.com/mindchain/mind-chain/mc-service/testlog"
	"github.com/stretchr/testify/require"
)

func TestResolutionMonitor_CheckResolutions(t *testing.T) {
	r, cl, m := newTestResolutionMonitor(t)
	games := newTestGames(uint64(cl.Now().Unix()))
	r.CheckResolutions(games)

	require.Equal(t, 1, m.calls[metrics.CompleteMaxDuration])
	require.Equal(t, 1, m.calls[metrics.CompleteBeforeMaxDuration])
	require.Equal(t, 1, m.calls[metrics.ResolvableMaxDuration])
	require.Equal(t, 1, m.calls[metrics.ResolvableBeforeMaxDuration])
	require.Equal(t, 1, m.calls[metrics.InProgressMaxDuration])
	require.Equal(t, 1, m.calls[metrics.InProgressBeforeMaxDuration])
}

func newTestResolutionMonitor(t *testing.T) (*ResolutionMonitor, *clock.DeterministicClock, *stubResolutionMetrics) {
	logger := testlog.Logger(t, log.LevelInfo)
	cl := clock.NewDeterministicClock(time.Unix(int64(time.Hour.Seconds()), 0))
	metrics := &stubResolutionMetrics{}
	return NewResolutionMonitor(logger, metrics, cl), cl, metrics
}

type stubResolutionMetrics struct {
	calls map[metrics.ResolutionStatus]int
}

func (s *stubResolutionMetrics) RecordGameResolutionStatus(status metrics.ResolutionStatus, count int) {
	if s.calls == nil {
		s.calls = make(map[metrics.ResolutionStatus]int)
	}
	s.calls[status] += count
}

func newTestGames(duration uint64) []*types.FaultGameData {
	newTestGame := func(duration uint64, status gameTypes.GameStatus, resolvable bool) *types.FaultGameData {
		game := &types.FaultGameData{
			MaxClockDuration: duration,
			CommonGameData:   types.CommonGameData{Status: status},
		}
		if !resolvable {
			game.Claims = []types.EnrichedClaim{
				{
					Resolved: false,
				},
			}
		}
		return game
	}
	return []*types.FaultGameData{
		newTestGame(duration/2, gameTypes.GameStatusInProgress, false),
		newTestGame(duration*5, gameTypes.GameStatusInProgress, false),
		newTestGame(duration/2, gameTypes.GameStatusInProgress, true),
		newTestGame(duration*5, gameTypes.GameStatusInProgress, true),
		newTestGame(duration/2, gameTypes.GameStatusDefenderWon, false),
		newTestGame(duration*5, gameTypes.GameStatusChallengerWon, false),
	}
}
