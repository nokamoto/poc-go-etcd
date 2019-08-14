package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
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

func Watch(cli *clientv3.Client) {
	// https://etcd.io/docs/v3.3.12/demo/#watch
	ctx := context.Background()
	size := 10

	watch := cli.Watch(ctx, "stock", clientv3.WithPrefix())

	go func() {
		for i := 0; i <= size; i++ {
			key := fmt.Sprintf("stock%d", i)
			value := fmt.Sprintf("value%d", i)

			res, err := cli.Put(ctx, key, value)
			if err != nil {
				panic(err)
			}

			log.Printf("put(%s, %s): %v", key, value, res)

			time.Sleep(3 * time.Second)
		}
	}()

	i := 0
	for observed := range watch {
		log.Printf("%d/%d: watch(stock) wth prefix: %v", i, size, observed)
		for j, ev := range observed.Events {
			log.Printf("%d/%d: events[%d]=%v", i, size, j, ev)
		}
		if i >= size {
			break
		}
		i = i + 1
	}

	// cli.Close() ?
}

func Lease(cli *clientv3.Client) {
	// https://etcd.io/docs/v3.3.12/demo/#lease
	ctx := context.Background()
	ttl := int64(300)
	lease, err := cli.Grant(ctx, ttl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("grant(%d): %v\n", ttl, lease)

	key := "sample"
	value := "value"

	res, err := cli.Put(ctx, key, value, clientv3.WithLease(lease.ID))
	if err != nil {
		panic(err)
	}

	fmt.Printf("put(%s, %s): %v\n", key, value, res)

	got, err := cli.Get(ctx, key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("get(%s): %v\n", key, got)

	ka, err := cli.KeepAliveOnce(ctx, lease.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("keepalive(%v): %v\n", lease.ID, ka)

	revoked, err := cli.Revoke(ctx, lease.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("revoke(%v): %v\n", lease.ID, revoked)

	got, err = cli.Get(ctx, key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("get(%s): %v\n", key, got)
}

func DistributedLocks(cli *clientv3.Client) {
	// https://etcd.io/docs/v3.3.12/demo/#distributed-locks
	s1, err := concurrency.NewSession(cli, concurrency.WithTTL(60))
	if err != nil {
		panic(err)
	}
	defer s1.Close()

	mu1 := concurrency.NewMutex(s1, "mutex1")
	fmt.Println("mu1 lock")
	err = mu1.Lock(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("mu1 acquired")

	go func() {
		fmt.Println("mu1 sleep 5 seconds...")
		time.Sleep(5 * time.Second)
		fmt.Println("mu1 unlock")
		mu1.Unlock(context.Background())
		fmt.Println("mu1 released")
	}()

	s2, err := concurrency.NewSession(cli, concurrency.WithTTL(60))
	if err != nil {
		panic(err)
	}
	defer s2.Close()

	mu2 := concurrency.NewMutex(s2, "mutex1")
	fmt.Println("mu2 lock")
	err = mu2.Lock(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("mu2 acquired")
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
	defer cli.Close()

	AccessEtcd(cli)
	GetByPrefix(cli)
	Delete(cli)
	TransactionalWrite(cli)
	Watch(cli)
	Lease(cli)
	DistributedLocks(cli)
}
