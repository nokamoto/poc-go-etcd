FROM golang:1.12.6-alpine3.9 AS build

RUN apk update && apk add git

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go list -m all

COPY main.go .

RUN go build .

FROM alpine:3.9

COPY --from=build /src/poc-go-etcd /bin/poc-go-etcd

ENTRYPOINT [ "/bin/poc-go-etcd" ]
