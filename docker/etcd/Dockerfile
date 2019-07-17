FROM golang:1.12.6-alpine3.9 AS build

RUN apk update && apk add git

RUN go get -v go.etcd.io/etcd
RUN go get -v go.etcd.io/etcd/etcdctl

RUN which etcd
RUN which etcdctl

FROM alpine:3.9

COPY --from=build /go/bin/etcd /bin/etcd
COPY --from=build /go/bin/etcdctl /bin/etcdctl

ENTRYPOINT [ "/bin/etcd" ]
