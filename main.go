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

func Delete(cli *clientv3.Client) {
	// https://etcd.io/docs/v3.3.12/demo/#delete
	ctx := context.Background()

	res, err := cli.Put(ctx, "foo", "bar")
	if err != nil {
		panic(err)
	}

	log.Printf("put(foo, bar): %v", res)

	deleted, err := cli.Delete(ctx, "foo")
	if err != nil {
		panic(err)
	}

	log.Printf("delete(foo): deleted=%v", deleted.Deleted)

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

	deleted, err = cli.Delete(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	log.Printf("delete(%s) with prefix: deleted=%v", prefix, deleted.Deleted)
}

func TransactionalWrite(cli *clientv3.Client) {
	//https://etcd.io/docs/v3.3.12/demo/#transactional-write
	ctx := context.TODO()

	res, err := cli.Put(context.Background(), "user1", "bad")
	if err != nil {
		panic(err)
	}

	log.Printf("put(user1, bad): %v", res)

	tx := cli.Txn(ctx)
	committed, err := tx.
		If(clientv3.Compare(clientv3.Value("user1"), "=", "bad")).
		Then(clientv3.OpDelete("user1")).
		Else(clientv3.OpPut("user1", "good")).
		Commit()
	if err != nil {
		panic(err)
	}

	log.Printf("compares(user1, =, bad) then delete(user1) else put(user1, good): %v", committed)

	got, err := cli.Get(context.Background(), "user1")
	if err != nil {
		panic(err)
	}

	log.Printf("get(user1): %v", got)
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
	Delete(cli)
	TransactionalWrite(cli)
}
