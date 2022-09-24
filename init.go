package proxy

import (
	"net/http"

	"github.com/elazarl/goproxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// required for cloud build
	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

var (
	proxy *goproxy.ProxyHttpServer
	log   *zap.Logger
)

func init() {

	// init logger
	logConfig := zap.NewProductionConfig()
	logConfig.Level.SetLevel(zap.DebugLevel)

	// this ensures google logs pick things up properly
	logConfig.EncoderConfig.MessageKey = "message"
	logConfig.EncoderConfig.LevelKey = "severity"
	logConfig.EncoderConfig.TimeKey = "time"
	logConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logConfig.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder

	var err error
	log, err = logConfig.Build()
	if err != nil {
		// this indicates a bug or some way that zap can fail i'm not aware of
		panic(err)
	}

	// set up proxy
	proxy = goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.Header.Set("User-Agent", RandomUserAgent())
			return r, nil
		})

}
