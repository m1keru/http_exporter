package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
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
	// Metrics
	metricResponseCode := promauto.NewGauge(prometheus.GaugeOpts{
		Name: fmt.Sprintf("%s_response_code", endpoint.MetricName),
		Help: "Last responseCode for endpoint",
	})
	metricResponseBodyAssert := promauto.NewGauge(prometheus.GaugeOpts{
		Name: fmt.Sprintf("%s_response_body_assert", endpoint.MetricName),
		Help: "Last response body regex assert for endpoint",
	})

	metricResponseTime := promauto.NewGauge(prometheus.GaugeOpts{
		Name: fmt.Sprintf("%s_response_time", endpoint.MetricName),
		Help: "Last response time for endpoint",
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
			start := time.Now()
			response, err := http.PostForm(endpoint.URL, data)
			if err != nil {
				log.SharedLogger.Errorf("crawler error:%+v\n", err)
				metricResponseCode.Set(float64(999))
				continue
			}
			elapsed := time.Since(start).Seconds()
			metricResponseTime.Set(float64(elapsed))
			processBody(*response, endpoint, metricResponseBodyAssert)
			metricResponseCode.Set(float64(response.StatusCode))
			response.Body.Close()
		case "POST-JSON":
			data, err := json.Marshal(endpoint.RequestData)
			if err != nil {
				log.SharedLogger.Error("unable to marshall requestData: %+v", endpoint.RequestData)
			}
			r, err := http.NewRequest("POST", endpoint.URL, bytes.NewBuffer(data))
			r.Header.Add("Content-Type", "application/json; charset=UTF-8")
			start := time.Now()
			response, err := client.Do(r)
			elapsed := time.Since(start).Seconds()
			metricResponseTime.Set(float64(elapsed))
			if err != nil {
				log.SharedLogger.Errorf("%+v\n", err)
				metricResponseCode.Set(float64(999))
				continue
			}
			processBody(*response, endpoint, metricResponseBodyAssert)
			metricResponseCode.Set(float64(response.StatusCode))
			response.Body.Close()
		default:
			urlData := "?"
			if endpoint.RequestData != nil {
				for key, val := range endpoint.RequestData {
					urlData += key + "=" + val + "&"
				}
			}
			start := time.Now()
			response, err := http.Get(endpoint.URL + urlData)
			elapsed := time.Since(start).Seconds()
			metricResponseTime.Set(float64(elapsed))
			if err != nil {
				log.SharedLogger.Errorf("%+v\n", err)
				metricResponseCode.Set(float64(999))
				continue
			}
			processBody(*response, endpoint, metricResponseBodyAssert)
			metricResponseCode.Set(float64(response.StatusCode))
			response.Body.Close()
		}
		log.SharedLogger.Debug(endpoint.MetricName)
		time.Sleep(time.Second * time.Duration(endpoint.ScrapeInverval))

		//Parse ResponseBody
		//Increase Metric
	}
}

func processBody(response http.Response, endpoint config.Endpoint, metricResponseBodyAssert prometheus.Gauge) {
	if endpoint.ResponseBodyRegex != "" {
		log.SharedLogger.Debug("ResponseBodyRegex:%s", endpoint.ResponseBodyRegex)
		responseRegex, err := regexp.Compile(endpoint.ResponseBodyRegex)
		if err != nil {
			log.SharedLogger.Errorf("crawler: unable to parse regex for %s", endpoint)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.SharedLogger.Errorf("crawler error: unable to read response body %+v\n")
		}
		if responseRegex.MatchString(string(body)) {
			metricResponseBodyAssert.Set(1)
		} else {
			metricResponseBodyAssert.Set(0)
		}
		log.SharedLogger.Debugf("%s", body)
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
