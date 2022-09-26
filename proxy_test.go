package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

var (
	timeout = 30 * time.Second
)

func TestRecreationGovRequestExternal(t *testing.T) {
	endpoint := "https://www.recreation.gov/api/search/campsites?start=0&size=1000&fq=asset_id%3A232450&include_non_site_specific_campsites=true"

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set(TargetHeader, endpoint)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// this relies on external server existing so if errors occur, just skip
		t.Skip(err)
		return
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

	if res.StatusCode != http.StatusOK {
		t.Error("bad status code", res.StatusCode)
	}
}

func TestRecreationGovRequest(t *testing.T) {
	endpoint := "https://www.recreation.gov/api/search/campsites?start=0&size=1000&fq=asset_id%3A232450&include_non_site_specific_campsites=true"

	// can send to any arbitrary url here (would normally be a proxy.)
	// we intercept the output
	req := httptest.NewRequest(http.MethodGet, "http://mickeymouse.gov", nil)
	// set header to be the actual desired target
	req.Header.Set(TargetHeader, endpoint)

	w := httptest.NewRecorder()

	proxy := InitProxy(zap.NewExample())
	proxy.ServeHTTP(w, req)

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

	t.Log(res.Status)

	if len(resBytes) == 0 {
		t.Error("got empty response")
	}
	t.Log(len(resBytes))
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
	req.Header.Set(TargetHeader, ts.URL)
	w := httptest.NewRecorder()

	proxy := InitProxy(zap.NewExample())
	proxy.ServeHTTP(w, req)

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
