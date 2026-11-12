package interopgen

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
	"github.com/mindchain/mind-chain/mc-core/devfeatures"
	"github.com/stretchr/testify/require"
)

func TestInitialDisputeAbsolutePrestates(t *testing.T) {
	cannonPrestate := common.HexToHash("0x01")
	cannonKonaPrestate := common.HexToHash("0x02")
	tests := []struct {
		name             string
		gameType         gameTypes.GameType
		selectedPrestate common.Hash
	}{
		{
			name:             "permissioned cannon selects cannon prestate",
			gameType:         gameTypes.PermissionedGameType,
			selectedPrestate: cannonPrestate,
		},
		{
			name:             "cannon kona selects kona prestate",
			gameType:         gameTypes.CannonKonaGameType,
			selectedPrestate: cannonKonaPrestate,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &L2Config{
				DisputeGameType:             uint32(test.gameType),
				DisputeAbsolutePrestate:     cannonPrestate,
				DisputeKonaAbsolutePrestate: cannonKonaPrestate,
			}
			selectedPrestate, fallbackPrestate := initialDisputeAbsolutePrestates(cfg)
			require.Equal(t, test.selectedPrestate, selectedPrestate)
			require.Equal(t, cannonPrestate, fallbackPrestate)
			if test.gameType == gameTypes.CannonKonaGameType {
				require.NotEqual(t, selectedPrestate, fallbackPrestate)
			}
		})
	}
}

func TestInteropAtGenesis(t *testing.T) {
	zero := hexutil.Uint64(0)
	nonzero := hexutil.Uint64(24)

	tests := []struct {
		name   string
		offset *hexutil.Uint64
		want   bool
	}{
		{"nil offset: Lagoon not scheduled", nil, false},
		{"zero offset: Lagoon activates interop at genesis", &zero, true},
		{"non-zero offset: Lagoon delayed activation", &nonzero, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, lagoonAtGenesis(tt.offset))
		})
	}
}

// devFeatureBitmapForL2Genesis sets the OptimismPortalInteropFlag when interop is enabled.
func TestDevFeatureBitmapForL2Genesis(t *testing.T) {
	interopOnly := devfeatures.EnableDevFeature(common.Hash{}, devfeatures.OptimismPortalInteropFlag)

	tests := []struct {
		name          string
		enableInterop bool
		want          common.Hash
	}{
		{"interop disabled", false, common.Hash{}},
		{"interop enabled", true, interopOnly},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, devFeatureBitmapForL2Genesis(tt.enableInterop))
		})
	}
}
