package main

import (
	"os"

	"github.com/athoune/fluentd-ambassador/server"
)

func main() {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "127.0.0.1:24224"
	}
	redis := os.Getenv("REDIS")
	if redis == "" {
		redis = "localhost:6379"
	}
	s, err := server.New(redis)
	if err != nil {
		panic(err)
	}
	err = s.Listen(listen)
	if err != nil {
		panic(err)
	}
}
