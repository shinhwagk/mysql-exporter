FROM golang:1.13.4

# ENV GOBIN /go/bin

# RUN git clone --depth=1 https://github.com/shinhwagk/mysqld_exporter
# WORKDIR mysqld_exporter

# RUN go get -v

COPY mysql_exporter .

ENTRYPOINT [ "./mysqld_exporter" ]
