package contracts

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	faultTypes "github.com/mindchain/mind-chain/mc-challenger/game/fault/types"
	"github.com/mindchain/mind-chain/mc-service/sources/batching"
	"github.com/mindchain/mind-chain/mc-service/sources/batching/rpcblock"
)

var (
	methodGameType = "gameType"

	gameTypeABI = mustParseAbi([]byte(`[{
		"inputs": [],
		"name": "gameType",
		"outputs": [{"type": "uint32"}],
		"stateMutability": "view",
		"type": "function"
	}]`))
)

func DetectGameType(ctx context.Context, addr common.Address, caller *batching.MultiCaller) (faultTypes.GameType, error) {
	result, err := caller.SingleCall(ctx, rpcblock.Latest, batching.NewContractCall(gameTypeABI, addr, methodGameType))
	if err != nil {
		return faultTypes.UnknownGameType, fmt.Errorf("failed to detect game type: %w", err)
	}
	gameType := faultTypes.GameType(result.GetUint32(0))
	switch gameType {
	case faultTypes.CannonGameType,
		faultTypes.PermissionedGameType,
		faultTypes.CannonKonaGameType,
		faultTypes.AsteriscGameType,
		faultTypes.AlphabetGameType,
		faultTypes.FastGameType,
		faultTypes.AsteriscKonaGameType,
		faultTypes.SuperCannonGameType,
		faultTypes.SuperPermissionedGameType,
		faultTypes.SuperCannonKonaGameType,
		faultTypes.SuperAsteriscKonaGameType:
		return gameType, nil
	default:
		return faultTypes.UnknownGameType, fmt.Errorf("unsupported game type: %d", gameType)
	}
}
