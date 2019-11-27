FROM alpine:latest

COPY mysql_exporter .

ENTRYPOINT [ "/mysqld_exporter" ]
