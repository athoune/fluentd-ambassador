package server

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/athoune/fluent-server/server"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	fluentd   *server.Server
	redis     *redis.Client
	Stream    string
	marshaler func(interface{}) ([]byte, error)
	MaxLen    int64
}

func New(redisHost string) (*Server, error) {
	s := &Server{
		Stream:    "fluentd",
		marshaler: json.Marshal,
		MaxLen:    1024,
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
	values := make([]interface{}, len(record)*2+2)
	values[0] = "@tag"
	values[1] = tag
	i := 1
	var err error
	for k, v := range record {
		values[i*2] = k
		if reflect.ValueOf(v).Kind() == reflect.String {
			values[i*2+1] = v
		} else {
			values[i*2+1], err = s.marshaler(v)
			if err != nil {
				return err
			}
		}
		i++
	}
	cmd := s.redis.XAdd(ctx, &redis.XAddArgs{
		MaxLen: s.MaxLen,
		Stream: s.Stream,
		Values: values,
	})
	err = s.redis.Process(ctx, cmd)
	if err != nil {
		return err
	}
	fmt.Println(cmd.Result())

	return nil
}

func (s *Server) Listen(listen string) error {
	return s.fluentd.ListenAndServe(listen)
}
