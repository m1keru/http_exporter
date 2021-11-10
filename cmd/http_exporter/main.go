package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"

	"github.com/m1keru/http_exporter/internal/config"
	"github.com/m1keru/http_exporter/internal/crawler"
	log "github.com/m1keru/http_exporter/internal/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	configpath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()
	var cfg config.Config
	if err := cfg.Setup(configpath); err != nil {
		panic(fmt.Errorf("config error: %+v", err))
	}
	log.InitLogger(cfg.Log)
	defer log.SharedLogger.Sync()

	var wg sync.WaitGroup
	crwlr := crawler.Crawler{Config: &cfg, WaitGroup: &wg}
	err := crwlr.Run()
	if err != nil {
		log.SharedLogger.Fatal(err.Error())
	}
	log.SharedLogger.Debug("waiting on WaitGroup")

	http.Handle("/metrics", promhttp.Handler())
	log.SharedLogger.Fatal(http.ListenAndServe(cfg.Daemon.Listen, nil))
	wg.Wait()
}
