## HTTP exporter 
####  It exposes `response codes`, `body regexes`, `request time` as metrics

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
