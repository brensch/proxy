package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleProxyRequest(t *testing.T) {

	endpoint := fmt.Sprintf("%s/api/search/campsites?start=0&size=1000&fq=asset_id%%3A%s&include_non_site_specific_campsites=true", "https://www.recreation.gov", "232450")
	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	w := httptest.NewRecorder()
	HandleProxyRequest(w, req)

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
	HandleProxyRequest(w, req)

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
