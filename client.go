package proxy

import (
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
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

func (c *Client) Do(req *http.Request) (*http.Response, error) {

	proxyString := c.uris[rand.Intn(len(c.uris))]
	proxyUrl, _ := url.Parse(proxyString)

	req.Header.Set("User-Agent", RandomUserAgent())
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
