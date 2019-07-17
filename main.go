package main

import (
	"context"
	"flag"
	"github.com/etcd-io/etcd/clientv3"
	"log"
	"time"
)

var (
	endpoint = flag.String("endpoint", "etcd:2379", "etcd endpoints")
)

func main() {
	flag.Parse()

	log.Println("hello")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*endpoint},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	res, err := cli.Put(ctx, "foo", "bar")
	if err != nil {
		panic(err)
	}

	log.Printf("put(foo, bar): %v", res)

	got, err := cli.Get(ctx, "foo")
	if err != nil {
		panic(err)
	}

	log.Printf("get(foo): %v", got)
}
