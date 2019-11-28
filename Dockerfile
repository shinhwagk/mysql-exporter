FROM quay.io/repository/prometheus/busybox-linux-amd64

COPY mysqld_exporter .

ENTRYPOINT [ "/mysqld_exporter" ]
