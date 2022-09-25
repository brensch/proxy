package main

import (
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// required for vendoring services on gcp
	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func init() {

}

func main() {
	// init logger
	logConfig := zap.NewProductionConfig()
	logConfig.Level.SetLevel(zap.DebugLevel)

	// this ensures google logs pick things up properly
	logConfig.EncoderConfig.MessageKey = "message"
	logConfig.EncoderConfig.LevelKey = "severity"
	logConfig.EncoderConfig.TimeKey = "time"
	logConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logConfig.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder

	log, err := logConfig.Build()
	if err != nil {
		panic(err)
	}

	proxy := InitProxy(log)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info("starting", zap.String("port", port))
	log.Error("got error serving",
		zap.Error(http.ListenAndServe(fmt.Sprintf(":%s", port), proxy)),
	)
}
