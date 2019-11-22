FROM golang:1.13.4

ENV GOBIN /go/bin

RUN git clone --depth=1 https://github.com/shinhwagk/mysqld_exporter
WORKDIR mysqld_exporter

RUN go get -v

RUN go build -o mysqld_exporter

ENTRYPOINT [ "./mysqld_exporter" ]
