package main

import (
	"log"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type zkPoller struct {
	interval time.Duration
	metrics  map[string]*prometheus.GaugeVec
	zk       *zooKeeper
}

func newPoller(interval time.Duration, metrics map[string]*prometheus.GaugeVec, zk *zooKeeper) *zkPoller {
	return &zkPoller{
		interval: interval,
		metrics:  metrics,
		zk:       zk,
	}
}

func (p *zkPoller) pollForMetrics() {
	for {
		log.Printf("poller: polling zookeeper [%v] for metrics\n", p.zk.addr)
		m, err := p.zk.fetchStats()
		if err != nil {
			log.Printf("poller: failed to fetch stats, err=%v\n", err)
		}
		p.refreshMetrics(m)
		<-time.After(p.interval)
	}
}

func (p *zkPoller) refreshMetrics(updated map[string]string) {
	for name, value := range updated {
		metric, ok := p.metrics[name]

		if !ok {
			log.Printf("poller: couldn't find metric for stat=%v\n", name)
			continue
		}

		if name == zkOK {
			switch value {
			case "imok":
				metric.WithLabelValues(p.zk.addr).Set(1)
			default:
				metric.WithLabelValues(p.zk.addr).Set(0)
			}
			continue
		}

		if name == zkServerState {
			state := getState(value)
			metric.WithLabelValues(p.zk.addr).Set(float64(state))
			continue
		}

		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("poller: failed to convert string value to float, value=%v\n", value)
			continue
		}

		metric.WithLabelValues(p.zk.addr).Set(f)
	}
}
