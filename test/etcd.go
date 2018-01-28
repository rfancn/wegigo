package main

import (
	"fmt"
	"log"
	"time"
	"context"
	"github.com/coreos/etcd/clientv3"
)

const ETCD_TIMEOUT = 5 * time.Second

func main() {
	url := fmt.Sprintf("http://%s:%d", "127.0.0.1", 2379)
	fmt.Println(url)
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{url},
		DialTimeout: ETCD_TIMEOUT,
	})
	defer etcd.Close()

	if err != nil {
		log.Fatal("Error connect etcd:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ETCD_TIMEOUT)
	resp, err := etcd.Get(ctx, "/app", clientv3.WithPrefix())
	cancel()
	if err != nil {
		log.Println("Error read app uuids from etcd:", err)
	}

	fmt.Println(resp)

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}


}
