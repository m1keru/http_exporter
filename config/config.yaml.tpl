---
daemon:
log:
  level: debug
  path: "./app.log"
endpoints:
  - url: "https://ya.ru"
    metricName: yandex
    responseCode: 200
    requestType: GET
    scrapeInterval: 2
    timeout: 2
  - url: "https://api-nonogram-android.easybrain.com/api.php"
    responseCode: 200
    metricName: api10_android
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
  - url: "https://api-aquapuzzle-ios.easybrain.com/api.php"
    responseCode: 200
    metricName: api10_ios
    requestType: "GET"
    requestData:
      action: config
      adid: monitoring-adid
      apiver: 100.0.0
      app: com.easybrain.monitoring
      app_version: 100.0.0
      device_model: monitoring
      devicetype: phone
      lat: 0
      loc: en,by
      os_version: 100.0.0
      resolution: 1x1
      uid: 1111111111111111111111111111111111111111111111111111111111111111
      utc_offset: 0
