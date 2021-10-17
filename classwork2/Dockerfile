FROM golang:latest AS build

ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on

WORKDIR /root/app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY *.go .

RUN CGO_ENABLED=0 go build -o /httpserver


FROM alpine:latest

COPY --from=build /httpserver /

CMD ["/httpserver"]

EXPOSE 8080
