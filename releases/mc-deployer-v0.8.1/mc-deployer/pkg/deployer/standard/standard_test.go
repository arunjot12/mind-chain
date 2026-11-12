package standard

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultHardforkSchedule(t *testing.T) {
	sched := DefaultHardforkSchedule()
	require.NotSame(t, sched, DefaultHardforkSchedule())
	require.NotNil(t, sched.RegolithTime(0))
	require.NotNil(t, sched.CanyonTime(0))
	require.NotNil(t, sched.DeltaTime(0))
	require.NotNil(t, sched.EcotoneTime(0))
	require.NotNil(t, sched.FjordTime(0))
	require.NotNil(t, sched.GraniteTime(0))
	require.NotNil(t, sched.HoloceneTime(0))
	require.NotNil(t, sched.IsthmusTime(0))
	require.NotNil(t, sched.JovianTime(0))
	require.NotNil(t, sched.KarstTime(0))
	require.Nil(t, sched.LagoonTime(0))
}

// TestSP1VerifierForUnsupportedChain pins that networks outside the OPCM release require an
// explicitly supplied verifier.
func TestSP1VerifierForUnsupportedChain(t *testing.T) {
	_, err := SP1VerifierFor(900)
	require.ErrorContains(t, err, "unsupported chain ID: 900")
}
