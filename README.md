![example workflow](https://github.com/m1keru/http_exporter/actions/workflows/release.yaml/badge.svg)
## HTTP exporter 
####  It exposes `response codes`, `body regexes`, `request time` as metrics

### Supporter types of endpoints:
* GET
* POST-FORM
* POST-JSON

### build
```bash 
make
```

### Install
```bash
* make clean && make
* sudo make install
* vim /etc/http_exporter/config.yaml
* sudo systemctl start http_exporter 
```

### Install from release
```bash
* download release https://github.com/m1keru/http_exporter/releases
* mkdir /tmp/http_exporter 
* mv http_exporter-x.x.x.tar.gz /tmp/http_exporter
* cd /tmp/http_exporter
* tar -zxvf http_exporter-x.x.x.tar.gz
* sudo make install_release
```
### Example output
```
# HELP response_code Last responseCode for endpoint
# TYPE response_code gauge
http_exporter_response_code{url="https://api-android.example.com"} 200
http_exporter_response_code{url="https://api-ios.example.com"} 400
# HELP response_time Last response time for endpoint
# TYPE response_time gauge
http_exporter_response_time{url="https://api-android.example.com"} 0.46648195
http_exporter_response_time{url="https://api-ios.example.com"} 0.1831524
# TYPE response_body_assert gauge
http_exporter_response_body_assert{url="https://example.com/healthcheck"} 1
```
