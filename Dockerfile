FROM golang:1.19.0 AS builder
RUN go env -w GO111MODULE=auto \
  && go env -w GOPROXY=https://goproxy.cn,direct 
WORKDIR /build
COPY ./ .
RUN cd /build && go build -tags netgo -ldflags="-w -s" -o clash-sub main.go 

FROM alpine
LABEL MAINTAINER=github.com/dezhiShen
WORKDIR /data
RUN apk add -U --repository http://mirrors.ustc.edu.cn/alpine/v3.13/main/ tzdata 
COPY --from=builder /build/clash-sub /usr/bin/clash-sub 
RUN chmod +x /usr/bin/clash-sub
VOLUME /data
ENTRYPOINT ["/usr/bin/clash-sub"]