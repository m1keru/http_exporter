package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/m1keru/http_exporter/internal/config"
	"github.com/m1keru/http_exporter/internal/crawler"
	log "github.com/m1keru/http_exporter/internal/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// AppVersion -release version
const AppVersion = "0.0.17"

func main() {
	configpath := flag.String("config", "config.yaml", "path to config file")
	version := flag.Bool("version", false, "current version")
	flag.Parse()
	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}
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
