package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	port         int
	servers      string
	pollInterval time.Duration
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
		log.Fatal("main: at least one zookeeper server is required")
	}

	for _, server := range servers {
		p := newPoller(pollInterval, metrics, newZooKeeper(server))
		go p.pollForMetrics()
	}

	http.Handle("/metrics", prometheus.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
}
