FROM quay.io/prometheus/busybox-linux-amd64:glibc

COPY mysqld_exporter .

ENTRYPOINT [ "/mysqld_exporter" ]
