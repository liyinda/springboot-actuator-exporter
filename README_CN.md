# springboot-actuator-exporter

English: Springboot Actuator Metrics and exports them via HTTP for Prometheus consumption

Convert JSON format to prometheus format and via HTTP standard output.

中文: 通过springboot actuator metrics获取的json格式信息转换为Prometheus exporter metrics格式

中文README.md (https://github.com/liyinda/springboot-actuator-exporter/README_CN.md)

## Table of Contents
* [Dependency](#dependency)
* [Download](#download)
* [Compile](#compile)
  * [build binary](#build-binary)
  * [build docker image](#build-docker-image)
* [Run](#run)
  * [run binary](#run-binary)
  * [run docker image](#run-docker-image)
* [Environment variables](#environment-variables)
* [Metrics](#metrics)
  * [Server main](#server-main)
  * [Server zones](#server-zones)
  * [Filter zones](#filter-zones)



## Dependency

* [lsof](http://www.linuxfromscratch.org/blfs/view/svn/general/lsof.html)
* [Prometheus](https://prometheus.io/)
* [Golang 1.9.4](https://golang.org/)


## Download

Binary can be downloaded from [Releases](https://github.com/liyinda/secuity_exporter/releases) page.

## Compile

### build binary

``` shell
go build security_exporter.go
```
### build docker image
``` shell
make docker
```

## Docker Hub Image
``` shell
DOCKER 部署方式作者会尽快补充 
docker pull 空:latest
```
### run docker
```
docker run  -ti 镜像地址 bin/security_exporter
```

## Environment variables

This image is configurable using different env variables

## Metrics

Documents about exposed Prometheus metrics.

``` 
# HELP fail_password_total Number of Fail Password in /var/log/secure.
# TYPE fail_password_total counter
fail_password_total{host="$hostname",zone="datacenter"} 3
# HELP file_change_total Number of Change in /etc.
# TYPE file_change_total counter
file_change_total{host="$hostname",zone="datacenter"} 21
# HELP reverse_shell_total Number of Reverse Shell.
# TYPE reverse_shell_total counter
reverse_shell_total{host="$hostname",zone="datacenter"} 0

```

### Grafana

![image](https://github.com/liyinda/security_exporter/blob/master/jpg/grafana.jpg)
