package proxy

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// HandleProxy doesn't need to do anything fancy, just some logging
// and the main HandlerFunc of the proxy package
func HandleProxy(w http.ResponseWriter, r *http.Request) {
	log := log.With(
		zap.String("dst", r.URL.String()),
		zap.String("method", r.Method),
	)
	log.Debug("received proxy request")
	start := time.Now()
	proxy.ServeHTTP(w, r)
	finish := time.Now()
	log.Debug("finished proxy request", zap.Duration("duration_ms", finish.Sub(start)))
}
