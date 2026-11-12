package sysgo

import (
	"net"
	"net/url"

	"github.com/mindchain/mind-chain/mc-service/testreq"
)

func ProxyAddr(require *testreq.Assertions, urlStr string) string {
	u, err := url.Parse(urlStr)
	require.NoError(err)
	return net.JoinHostPort(u.Hostname(), u.Port())
}
