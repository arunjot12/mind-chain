package derive

import "github.com/mindchain/mind-chain/mc-service/testutils"

var _ L1Fetcher = (*testutils.MockL1Source)(nil)

var _ Metrics = (*testutils.TestDerivationMetrics)(nil)
