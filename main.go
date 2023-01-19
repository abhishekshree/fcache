package main

import "github.com/abhishekshree/fcache/cache"

func main() {
	server := NewServer(ServerConfig{
		ListenAddr: ":8080",
		IsLeader:   true,
	}, cache.New())

	server.Start()
}
