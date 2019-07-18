package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/etcd-io/etcd/clientv3"
	"log"
	"time"
)

var (
	endpoint = flag.String("endpoint", "etcd:2379", "etcd endpoints")
)

func AccessEtcd(cli *clientv3.Client) {
	// https://etcd.io/docs/v3.3.12/demo/#access-etcd
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

func GetByPrefix(cli *clientv3.Client) {
	// https://etcd.io/docs/v3.3.12/demo/#get-by-prefix

	ctx := context.Background()
	prefix := "foo"

	for i := 0; i < 10; i++ {
		k := fmt.Sprintf("%s%d", prefix, i)
		v := fmt.Sprintf("bar%d", i)

		res, err := cli.Put(ctx, k, v)
		if err != nil {
			panic(err)
		}

		log.Printf("put(%s, %s): %v", k, v, res)
	}

	res, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	log.Printf("get(%s) with prefix: %v", prefix, res)
}

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

	AccessEtcd(cli)
	GetByPrefix(cli)
}
