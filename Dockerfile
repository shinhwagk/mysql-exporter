FROM quay.io/prometheus/busybox-linux-amd64

COPY mysqld_exporter .

ENTRYPOINT [ "/mysqld_exporter" ]
