package main

import (
	"fmt"
	"testing"
	"time"
)

type testPoller struct {
	addr         string
	pollInterval time.Duration
	pollCalls    int
}

func mockStats(cmd, addr string) ([]byte, error) {
	if cmd == monitorCMD {
		return []byte(`zk_version	3.4.6-1569965, built on 02/20/2014 09:09 GMT
zk_avg_latency	0
zk_max_latency	382
zk_min_latency	0
zk_packets_received	98962360
zk_packets_sent	98969748
zk_num_alive_connections	90
zk_outstanding_requests	0
zk_server_state	follower
zk_znode_count	9002
zk_watch_count	851
zk_ephemerals_count	1908
zk_approximate_data_size	985524
zk_open_file_descriptor_count	119
zk_max_file_descriptor_count	1048576
`), nil
	}
	return []byte(`imok`), nil
}

func mockOK(cmd, addr string) ([]byte, error) {
	return []byte(`imok`), nil
}

func mockNotOK(cmd, addr string) ([]byte, error) {
	return []byte(`imdead`), nil
}

func TestUpdatesMntrMetrics(t *testing.T) {
	zk := zkPoller{addr: "test-zk:2181", pollInterval: 1 * time.Minute, fetchStats: mockStats}

	stats, err := zk.fetchAll()
	if err != nil {
		t.Errorf("failed to fetch stats, err=%#v", err)
	}

	fmt.Printf("stats: %#v", stats)
}
