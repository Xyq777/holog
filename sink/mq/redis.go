package mq

import (
	"github.com/go-redis/redis/v8"
)

type RedisStream struct {
	*redis.Client
}
