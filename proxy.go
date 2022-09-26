package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
)

const (
	TargetHeader = "X-Target"
)

func InitProxy(log *zap.Logger) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: Director(log),
	}
}

func Director(log *zap.Logger) func(req *http.Request) {
	return func(req *http.Request) {
		userAgent := RandomUserAgent()

		// the desired destination is read from this header
		target := req.Header.Get(TargetHeader)

		log.Debug("proxy request received",
			zap.String("user_agent", userAgent),
			zap.String("target", target),
		)

		targetURL, err := url.Parse(target)
		if err != nil {
			log.Error("failed to parse url in X-Target", zap.Error(err))
			return
		}

		// the url is updated to be the value retrieved from the header in the request
		req.URL = targetURL
		req.Host = targetURL.Host

		fmt.Println(req.Header.Get("Authorization"))

		// reset all headers for maximum incognito
		req.Header.Del("X-Forwarded-For")
		req.Header.Del(TargetHeader)

		// use a randomised useragent
		req.Header.Set("User-Agent", userAgent)
	}
}
