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

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port         int
	servers      string
	pollInterval time.Duration
	listenaddr   string
)

func init() {
	flag.IntVar(&port, "port", 9120, "The port to serve the endpoint from.")
	flag.StringVar(&listenaddr, "listenaddr", "", "Address to listen on")
	flag.StringVar(&servers, "servers", "", "Comma separated list of zk servers in the format host:port")
	flag.DurationVar(&pollInterval, "pollinterval", 10*time.Second, "How often to poll zookeeper for metrics.")
	flag.Parse()
}

func main() {
	ss := strings.Split(servers, ",")
	if servers == "" || len(ss) == 0 {
		log.Fatal("main: at least one zookeeper server is required")
	}

	metrics := initMetrics()

	for _, server := range ss {
		p := newPoller(pollInterval, metrics, newZooKeeper(server))
		go p.pollForMetrics()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", listenaddr, port),
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
		log.Printf("main: failure while serving endpoint, err=%#v\n", err)
	}
}
