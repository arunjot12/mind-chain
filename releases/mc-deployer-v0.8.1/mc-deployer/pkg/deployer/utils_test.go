package deployer

import (
	"testing"

	"github.com/ethereum-optimism/optimism/mc-deployer/pkg/deployer/flags"
	"github.com/stretchr/testify/require"
)

func TestEnsureDefaultCacheDir(t *testing.T) {
	cacheDir := flags.DefaultCacheDir()
	require.NotNil(t, cacheDir)
}
