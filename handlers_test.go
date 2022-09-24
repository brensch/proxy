package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

func TestRecreationGovRequestExternal(t *testing.T) {
	endpoint := fmt.Sprintf("%s/api/search/campsites?start=0&size=1000&fq=asset_id%%3A%s&include_non_site_specific_campsites=true", "https://www.recreation.gov", "232450")
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	audience := "https://proxyrequest-fczsqdxnba-uw.a.run.app"
	c, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		t.Error("failed to get client", err)
	}
	proxyUrl, err := url.Parse(audience)
	// c := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	c.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

	res, err := c.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Error("got bad status code", res.StatusCode)
		return
	}

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	if len(resBytes) == 0 {
		t.Error("got empty response")
	}
}

func TestRecreationGovRequest(t *testing.T) {
	endpoint := fmt.Sprintf("%s/api/search/campsites?start=0&size=1000&fq=asset_id%%3A%s&include_non_site_specific_campsites=true", "https://www.recreation.gov", "232450")
	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	w := httptest.NewRecorder()
	HandleProxy(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Error("got bad status code", res.StatusCode)
		return
	}

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	if len(resBytes) == 0 {
		t.Error("got empty response")
	}
}

func TestProxyRequest(t *testing.T) {

	payload := []byte("strings are cool")

	// set up the testserver to receive the proxied request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		if strings.HasPrefix(r.Header["User-Agent"][0], "Go-http-client") {
			t.Error("got Go user agent")
		}
		w.WriteHeader(http.StatusOK)
		// write payload for testing at output
		w.Write(payload)
	}))
	defer ts.Close()

	// fire the proxy request to the testserver to check headers
	req := httptest.NewRequest(http.MethodGet, ts.URL, nil)
	w := httptest.NewRecorder()
	HandleProxy(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Error("got bad status code")
		return
	}

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(resBytes, payload) {
		t.Error("got different payload than was sent by test server")
	}

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
