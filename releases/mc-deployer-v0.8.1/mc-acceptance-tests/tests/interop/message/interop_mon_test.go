package msg

import (
	"math/rand"
	"testing"
	"time"

	"github.com/mindchain/mind-chain/mc-acceptance-tests/tests/interop"
	"github.com/mindchain/mind-chain/mc-devstack/devtest"
	"github.com/mindchain/mind-chain/mc-devstack/presets"
	"github.com/mindchain/mind-chain/mc-interop-mon/monitor"
	"github.com/mindchain/mind-chain/mc-service/eth"
	opmetrics "github.com/mindchain/mind-chain/mc-service/metrics"
	"github.com/mindchain/mind-chain/mc-service/sources"
	"github.com/stretchr/testify/assert"
)

// TestInteropMon is testing that the mc-interop-mon metrics are correctly collected
func TestInteropMon(gt *testing.T) {
	t := devtest.ParallelT(gt)
	sys := presets.NewTwoL2SupernodeInterop(t, 0)

	clients := map[eth.ChainID]*sources.EthClient{
		sys.L2ELA.Escape().ChainID(): sys.L2ELA.Escape().EthClient().(*sources.EthClient),
		sys.L2ELB.Escape().ChainID(): sys.L2ELB.Escape().EthClient().(*sources.EthClient),
	}

	require := t.Require()

	// Start mc-interop-mon in the test context and attach to the devstack
	im, err := monitor.InteropMonitorServiceFromClients(t.Ctx(), "test", &monitor.CLIConfig{
		PollInterval: 50 * time.Millisecond,
		L2Rpcs:       []string{}, // unused here
		MetricsConfig: opmetrics.CLIConfig{
			Enabled: true,
		},
	}, clients, t.Logger())
	t.Require().NoError(err)
	require.NoError(im.Start(t.Ctx()))

	// two EOAs for triggering the init and exec interop txs
	alice := sys.FunderA.NewFundedEOA(eth.OneHundredthEther)
	bob := sys.FunderB.NewFundedEOA(eth.OneHundredthEther)
	eventLoggerAddress := alice.DeployEventLogger()

	// send initiating message on chain A
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	initMsg := alice.SendInitMessage(interop.RandomInitTrigger(rng, eventLoggerAddress, rng.Intn(3), rng.Intn(10)))

	// send executing message on chain B
	_ = bob.SendExecMessage(initMsg)

	// Ensure the metrics are generated
	require.EventuallyWithT(func(t *assert.CollectT) {
		checker := opmetrics.NewMetricChecker(t, im.Metrics.(opmetrics.RegistryMetricer).Registry())
		checker.FindByName("op_interop_mon_default_message_status")
	}, 2*time.Minute, 100*time.Millisecond)
	t.Log("mc-interop-mon metrics check successful")

	require.NoError(im.Stop(t.Ctx()))
}
