package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const (
	monitorCMD = "mntr"
	okCMD      = "ruok"
)

type zooKeeper struct {
	addr string
}

func newZooKeeper(addr string) *zooKeeper {
	return &zooKeeper{addr: addr}
}

func (zk *zooKeeper) fetchStats() (map[string]string, error) {
	stats, err := zk.fetchMntrStats()
	if err != nil {
		return stats, err
	}

	isOK, err := zk.fetchOKStat()
	if err != nil {
		return stats, err
	}

	stats[zkOK] = isOK
	return stats, nil
}

func (zk *zooKeeper) fetchMntrStats() (map[string]string, error) {
	stats := make(map[string]string)
	byts, err := zk.sendCommand(monitorCMD)
	if err != nil {
		return stats, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(byts))
	for scanner.Scan() {
		splits := strings.Split(scanner.Text(), "\t")
		if splits[0] == "zk_version" {
			continue
		}
		if len(splits) != 2 {
			log.Printf("zookeeper: expected a key value pair separated by a tab, got [%v]\n", splits)
			continue
		}
		stats[splits[0]] = splits[1]
	}
	return stats, nil
}

func (zk *zooKeeper) fetchOKStat() (string, error) {
	byts, err := zk.sendCommand(okCMD)
	return string(byts), err
}

func (zk *zooKeeper) sendCommand(cmd string) ([]byte, error) {
	conn, err := net.Dial("tcp", zk.addr)
	if err != nil {
		return nil, fmt.Errorf("zookeeper: dial failed, err=%#v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("Failed to close connection, err=%#v\n", err)
		}
	}()

	if _, err = fmt.Fprintf(conn, fmt.Sprintf("%s\n", cmd)); err != nil {
		return nil, fmt.Errorf("zookeeper: command send failed, err=%#v", err)
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, conn); err != nil {
		return nil, fmt.Errorf("zookeeper: fetch response failed, err=%#v", err)
	}

	return buf.Bytes(), nil
}
