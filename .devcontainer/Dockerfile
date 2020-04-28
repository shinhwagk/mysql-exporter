FROM golang:1

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct

RUN go get -u -v golang.org/x/tools/gopls@latest