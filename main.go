package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/abhishekshree/fcache/cache"
	"github.com/abhishekshree/fcache/client"
)

func main() {
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader")
	)
	flag.Parse()

	config := ServerConfig{
		ListenAddr: *listenAddr,
		IsLeader:   true,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(5 * time.Second)
		if config.IsLeader {
			SendStuff()
		}
	}()

	server := NewServer(config, cache.New())
	server.Start()
}

func SendStuff() {
	for i := 0; i < 100; i++ {
		go func(i int) {
			client, err := client.New(":3000", client.Options{})
			if err != nil {
				log.Fatal(err)
			}

			var (
				key   = []byte(fmt.Sprintf("key_%d", i))
				value = []byte(fmt.Sprintf("val_%d", i))
			)

			err = client.Set(context.Background(), key, value, 0)
			if err != nil {
				log.Fatal(err)
			}

			fetchedValue, err := client.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(fetchedValue))

			client.Close()
		}(i)
	}
}
