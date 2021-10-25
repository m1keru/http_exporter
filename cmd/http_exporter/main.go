package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/m1keru/http_exporter/internal/config"
	"github.com/m1keru/http_exporter/internal/crawler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger

//InitLogger - InitLogger
func InitLogger(level string, path string) {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	var logLevel zapcore.Level
	logLevel.Set(level)
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	logger := zap.New(core)
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./app.log")
	return zapcore.AddSync(file)
}

func main() {
	configpath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()
	var cfg config.Config
	if err := cfg.Setup(configpath); err != nil {
		log.Fatalf("config error: %+v", err)
	}
	InitLogger(cfg.Log.Level, cfg.Log.Path)
	defer sugarLogger.Sync()

	msgChan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	crwlr := crawler.Crawler{Config: &cfg, WaitGroup: &wg, MsgChannel: &msgChan}
	err := crwlr.Run()
	if err != nil {
		sugarLogger.Fatal(err.Error())
	}
	sugarLogger.Debug("waiting on WaitGroup")

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
	wg.Wait()
}
