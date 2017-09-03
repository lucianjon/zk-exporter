# ZooKeeper Exporter [![Build Status](https://travis-ci.org/lucianjon/zk-exporter.svg?branch=master)](https://travis-ci.org/lucianjon/zk-exporter)

Exposes prometheus metrics scraped from ZooKeeper via four letter word commands, see https://zookeeper.apache.org/doc/trunk/zookeeperAdmin.html#The+Four+Letter+Words

Currently exposes metrics from the "mntr" and "ruok" (requires ZooKeeper 3.4.0 or above)

Still a WIP

## Installing

```
go get -u github.com/lucianjon/zk-exporter
```

## Getting Started

* By default metrics are exposed on `localhost:9120/metrics`
