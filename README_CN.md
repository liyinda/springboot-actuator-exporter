# springboot-actuator-exporter

通过springboot actuator metrics获取的json格式信息转换为Prometheus exporter metrics格式,

将prometheus获取的监控信息（如springboot程序的性能指标）使用grafana展示。


## 目录列表
* [依赖](#dependency)
* [下载](#download)
* [编译](#compile)
  * [build binary](#build-binary)
  * [build docker image](#build-docker-image)
* [运行](#run)
  * [run binary](#run-binary)
  * [run docker image](#run-docker-image)
  * [run parameter](#run-parameter)
* [环境变量](#environment-variables)
* [指标](#metrics)
  * [springboot_monitor_performance](#springboot_monitor_performance)
* [Grafana](#grafana)



## Dependency

```text
需要springboot开启actuator监控，通过HTTP输出JSON格式信息
E.g: curl http://localhost/management/metrics (springboot actuator metrics)
{
    "mem": 458972,
    "processors": 24,
    "uptime": 16774475011,
    "systemload.average": 0.87,
    "heap.used": 184541,
    "threads": 39,
    ...
} 

```

* [Springboot Actuator](https://docs.spring.io/spring-boot/docs/current/reference/html/production-ready-endpoints.html)
* [Prometheus](https://prometheus.io/)
* [Golang 1.11](https://golang.org/)


## Download

二进制程序下载地址如 [Releases](https://github.com/liyinda/springboot-actuator-exporter/releases) page.

## Compile

### build binary
```text
docker方式作者还未添加^-^
```

``` shell
go build springboot_actuatorMetrics_exporter.go
```
### build docker image
``` shell
make docker
```

## Docker Hub Image
``` shell
DOCKER The deployment method author will add as soon as possible 
docker pull null:latest
```
### run docker
```
docker run  -ti image  bin/springboot-actuator-exporter
```

### run parameter
```shell
程序运行参数
-springboot.scrape_uri string
    URI to stringboot metrics stub status page (default "http://localhost/management/metrics")
    获取stringboot actuator监控信息，JSON格式
    E.g: curl http://localhost/management/metrics (springboot actuator metrics)
    {
        "mem": 458972,
        "processors": 24,
        "uptime": 16774475011,
        "systemload.average": 0.87,
        "heap.used": 184541,
        "threads": 39,
        ...
    } 

-springboot.service string
    springboot服务名称
    springboot services (default "service")

-telemetry.address string
    启用端口号
    Address on which to expose metrics. (default ":9933")

-telemetry.endpoint string
    exporter endpoint位置
    Path under which to expose metrics. (default "/metrics")

```

## Environment variables

环境变量（无）

## Metrics

### springboot_monitor_performance
转换后的prometheus exporter metrics信息如下：

``` 
# TYPE springboot_actuator_exporter_build_info gauge
springboot_actuator_exporter_build_info{branch="",goversion="go1.11",revision="",version=""} 1
# HELP springboot_monitor_info springboot info
# TYPE springboot_monitor_info gauge
springboot_monitor_info{Processors="processors"} 16
# HELP springboot_monitor_performance springboot performance
# TYPE springboot_monitor_performance gauge
springboot_monitor_performance{hostname="$hostname",service="$service",sys="heap"} 250592
springboot_monitor_performance{hostname="$hostname",service="$service",sys="memory"} 1.11207e+06
springboot_monitor_performance{hostname="$hostname",service="$service",sys="systemload"} 0.06
springboot_monitor_performance{hostname="$hostname",service="$service",sys="threads"} 38
springboot_monitor_performance{hostname="$hostname",service="$service",sys="uptime"} 9.252913831e+09

```

### Grafana

![image](https://github.com/liyinda/springboot-actuator-exporter/blob/master/jpg/grafana.jpg)
