FROM golang:alpine
MAINTAINER JiangHongJie "jhj767658181@gmail.com"
WORKDIR $GOPATH/src/LogForce
ADD . ./
ENV GO111MODULE=on
# 使用代理拉取包
ENV GOPROXY="https://goproxy.io"
RUN go build .
EXPOSE 9992
ENTRYPOINT  ["./LogForce"]