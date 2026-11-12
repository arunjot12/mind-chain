package extract

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	monTypes "github.com/mindchain/mind-chain/mc-dispute-mon/mon/types"
	"github.com/mindchain/mind-chain/mc-service/clock"
	"github.com/mindchain/mind-chain/mc-service/sources/batching/rpcblock"
)

// ZKAgreementEnricher reconciles a ZK proposal with configured super-root sources.
type ZKAgreementEnricher struct {
	superAgreement SuperAgreementEnricher
}

var _ ZKEnricher = (*ZKAgreementEnricher)(nil)

func NewZKAgreementEnricher(logger log.Logger, metrics OutputMetrics, clients []SuperRootProvider, cl clock.Clock) *ZKAgreementEnricher {
	return &ZKAgreementEnricher{
		superAgreement: *NewSuperAgreementEnricher(logger, metrics, clients, cl),
	}
}

func (e *ZKAgreementEnricher) Enrich(ctx context.Context, _ rpcblock.Block, _ ZKGameCaller, zkGame *monTypes.ZKGameData) error {
	return e.superAgreement.enrich(ctx, zkGame.Common(), zkGameAgreement)
}
