package client

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"go.uber.org/zap"
)

func TestUseClient(t *testing.T) {
	c, err := InitClient("proxy-362608")
	if err != nil {
		t.Error("failed to init client", err)
		return
	}

	log := zap.NewExample()

	endpoint := "https://www.recreation.gov/api/search/campsites?start=0&size=1000&fq=asset_id%3A232450&include_non_site_specific_campsites=true"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := c.Do(req, log)
	if err != nil {
		t.Error("failed to do request", err)
		return
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	t.Log(string(body))
	_ = body
	t.Log(res.Status)

	if res.StatusCode != http.StatusOK {
		t.Error("got bad status code", res.StatusCode)
	}
}

func TestTokenSource(t *testing.T) {

	source, err := IDTokenTokenSource(context.Background(), "http://cool.com")
	if err != nil {
		t.Error("bad token", err)
		return
	}

	token, err := source.Token()
	if err != nil {
		t.Error("failed to get token", err)
		return
	}

	t.Log(token.AccessToken)

}
