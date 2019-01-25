FROM golang:latest AS go-build

WORKDIR /go/src/github.com/hunterhug/fafacms

COPY core /go/src/github.com/hunterhug/fafacms/core
COPY vendor /go/src/github.com/hunterhug/fafacms/vendor
COPY main.go /go/src/github.com/hunterhug/fafacms/main.go

RUN go build -ldflags "-s -w" -o fafacms main.go

FROM ubuntu:16.04 AS prod

WORKDIR /root/

COPY --from=go-build /go/src/github.com/hunterhug/fafacms/fafacms /bin/fafacms
RUN chmod 777 /bin/fafacms