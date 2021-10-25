package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

//Daemon - DaemonConfig
type Daemon struct {
	LogFile string `yaml:"logfile"`
	Debug   bool   `yaml:"debug"`
}

//Endpoint - Endpoint
type Endpoint struct {
	URL            string            `yaml:"url"`
	ResponseCode   string            `yaml:"responseCode"`
	MetricName     string            `yaml:"metricName"`
	RequestType    string            `yaml:"requestType"`
	RequestData    map[string]string `yaml:"requestData,flow"`
	ScrapeInverval int               `yaml:"scrapeInterval"`
	Timeout        int               `yaml:"timeout"`
}

//Log - Log
type Log struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}

//Config - Config
type Config struct {
	Daemon    Daemon     `yaml:"daemon"`
	Log       Log        `yaml:"log"`
	Endpoints []Endpoint `yaml:"endpoints,flow"`
}

//Setup - Setup
func (cfg *Config) Setup(filename *string) error {
	configFile, err := os.Open(*filename)
	if err != nil {
		return errors.New("unable to read config file")
	}
	defer configFile.Close()
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		return errors.New("Unable to Unmarshal Config")
	}
	return nil
}