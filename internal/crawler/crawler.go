package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/m1keru/http_exporter/internal/config"
	log "github.com/m1keru/http_exporter/internal/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Crawler - Crawler
type Crawler struct {
	Config    *config.Config
	WaitGroup *sync.WaitGroup
}

func (crawler Crawler) endpointLoop(endpoint config.Endpoint) {
	log.SharedLogger.Debugf("started endpoint %s\n", endpoint.MetricName)
	counter := promauto.NewGauge(prometheus.GaugeOpts{
		Name: fmt.Sprintf("%s_response_code", endpoint.MetricName),
		Help: "Last responseCode for endpoint",
	})
	client := &http.Client{
		Timeout: time.Second * time.Duration(endpoint.Timeout),
	}

	if endpoint.ScrapeInverval == 0 {
		endpoint.ScrapeInverval = 15
	}

	if endpoint.Timeout == 0 {
		endpoint.Timeout = 3
	}

	for {
		switch endpoint.RequestType {
		case "POST-FORM":
			data := url.Values{}
			for key, val := range endpoint.RequestData {
				data.Set(key, val)
			}
			response, err := http.PostForm(endpoint.URL, data)
			if err != nil {
				log.SharedLogger.Errorf("%+v\n", err)
				counter.Set(float64(999))
				continue
			}
			counter.Set(float64(response.StatusCode))
			response.Body.Close()
		case "POST-JSON":
			data, err := json.Marshal(endpoint.RequestData)
			if err != nil {
				log.SharedLogger.Error("unable to marshall requestData: %+v", endpoint.RequestData)
			}
			r, err := http.NewRequest("POST", endpoint.URL, bytes.NewBuffer(data))
			r.Header.Add("Content-Type", "application/json; charset=UTF-8")
			response, err := client.Do(r)
			if err != nil {
				log.SharedLogger.Errorf("%+v\n", err)
				counter.Set(float64(999))
				continue
			}
			counter.Set(float64(response.StatusCode))
			response.Body.Close()
		default:
			urlData := "?"
			if endpoint.RequestData != nil {
				for key, val := range endpoint.RequestData {
					urlData += key + "=" + val + "&"
				}
			}
			response, err := http.Get(endpoint.URL + urlData)
			if err != nil {
				log.SharedLogger.Errorf("%+v\n", err)
				counter.Set(float64(999))
				continue
			}
			counter.Set(float64(response.StatusCode))
			response.Body.Close()
		}
		log.SharedLogger.Debug(endpoint.MetricName)
		time.Sleep(time.Second * time.Duration(endpoint.ScrapeInverval))

		//Parse ResponseBody
		//Increase Metric
	}
}

func (crawler Crawler) recordMetrics() {
	crawler.WaitGroup.Add(1)
	for _, endpoint := range crawler.Config.Endpoints {
		go crawler.endpointLoop(endpoint)
	}
	crawler.WaitGroup.Done()
}

// Run - Run
func (crawler Crawler) Run() error {
	crawler.recordMetrics()
	return nil
}
