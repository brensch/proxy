package client

import (
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

var (
	timeout = 30 * time.Second
)

// Client contains all the proxy servers and satisfies the http.Client interface.
// Any request done by the client selects a random URI
type Client struct {
	uris []string
	opts []option.ClientOption
}

func InitClient(projectID string) (*Client, error) {

	uris, err := AuditProxies(projectID)
	if err != nil {
		return nil, err
	}

	client := &Client{
		uris: uris,
	}

	// detect if in cloud, use service account if so
	// idtoken doesn't allow user account types, needs to be a service account
	if os.Getenv("K_SERVICE") == "" {
		client.opts = []option.ClientOption{idtoken.WithCredentialsFile("credentials.json")}
	}

	return client, nil

}

func (c *Client) Do(req *http.Request, log *zap.Logger) (*http.Response, error) {

	proxy := c.uris[rand.Intn(len(c.uris))]
	log.Debug("doing proxy request",
		zap.String("proxy", proxy),
	)

	proxyUrl, _ := url.Parse(proxy)

	// not sure if this constant should live in this package or main
	req.Header.Set("X-Target", req.URL.String())
	req.Host = proxyUrl.Host
	req.URL.Host = proxyUrl.Host

	client, err := idtoken.NewClient(
		req.Context(),
		req.URL.String(),
		c.opts...,
	)
	if err != nil {
		return nil, err
	}

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: timeout,
		}).Dial,
		TLSHandshakeTimeout: timeout,
	}
	client.Timeout = timeout
	client.Transport = netTransport

	return client.Do(req)

}
