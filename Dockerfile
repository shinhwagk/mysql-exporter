FROM alpine:latest

COPY mysqld_exporter .

ENTRYPOINT [ "/mysqld_exporter" ]
