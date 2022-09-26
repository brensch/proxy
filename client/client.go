package client

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

var (
	timeout = 30 * time.Second
)

// Client contains all the proxy servers and satisfies the http.Client interface.
// Any request done by the client selects a random URI
type Client struct {
	uris []string

	c *http.Client

	tokenMu sync.Mutex
	token   *oauth2.Token
}

// InitClient looks for all the available proxies, and initialises the auth to use them.
func InitClient(projectID string) (*Client, error) {

	uris, err := AuditProxies(projectID)
	if err != nil {
		return nil, err
	}

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: timeout,
		}).Dial,
		TLSHandshakeTimeout: timeout,
	}

	netClient := &http.Client{
		Transport: netTransport,
		Timeout:   timeout,
	}

	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	client := &Client{
		uris:  uris,
		c:     netClient,
		token: token,
	}

	return client, nil

}

func GetAccessToken() (*oauth2.Token, error) {
	// NB it seemed like i would need to be specific about the audience,
	// but testing shows any audience will work for auth to our services.
	source, err := IDTokenTokenSource(context.Background(), "a.run.app")
	if err != nil {
		return nil, err
	}

	return source.Token()

}

func (c *Client) Do(req *http.Request, olog *zap.Logger) (*http.Response, error) {

	start := time.Now()

	proxy := c.uris[rand.Intn(len(c.uris))]
	log := olog.With(zap.String("proxy", proxy))
	log.Debug("started proxy request")

	proxyUrl, _ := url.Parse(proxy)

	// set the header to be the desired target that was the original url of the request
	req.Header.Set("X-Target", req.URL.String())

	// get auth token
	var newExpiry time.Time
	c.tokenMu.Lock()
	// check if we're valid while we're locked, revalidate if not
	if !c.token.Valid() {
		token, err := GetAccessToken()
		if err != nil {
			return nil, err
		}
		c.token = token

		// record time instead of log here to minimise time with held mutex
		newExpiry = token.Expiry
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))
	c.tokenMu.Unlock()

	if !newExpiry.IsZero() {
		log.Info("reauthed client", zap.Time("new_expiry", newExpiry))
	}

	req.Host = proxyUrl.Host

	// reset the url for the call to only be the host and scheme
	// (path and query are all captured in the header)
	url := &url.URL{
		Host:   proxyUrl.Host,
		Scheme: proxyUrl.Scheme,
	}
	req.URL = url

	res, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	log.Debug("finished proxy request", zap.Duration("execution_time_ms", time.Since(start)))

	return res, nil
}

func IDTokenTokenSource(ctx context.Context, audience string) (oauth2.TokenSource, error) {
	// First we try the idtoken package, which only works for service accounts
	ts, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		if err.Error() != `idtoken: credential must be service_account, found "authorized_user"` {
			return nil, err
		}
		// If that fails, we use our Application Default Credentials to fetch an id_token on the fly
		gts, err := google.DefaultTokenSource(ctx)
		if err != nil {
			return nil, err
		}
		ts = oauth2.ReuseTokenSource(nil, &idTokenSource{TokenSource: gts})
	}
	return ts, nil
}

// idTokenSource is an oauth2.TokenSource that wraps another
// It takes the id_token from TokenSource and passes that on as a bearer token
type idTokenSource struct {
	TokenSource oauth2.TokenSource
}

func (s *idTokenSource) Token() (*oauth2.Token, error) {
	token, err := s.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("token did not contain an id_token")
	}

	return &oauth2.Token{
		AccessToken: idToken,
		TokenType:   "Bearer",
		Expiry:      token.Expiry,
	}, nil
}
