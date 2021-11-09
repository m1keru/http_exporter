---
daemon:
log:
  level: debug
  path: "./app.log"
endpoints:
  - url: "https://example.com"
    metricName: example_1
    responseCode: 200
    requestType: GET
    scrapeInterval: 2
    timeout: 2
  - urlList:
      - "https://1.example.com/api.php"
      - "https://2.example.com/api.php"
    responseCode: 200
    metricName: example_2
    requestType: "POST-FORM"
    requestData:
      action: "config"
      adid: "monitoring"
      ads_version: "100.0.0"
      android_id: "monitoring"
      app_id: "monitoring"
      app_version: "100.0.0"
      device_brand: "monitoring"
      device_codename: "monitoring"
      device_manufacturer: "monitoring"
      device_model: "monitoring"
      devicetype: "phone"
      google_ad_id: "monitoring"
      instance_id: "monitoring"
      is_old_user: "0"
      limited_ad_tracking: "0"
      locale: "en,by"
      os_version: "100.0.0"
      platform: "android"
      resolution_app: "1x1"
      resolution_real: "1x1"
  - url: "https://example.com/api.php"
    responseCode: 200
    metricName: example_3
    requestType: "GET"
    requestData:
      action: config
      adid: monitoring-adid
      apiver: 100.0.0
      app_version: 100.0.0
      device_model: monitoring
      devicetype: phone
      lat: 0
      loc: en,by
      os_version: 100.0.0
      resolution: 1x1
      uid: 1111111111111111111111111111111111111111111111111111111111111111
      utc_offset: 0
  - url: "https://example.com/healthcheck"
    responseCode: 200
    metricName: example_4
    requestType: "GET"
    responseBodyRegex: ".*enabled.*"
  - url: "http://example.com/test/"
    responseCode: 200
    metricName: json_test
    requestType: "POST-JSON"
    requestData:
      id: "2"
      title: "ping"
      author: "http_exporter"

