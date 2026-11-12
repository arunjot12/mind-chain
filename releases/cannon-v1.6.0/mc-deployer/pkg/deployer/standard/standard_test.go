package standard

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultHardforkScheduleForTag(t *testing.T) {
	sched := DefaultHardforkScheduleForTag(ContractsV160Tag)
	require.Nil(t, sched.HoloceneTime(0))
	require.Nil(t, sched.IsthmusTime(0))

	sched = DefaultHardforkScheduleForTag(ContractsV180Tag)
	require.NotNil(t, sched.HoloceneTime(0))
	require.Nil(t, sched.IsthmusTime(0))

	sched = DefaultHardforkScheduleForTag(ContractsV200Tag)
	require.NotNil(t, sched.HoloceneTime(0))
	require.Nil(t, sched.IsthmusTime(0))

	sched = DefaultHardforkScheduleForTag(ContractsV300Tag)
	require.NotNil(t, sched.HoloceneTime(0))
	require.Nil(t, sched.IsthmusTime(0))

	sched = DefaultHardforkScheduleForTag("")
	require.NotNil(t, sched.HoloceneTime(0))
	require.NotNil(t, sched.IsthmusTime(0))
}
