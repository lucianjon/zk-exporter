package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"time"

	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	port         int
	servers      string
	pollInterval time.Duration
)

const (
	monitorCMD = "mntr"
	okCMD      = "ruok"
)

func init() {
	flag.IntVar(&port, "port", 9120, "The port to serve the endpoint from.")
	flag.StringVar(&servers, "servers", "localhost:2181", "Comma separated list of zk servers in the format host:port")
	flag.DurationVar(&pollInterval, "pollinterval", 10*time.Second, "How often to poll zookeeper for metrics.")
	flag.Parse()
}

func main() {
	metrics := initMetrics()
	servers := strings.Split(servers, ",")
	if len(servers) == 0 {
		log.Fatalf("At least one zookeeper server is required.")
	}

	for _, server := range servers {
		zk := zkPoller{addr: server, pollInterval: pollInterval, metrics: metrics, fetchStats: fetch4LWStats}
		go zk.start()
	}

	http.Handle("/metrics", prometheus.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
}
