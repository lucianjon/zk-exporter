package main

import "github.com/prometheus/client_golang/prometheus"

type serverState float64

const (
	zkAvgLatency              = "zk_avg_latency"
	zkMinLatency              = "zk_min_latency"
	zkMaxLatency              = "zk_max_latency"
	zkPacketsReceived         = "zk_packets_received"
	zkPacketsSent             = "zk_packets_sent"
	zkNumAliveConnections     = "zk_num_alive_connections"
	zkOutstandingRequests     = "zk_outstanding_requests"
	zkZnodeCount              = "zk_znode_count"
	zkWatchCount              = "zk_watch_count"
	zkEphemeralsCount         = "zk_ephemerals_count"
	zkApproximateDataSize     = "zk_approximate_data_size"
	zkOpenFileDescriptorCount = "zk_open_file_descriptor_count"
	zkMaxFileDescriptorCount  = "zk_max_file_descriptor_count"
	zkFollowers               = "zk_followers"
	zkSyncedFollowers         = "zk_synced_followers"
	zkPendingSyncs            = "zk_pending_syncs"
	zkServerState             = "zk_server_state"

	zkOK = "zk_ok"

	// server states
	unknown    serverState = -1
	follower   serverState = 1
	leader     serverState = 2
	standalone serverState = 3
)

func getState(s string) serverState {
	switch s {
	case "follower":
		return follower
	case "leader":
		return leader
	case "standalone":
		return standalone
	default:
		return unknown
	}
}

func initMetrics() map[string]*prometheus.GaugeVec {

	metrics := make(map[string]*prometheus.GaugeVec)

	metrics[zkAvgLatency] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkAvgLatency,
		Help: "Average Latency for ZooKeeper network requests.",
	}, []string{"instance"})

	metrics[zkMinLatency] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkMinLatency,
		Help: "Minimum latency for Zookeeper network requests.",
	}, []string{"instance"})

	metrics[zkMaxLatency] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkMaxLatency,
		Help: "Maximum latency for ZooKeeper network requests",
	}, []string{"instance"})

	metrics[zkPacketsReceived] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkPacketsReceived,
		Help: "Number of network packets received by ZooKeeper instance.",
	}, []string{"instance"})

	metrics[zkPacketsSent] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkPacketsSent,
		Help: "Number of network packets sent by ZooKeeper instance.",
	}, []string{"instance"})

	metrics[zkNumAliveConnections] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkNumAliveConnections,
		Help: "Number of currently alive connections to the ZooKeeper instance.",
	}, []string{"instance"})

	metrics[zkOutstandingRequests] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkOutstandingRequests,
		Help: "Number of requests currently waiting in the queue.",
	}, []string{"instance"})

	metrics[zkZnodeCount] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkZnodeCount,
		Help: "Znode count",
	}, []string{"instance"})

	metrics[zkWatchCount] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkWatchCount,
		Help: "Watch count",
	}, []string{"instance"})

	metrics[zkEphemeralsCount] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkEphemeralsCount,
		Help: "Ephemerals Count",
	}, []string{"instance"})

	metrics[zkApproximateDataSize] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkApproximateDataSize,
		Help: "Approximate data size",
	}, []string{"instance"})

	metrics[zkOpenFileDescriptorCount] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkOpenFileDescriptorCount,
		Help: "Number of currently open file descriptors",
	}, []string{"instance"})

	metrics[zkMaxFileDescriptorCount] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkMaxFileDescriptorCount,
		Help: "Maximum number of open file descriptors",
	}, []string{"instance"})

	metrics[zkServerState] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkServerState,
		Help: "test help",
	}, []string{"instance"})

	metrics[zkFollowers] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkFollowers,
		Help: "Leader only: number of followers.",
	}, []string{"instance"})

	metrics[zkSyncedFollowers] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkSyncedFollowers,
		Help: "Leader only: number of followers currenty in sync",
	}, []string{"instance"})

	metrics[zkPendingSyncs] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkPendingSyncs,
		Help: "pending syncs",
	}, []string{"instance"})

	metrics[zkOK] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: zkOK,
		Help: "is zk ok",
	}, []string{"instance"})

	for _, metric := range metrics {
		prometheus.MustRegister(metric)
	}

	return metrics
}
