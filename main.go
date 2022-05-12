package main

import (
	"net/http"
	"os"

	"github.com/athoune/fluentd-ambassador/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	prom := os.Getenv("PROMETHEUS")
	if prom == "" {
		prom = ":2112"
	}
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(prom, nil)
	err = s.Listen(listen)
	if err != nil {
		panic(err)
	}
}
