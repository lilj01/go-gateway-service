FROM golang:1.17.5

MAINTAINER "lilj"
WORKDIR /usr/go-gateway
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
ADD . /usr/go-gateway
RUN go mod tidy
RUN go build main.go
EXPOSE 8880
ENTRYPOINT ["./main"]