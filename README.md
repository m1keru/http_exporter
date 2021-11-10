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
```
* make clean && make
* sudo make install
* vim /etc/http_exporter/config.yaml
* sudo systemctl start http_exporter 
```

### Install from release
* download release https://github.com/m1keru/http_exporter/releases
* mkdir /tmp/http_exporter 
* mv http_exporter-x.x.x.tar.gz /tmp/http_exporter
* cd /tmp/http_exporter
* tar -zxvf http_exporter-x.x.x.tar.gz
* sudo make install_release

### Example output
```
# HELP sudoku_health_response_body_assert Last response body regex assert for endpoint
# TYPE sudoku_health_response_body_assert gauge
sudoku_health_response_body_assert 1
# HELP sudoku_health_response_code Last responseCode for endpoint
# TYPE sudoku_health_response_code gauge
sudoku_health_response_code 200
# HELP sudoku_health_response_time Last response time for endpoint
# TYPE sudoku_health_response_time gauge
sudoku_health_response_time 0.436875086
```
