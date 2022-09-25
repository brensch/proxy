package client

import (
	"context"
	"io/ioutil"
	"testing"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/run/v1"
)

// TestAuditProxies works to validate the method used to audit proxies on initialisation of the client,
// and also to make sure you're deployed your proxies to the correct project and have given valid credentials.
func TestAuditProxies(t *testing.T) {
	ctx := context.Background()
	credentialBytes, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		t.Error("failed to read credentials file", err)
		return
	}

	credentials, err := google.CredentialsFromJSON(ctx, credentialBytes, run.CloudPlatformScope)
	if err != nil {
		t.Error("failed to get credentials", err)
		return
	}

	proxies, err := AuditProxies(credentials.ProjectID)
	if err != nil {
		t.Error("failed to audit proxies")
	}

	t.Log(len(proxies))
	if len(proxies) == 0 {
		t.Error("your proxies aren't deployed yet")
	}
}
