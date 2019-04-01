package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

type SpringbootAcuator struct {
        //metrics
	Processors     uint64 `json:"processors"`
	Mem     uint64 `json:"mem"`
	Uptime     uint64 `json:"uptime"`
	Systemload     float64 `json:"systemload.average"`
	Heap     uint64 `json:"heap.used"`
	Threads     uint64 `json:"threads"`
}


type Exporter struct {
	URI string

	infoMetric                                                  *prometheus.Desc
	serverMetrics   map[string]*prometheus.Desc

}

func newServerMetric(metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(*metricsNamespace, "monitor", metricName),
		docString, labels, nil,
	)
}


func NewExporter(uri string) *Exporter {
	return &Exporter{
		URI:        uri,
                infoMetric: newServerMetric("info", "springboot info", []string{"Processors"}),
		serverMetrics: map[string]*prometheus.Desc{
			"performance": newServerMetric("performance", "springboot performance", []string{"sys", "hostname", "service"}),
		},
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.serverMetrics {
		ch <- m
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	body, err := fetchHTTP(e.URI, time.Duration(*springbootScrapeTimeout)*time.Second)()
	if err != nil {
		log.Println("fetchHTTP failed", err)
		return
	}
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println("ioutil.ReadAll failed", err)
		return
	}

	var springBoot SpringbootAcuator
	err = json.Unmarshal(data, &springBoot)
	if err != nil {
		log.Println("json.Unmarshal failed", err)
		return
	}

        hostname := getHostname()

	// info
        ch <- prometheus.MustNewConstMetric(e.infoMetric, prometheus.GaugeValue, float64(springBoot.Processors), "processors")

	// performance
	ch <- prometheus.MustNewConstMetric(e.serverMetrics["performance"], prometheus.GaugeValue, float64(springBoot.Uptime), "uptime", hostname, *springbootService)
	ch <- prometheus.MustNewConstMetric(e.serverMetrics["performance"], prometheus.GaugeValue, float64(springBoot.Mem), "memory", hostname, *springbootService)
	ch <- prometheus.MustNewConstMetric(e.serverMetrics["performance"], prometheus.GaugeValue, float64(springBoot.Systemload), "systemload", hostname, *springbootService)
	ch <- prometheus.MustNewConstMetric(e.serverMetrics["performance"], prometheus.GaugeValue, float64(springBoot.Heap), "heap", hostname, *springbootService)
	ch <- prometheus.MustNewConstMetric(e.serverMetrics["performance"], prometheus.GaugeValue, float64(springBoot.Threads), "threads", hostname, *springbootService)

}

func fetchHTTP(uri string, timeout time.Duration) func() (io.ReadCloser, error) {
	http.DefaultClient.Timeout = timeout
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: *insecure}

	return func() (io.ReadCloser, error) {
		resp, err := http.DefaultClient.Get(uri)
		if err != nil {
			return nil, err
		}
		if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
		}
                //io.Copy(os.Stdout, resp.Body)
		return resp.Body, nil

	}
}

var (
	showVersion        = flag.Bool("version", false, "Print version information.")
	listenAddress      = flag.String("telemetry.address", ":9933", "Address on which to expose metrics.")
	metricsEndpoint    = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
	metricsNamespace   = flag.String("metrics.namespace", "springboot", "Prometheus metrics namespace.")
	springbootScrapeURI     = flag.String("springboot.scrape_uri", "http://localhost/management/metrics", "URI to stringboot metrics stub status page")
	springbootService     = flag.String("springboot.service", "service", "springboot services")
	insecure           = flag.Bool("insecure", true, "Ignore server certificate if using https")
	springbootScrapeTimeout = flag.Int("springboot.scrape_timeout", 2, "The number of seconds to wait for an HTTP response from the stringboot.scrape_uri")
)

//Get the hostname
func getHostname() string {
    host, err := os.Hostname()
    if err != nil {
        fmt.Printf("%s", err)
    } 
    return host
}


func init() {
	prometheus.MustRegister(version.NewCollector("springboot_actuator_exporter"))
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("springboot_actuator_exporter"))
		os.Exit(0)
	}

	log.Printf("Starting springboot_actuator_exporter %s", version.Info())
	log.Printf("Build context %s", version.BuildContext())

	exporter := NewExporter(*springbootScrapeURI)
	prometheus.MustRegister(exporter)
	//prometheus.Unregister(prometheus.NewProcessCollector(os.Getpid(), ""))
	prometheus.Unregister(prometheus.NewGoCollector())

	http.Handle(*metricsEndpoint, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Nginx Exporter</title></head>
			<body>
			<h1>Nginx Exporter</h1>
			<p><a href="` + *metricsEndpoint + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Printf("Starting Server at : %s", *listenAddress)
	log.Printf("Metrics endpoint: %s", *metricsEndpoint)
	log.Printf("Metrics namespace: %s", *metricsNamespace)
	log.Printf("Scraping information from : %s", *springbootScrapeURI)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
