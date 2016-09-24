package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type zkPoller struct {
	addr         string
	pollInterval time.Duration
	metrics      map[string]*prometheus.GaugeVec
	fetchStats   func(string, string) ([]byte, error)
}

func (zk zkPoller) start() {
	for {
		fmt.Printf("Polling zk %#v for metrics", zk)
		stats, err := zk.fetchAll()
		if err != nil {
			fmt.Printf("failed to fetch stats, err=%#v", err)
		}
		zk.updateMetrics(stats)
		<-time.After(zk.pollInterval)
	}
}

func (zk zkPoller) updateMetrics(stats map[string]string) {
	for stat, value := range stats {
		metric, ok := zk.metrics[stat]

		if !ok {
			fmt.Printf("couldn't find metric for stat %#v", stat)
			continue
		}

		if stat == zkOK {
			switch value {
			case "imok":
				metric.WithLabelValues(zk.addr).Set(1)
			default:
				metric.WithLabelValues(zk.addr).Set(0)
			}
			continue
		}

		if stat == zkServerState {
			state := getState(value)
			metric.WithLabelValues(zk.addr).Set(float64(state))
			continue
		}

		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("Failed to convert string value %#v to float\n", value)
			continue
		}

		metric.WithLabelValues(zk.addr).Set(f)
	}
}

func (zk zkPoller) fetchAll() (stats map[string]string, err error) {
	stats = make(map[string]string)

	mntrStats, err := zk.fetchStats(monitorCMD, zk.addr)
	if err != nil {
		return stats, err
	}

	okStats, err := zk.fetchStats(okCMD, zk.addr)
	if err != nil {
		return stats, err
	}

	r := io.MultiReader(bytes.NewReader(mntrStats), bytes.NewReader(okStats))

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		splits := strings.Split(scanner.Text(), "\t")
		if len(splits) != 2 {
			log.Printf("Expected a key value pair separated by a tab, got %#v\n", splits)
			continue
		}
		stats[splits[0]] = splits[1]
	}

	return stats, nil
}

func fetch4LWStats(cmd, addr string) ([]byte, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "dial failed")
	}

	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("Failed to close connection, err=%#v", err)
		}
	}()

	_, err = fmt.Fprintf(conn, fmt.Sprintf("%s\n", cmd))
	if err != nil {
		return nil, errors.Wrap(err, "cmd send failed")
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, conn)
	if err != nil {
		return nil, errors.Wrap(err, "fetch response failed")
	}

	return buf.Bytes(), nil
}
