package standard

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

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

// TestStandardSuperchainUnavailable pins that Mind Chain ships no upstream
// Superchain values: every standard lookup fails with ErrNoStandardSuperchain.
func TestStandardSuperchainUnavailable(t *testing.T) {
	for _, chainID := range []uint64{1, 11155111} {
		for name, f := range map[string]func(uint64) (common.Address, error){
			"GuardianAddressFor":   GuardianAddressFor,
			"ChallengerAddressFor": ChallengerAddressFor,
			"L1ProxyAdminOwner":    L1ProxyAdminOwner,
			"L2ProxyAdminOwner":    L2ProxyAdminOwner,
		} {
			_, err := f(chainID)
			require.ErrorIs(t, err, ErrNoStandardSuperchain, name)
		}
		_, err := OPCMImplAddressFor(chainID, CurrentTag)
		require.ErrorIs(t, err, ErrNoStandardSuperchain)
	}
}

func TestSP1VerifierFor(t *testing.T) {
	for _, chainID := range []uint64{1, 11155111} {
		addr, err := SP1VerifierFor(chainID)
		require.NoError(t, err)
		require.Equal(t, common.HexToAddress("0xc3c6dDDAc8829b233Dc6536Ec024775a57b0AF2A"), addr)
	}
}
