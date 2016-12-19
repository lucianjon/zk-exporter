FROM alpine:3.5

ADD zk-exporter /zk-exporter

EXPOSE 9120

ENTRYPOINT ["/zk-exporter"]