package zk

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/mindchain/mind-chain/mc-challenger/config"
	"github.com/mindchain/mind-chain/mc-challenger/game/client"
	"github.com/mindchain/mind-chain/mc-challenger/game/fault/claims"
	"github.com/mindchain/mind-chain/mc-challenger/game/fault/contracts"
	"github.com/mindchain/mind-chain/mc-challenger/game/generic"
	"github.com/mindchain/mind-chain/mc-challenger/game/scheduler"
	gameTypes "github.com/mindchain/mind-chain/mc-challenger/game/types"
	"github.com/mindchain/mind-chain/mc-challenger/metrics"
	"github.com/mindchain/mind-chain/mc-service/sources/batching"
	"github.com/mindchain/mind-chain/mc-service/txmgr"
)

type ClockReader interface {
	Now() time.Time
}

type Registry interface {
	RegisterGameType(gameType gameTypes.GameType, creator scheduler.PlayerCreator)
	RegisterBondContract(gameType gameTypes.GameType, creator claims.BondContractCreator)
}

func RegisterBondContracts(m metrics.Metricer, registry Registry, caller *batching.MultiCaller) {
	registry.RegisterBondContract(gameTypes.ZKDisputeGameType, func(game gameTypes.GameMetadata) (claims.BondContract, error) {
		return contracts.NewZKDisputeGameContract(m, game.Proxy, caller)
	})
}

type TxSender interface {
	SendAndWaitSimple(txPurpose string, txs ...txmgr.TxCandidate) error
}

func RegisterGameTypes(
	ctx context.Context,
	l1Clock ClockReader,
	logger log.Logger,
	m metrics.Metricer,
	cfg *config.Config,
	registry Registry,
	txSender TxSender,
	clients *client.Provider,
	gameStatusProvider GameStatusProvider,
) error {
	if cfg.GameTypeEnabled(gameTypes.ZKDisputeGameType) {
		registry.RegisterGameType(gameTypes.ZKDisputeGameType, func(game gameTypes.GameMetadata, dir string) (scheduler.GamePlayer, error) {
			superNodeClient, syncValidator, err := clients.SuperchainClients()
			if err != nil {
				return nil, fmt.Errorf("failed to create superchain clients: %w", err)
			}
			contract, err := contracts.NewZKDisputeGameContract(m, game.Proxy, clients.MultiCaller())
			if err != nil {
				return nil, fmt.Errorf("failed to create zk dispute game bindings: %w", err)
			}
			return generic.NewGenericGamePlayer(
				ctx,
				logger,
				game.Proxy,
				contract,
				syncValidator,
				nil,
				clients.L1Client(),
				ActorCreator(l1Clock, superNodeClient, gameStatusProvider, contract, txSender),
			)
		})
	}
	return nil
}
