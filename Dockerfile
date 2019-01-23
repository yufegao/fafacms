FROM golang:latest AS go-build

WORKDIR /go/src/github.com/hunterhug/fafa

COPY core /go/src/github.com/hunterhug/fafa/core
COPY vendor /go/src/github.com/hunterhug/fafa/vendor
COPY main.go /go/src/github.com/hunterhug/fafa/main.go

RUN go build -ldflags "-s -w" -o fafa main.go

FROM ubuntu:16.04 AS prod

WORKDIR /root/

COPY --from=go-build /go/src/github.com/hunterhug/fafa/fafa /bin/fafa
RUN chmod 777 /bin/fafa