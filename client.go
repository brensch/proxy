package proxy

import (
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

var (
	timeout = 30 * time.Second
)

// Client contains all the proxy servers and satisfies the http.Client interface.
// Any request done by the client selects a random URI
type Client struct {
	uris []string
}

func InitClient(projectID string) (*Client, error) {

	uris, err := AuditProxies(projectID)
	if err != nil {
		return nil, err
	}

	return &Client{
		uris: uris,
	}, nil

}

func (c *Client) Do(req *http.Request, log *zap.Logger) (*http.Response, error) {

	proxy := c.uris[rand.Intn(len(c.uris))]
	userAgent := RandomUserAgent()
	log.Debug("doing proxy request",
		zap.String("proxy", proxy),
		zap.String("user_agent", userAgent),
	)

	proxyUrl, _ := url.Parse(proxy)

	req.Header.Set("User-Agent", userAgent)
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: timeout,
		}).Dial,
		TLSHandshakeTimeout: timeout,
		Proxy:               http.ProxyURL(proxyUrl),
	}
	var netClient = &http.Client{
		Timeout:   timeout,
		Transport: netTransport,
	}

	return netClient.Do(req)
}
