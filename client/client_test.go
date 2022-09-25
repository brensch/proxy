package client

import (
	"io/ioutil"
	"net/http"
	"testing"

	"go.uber.org/zap"
)

func TestUseClient(t *testing.T) {
	c, err := InitClient("763810810662")
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

	if c == nil {
		t.Error("what")
		return
	}

	res, err := c.Do(req, log)
	if err != nil {
		t.Error("failed to do request", err)
		return
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	t.Log(len(body))

	t.Log(res.Status)
}
