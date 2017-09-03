package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.Handle("/metrics", prometheus.Handler())

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		_ = <-sigs
		log.Println("main: received SIGINT or SIGTERM, shutting down")
		context, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := srv.Shutdown(context); err != nil {
			log.Printf("main: failed to shutdown endpoint with err=%#v\n", err)
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("main: failure while serving endpoint, err=%#v", err)
	}
}
