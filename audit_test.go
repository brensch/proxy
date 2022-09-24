package proxy

import "testing"

func TestAuditProxies(t *testing.T) {
	proxies, err := AuditProxies("763810810662")
	if err != nil {
		t.Error("failed to audit proxies")
	}

	t.Log(len(proxies))
}
