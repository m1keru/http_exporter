package crawler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/m1keru/http_exporter/internal/config"
	log "github.com/m1keru/http_exporter/internal/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Crawler - Crawler
type Crawler struct {
	Config                   *config.Config
	WaitGroup                *sync.WaitGroup
	metricResponseCode       *prometheus.GaugeVec
	metricResponseBodyAssert *prometheus.GaugeVec
	metricResponseTime       *prometheus.GaugeVec
}

func (crawler Crawler) endpointLoop(endpoint config.Endpoint) {
	log.SharedLogger.Infof("started endpoint name:%s url:%s", endpoint.MetricName, endpoint.URL)

	if endpoint.URL != "" && endpoint.URLList != nil {
		log.SharedLogger.Fatalf("crawler: error in config - can not use `url` and `urlList` together")
	}
	if endpoint.URL == "" && endpoint.URLList != nil {
		re, _ := regexp.Compile(`[^\w]`)
		for _, subURL := range endpoint.URLList {
			subEndpoint := endpoint
			subEndpoint.URL = subURL
			id := re.ReplaceAllString(subURL, "")
			subEndpoint.MetricName = subEndpoint.MetricName + "_" + id
			subEndpoint.URLList = nil
			log.SharedLogger.Infof("%+v", subEndpoint)
			go crawler.endpointLoop(subEndpoint)
		}
		return
	}

	client := &http.Client{
		Timeout: time.Second * time.Duration(endpoint.Timeout),
	}

	if endpoint.ScrapeInverval == 0 {
		endpoint.ScrapeInverval = 15
	}

	if endpoint.Timeout == 0 {
		endpoint.Timeout = 3
	}

	if endpoint.Severity == "" {
		endpoint.Severity = "warning"
	}

	for {
		log.SharedLogger.Debugf("Endpoint: %+v", endpoint)
		switch endpoint.RequestType {
		case "POST-FORM":
			data := url.Values{}
			for key, val := range endpoint.RequestData {
				data.Set(key, val)
			}
			r, err := http.NewRequest("POST", endpoint.URL, strings.NewReader(data.Encode()))
			if err != nil {
				log.SharedLogger.Errorf("%+v\n", err)
			}
			if endpoint.BasicAuthUserName != "" && endpoint.BasicAuthPassword != "" {
				r.SetBasicAuth(endpoint.BasicAuthUserName, endpoint.BasicAuthPassword)
			}
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
			start := time.Now()
			response, err := client.Do(r)
			if err != nil {
				log.SharedLogger.Errorf("crawler error: endpoint: %+v\n %+v\n", endpoint, err)
				crawler.metricResponseCode.WithLabelValues(endpoint.URL, endpoint.Severity).Set(999)
				time.Sleep(time.Second * time.Duration(endpoint.ScrapeInverval))
				continue
			}
			elapsed := time.Since(start).Seconds()
			crawler.metricResponseTime.WithLabelValues(endpoint.URL, endpoint.Severity).Set(float64(elapsed))
			processBody(*response, endpoint, crawler.metricResponseBodyAssert)
			crawler.metricResponseCode.WithLabelValues(endpoint.URL, endpoint.Severity).Set(float64(response.StatusCode))
			response.Body.Close()
		case "POST-JSON":
			data, err := json.Marshal(endpoint.RequestData)
			if err != nil {
				log.SharedLogger.Error("unable to marshall requestData: %+v", endpoint.RequestData)
			}
			r, err := http.NewRequest("POST", endpoint.URL, bytes.NewBuffer(data))
			if err != nil {
				log.SharedLogger.Error("unable to do request err: %+v requestData: %+v", err, endpoint.RequestData)
			}
			if endpoint.BasicAuthUserName != "" && endpoint.BasicAuthPassword != "" {
				r.SetBasicAuth(endpoint.BasicAuthUserName, endpoint.BasicAuthPassword)
			}
			r.Header.Add("Content-Type", "application/json; charset=UTF-8")
			start := time.Now()
			response, err := client.Do(r)
			if err != nil {
				log.SharedLogger.Errorf("unable to do request to %s error:%+v\n", endpoint.URL, err)
			}
			elapsed := time.Since(start).Seconds()
			crawler.metricResponseTime.WithLabelValues(endpoint.URL, endpoint.Severity).Set(float64(elapsed))
			if err != nil {
				log.SharedLogger.Errorf("%+v\n", err)
				crawler.metricResponseCode.WithLabelValues(endpoint.URL, endpoint.Severity).Set(999)
				time.Sleep(time.Second * time.Duration(endpoint.ScrapeInverval))
				continue
			}
			processBody(*response, endpoint, crawler.metricResponseBodyAssert)
			crawler.metricResponseCode.WithLabelValues(endpoint.URL, endpoint.Severity).Set(float64(response.StatusCode))
			response.Body.Close()
		default:
			urlData := "?"
			if endpoint.RequestData != nil {
				for key, val := range endpoint.RequestData {
					urlData += key + "=" + val + "&"
				}
			}
			start := time.Now()
			req, err := http.NewRequest("GET", endpoint.URL+urlData, nil)
			if endpoint.BasicAuthUserName != "" && endpoint.BasicAuthPassword != "" {
				req.SetBasicAuth(endpoint.BasicAuthUserName, endpoint.BasicAuthPassword)
			}
			response, err := client.Do(req)
			if err != nil {
				log.SharedLogger.Errorf("unable to do request to %s error:%+v\n", endpoint.URL, err)
			}
			elapsed := time.Since(start).Seconds()
			crawler.metricResponseTime.WithLabelValues(endpoint.URL, endpoint.Severity).Set(float64(elapsed))
			if err != nil {
				log.SharedLogger.Errorf("%+v\n", err)
				crawler.metricResponseCode.WithLabelValues(endpoint.URL, endpoint.Severity).Set(999)
				time.Sleep(time.Second * time.Duration(endpoint.ScrapeInverval))
				continue
			}
			processBody(*response, endpoint, crawler.metricResponseBodyAssert)
			crawler.metricResponseCode.WithLabelValues(endpoint.URL, endpoint.Severity).Set(float64(response.StatusCode))
			response.Body.Close()
		}
		log.SharedLogger.Debug(endpoint.MetricName)
		time.Sleep(time.Second * time.Duration(endpoint.ScrapeInverval))
	}
}

func processBody(response http.Response, endpoint config.Endpoint, metricResponseBodyAssert *prometheus.GaugeVec) {
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
			metricResponseBodyAssert.WithLabelValues(endpoint.URL, endpoint.Severity).Set(1)
		} else {
			metricResponseBodyAssert.WithLabelValues(endpoint.URL, endpoint.Severity).Set(0)
		}
		log.SharedLogger.Debugf("endpoint: %s, response: %s", endpoint.URL, body)
	}
}

func (crawler Crawler) recordMetrics() {
	crawler.WaitGroup.Add(1)
	crawler.metricResponseCode = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_exporter_response_code",
		Help: "Last responseCode for endpoint",
	}, []string{"url", "severity"})
	crawler.metricResponseBodyAssert = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_exporter_response_body_assert",
		Help: "Last response body regex assert for endpoint",
	}, []string{"url", "severity"})

	crawler.metricResponseTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_exporter_response_time",
		Help: "Last response time for endpoint",
	}, []string{"url", "severity"})

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
