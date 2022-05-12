package server

import (
	"context"
	"fmt"
	"time"

	"github.com/athoune/fluent-server/server"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	fluentd *server.Server
	redis   *redis.Client
	stream  string
}

func New(redisHost string) (*Server, error) {
	s := &Server{
		stream: "fluentd",
	}
	var err error
	s.fluentd, err = server.New(s.Handle)
	if err != nil {
		return nil, err
	}

	s.redis = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})

	_, err = s.redis.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) Handle(tag string, time *time.Time, record map[string]interface{}) error {
	ctx := context.TODO()
	cmd := s.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: s.stream,
		Values: record,
	})
	err := s.redis.Process(ctx, cmd)
	if err != nil {
		return err
	}
	fmt.Println(cmd.Result())

	return nil
}

func (s *Server) Listen(listen string) error {
	return s.fluentd.ListenAndServe(listen)
}
