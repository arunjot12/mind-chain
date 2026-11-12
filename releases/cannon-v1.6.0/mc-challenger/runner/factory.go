package runner

import (
	"context"
	"errors"
	"net/url"

	"github.com/ethereum/go-ethereum/log"
	"github.com/mindchain/mind-chain/mc-challenger/config"
	"github.com/mindchain/mind-chain/mc-challenger/game/fault/trace/asterisc"
	"github.com/mindchain/mind-chain/mc-challenger/game/fault/trace/cannon"
	"github.com/mindchain/mind-chain/mc-challenger/game/fault/trace/utils"
	"github.com/mindchain/mind-chain/mc-challenger/game/fault/trace/vm"
	"github.com/mindchain/mind-chain/mc-challenger/game/fault/types"
)

type prestateFetcher interface {
	getPrestate(ctx context.Context, logger log.Logger, prestateBaseUrl *url.URL, prestatePath string, dataDir string, stateConverter vm.StateConverter) (string, error)
}

func createTraceProvider(
	ctx context.Context,
	logger log.Logger,
	m vm.Metricer,
	cfg *config.Config,
	prestateSource prestateFetcher,
	traceType types.TraceType,
	localInputs utils.LocalGameInputs,
	dir string,
) (types.TraceProvider, error) {
	switch traceType {
	case types.TraceTypeCannon, types.TraceTypeSuperCannon:
		serverExecutor := vm.NewOpProgramServerExecutor(logger)
		stateConverter := cannon.NewStateConverter(cfg.Cannon)
		prestate, err := prestateSource.getPrestate(ctx, logger, cfg.CannonAbsolutePreStateBaseURL, cfg.CannonAbsolutePreState, dir, stateConverter)
		if err != nil {
			return nil, err
		}
		prestateProvider := vm.NewPrestateProvider(prestate, stateConverter)
		return cannon.NewTraceProvider(logger, m, cfg.Cannon, serverExecutor, prestateProvider, prestate, localInputs, dir, 42), nil
	case types.TraceTypeCannonKona:
		serverExecutor := vm.NewKonaExecutor()
		stateConverter := cannon.NewStateConverter(cfg.CannonKona)
		prestate, err := prestateSource.getPrestate(ctx, logger, cfg.CannonKonaAbsolutePreStateBaseURL, cfg.CannonKonaAbsolutePreState, dir, stateConverter)
		if err != nil {
			return nil, err
		}
		prestateProvider := vm.NewPrestateProvider(prestate, stateConverter)
		return cannon.NewTraceProvider(logger, m, cfg.CannonKona, serverExecutor, prestateProvider, prestate, localInputs, dir, 42), nil
	case types.TraceTypeAsterisc:
		serverExecutor := vm.NewOpProgramServerExecutor(logger)
		stateConverter := asterisc.NewStateConverter(cfg.Asterisc)
		prestate, err := prestateSource.getPrestate(ctx, logger, cfg.AsteriscAbsolutePreStateBaseURL, cfg.AsteriscAbsolutePreState, dir, stateConverter)
		if err != nil {
			return nil, err
		}
		prestateProvider := vm.NewPrestateProvider(prestate, stateConverter)
		return asterisc.NewTraceProvider(logger, m, cfg.Asterisc, serverExecutor, prestateProvider, prestate, localInputs, dir, 42), nil
	case types.TraceTypeAsteriscKona:
		serverExecutor := vm.NewKonaExecutor()
		stateConverter := asterisc.NewStateConverter(cfg.AsteriscKona)
		prestate, err := prestateSource.getPrestate(ctx, logger, cfg.AsteriscKonaAbsolutePreStateBaseURL, cfg.AsteriscKonaAbsolutePreState, dir, stateConverter)
		if err != nil {
			return nil, err
		}
		prestateProvider := vm.NewPrestateProvider(prestate, stateConverter)
		return asterisc.NewTraceProvider(logger, m, cfg.AsteriscKona, serverExecutor, prestateProvider, prestate, localInputs, dir, 42), nil
	case types.TraceTypeSuperAsteriscKona:
		serverExecutor := vm.NewKonaSuperExecutor()
		stateConverter := asterisc.NewStateConverter(cfg.AsteriscKona)
		prestate, err := prestateSource.getPrestate(ctx, logger, cfg.AsteriscKonaAbsolutePreStateBaseURL, cfg.AsteriscKonaAbsolutePreState, dir, stateConverter)
		if err != nil {
			return nil, err
		}
		prestateProvider := vm.NewPrestateProvider(prestate, stateConverter)
		return asterisc.NewTraceProvider(logger, m, cfg.AsteriscKona, serverExecutor, prestateProvider, prestate, localInputs, dir, 42), nil
	}
	return nil, errors.New("invalid trace type")
}
