#!/bin/bash
cp -f dist/http_exporter /usr/local/bin/
mkdir -p /etc/http_exporter
cp config/config.yaml.tpl /etc/http_exporter/config.yaml
cp init/http_exporter.service /etc/systemd/system
systemctl enable http_exporter.service
