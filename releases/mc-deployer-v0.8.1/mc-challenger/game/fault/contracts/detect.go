package contracts

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
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

func DetectGameType(ctx context.Context, addr common.Address, caller *batching.MultiCaller) (gameTypes.GameType, error) {
	result, err := caller.SingleCall(ctx, rpcblock.Latest, batching.NewContractCall(gameTypeABI, addr, methodGameType))
	if err != nil {
		return gameTypes.UnknownGameType, fmt.Errorf("failed to detect game type: %w", err)
	}
	gameType := gameTypes.GameType(result.GetUint32(0))
	switch gameType {
	case gameTypes.CannonGameType,
		gameTypes.PermissionedGameType,
		gameTypes.CannonKonaGameType,
		gameTypes.AlphabetGameType,
		gameTypes.FastGameType,
		gameTypes.SuperPermissionedGameType,
		gameTypes.SuperCannonKonaGameType:
		return gameType, nil
	default:
		return gameTypes.UnknownGameType, fmt.Errorf("unsupported game type: %d", gameType)
	}
}
