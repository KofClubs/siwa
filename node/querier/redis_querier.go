package querier

import (
	"context"
	"reflect"

	"github.com/MonteCarloClub/log"
	"github.com/MonteCarloClub/utils"
	"github.com/go-redis/redis/v8"
)

type RedisQuerier struct {
	RedisClient *redis.Client
}

func (redisQuerier *RedisQuerier) Init(args ...interface{}) {
	if redisQuerier == nil {
		log.Error("nil redis querier", "err", utils.NilPtrDerefErr)
		return
	}

	if len(args) < 1 {
		log.Error("fail to init redis querier", "argc", len(args))
		return
	}
	if len(args) > 1 {
		log.Warn("too many arguments to init redis querier", "argc", len(args))
	}
	var redisAddrString string
	var ok bool
	if redisAddrString, ok = args[0].(string); !ok {
		log.Error("wrong arg type to init redis querier", "arg", args[0],
			"arg type", reflect.TypeOf(args[0]).String())
		return
	}

	redisQuerier.RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddrString,
	})
}

func (redisQuerier *RedisQuerier) Do(expression string) string {
	if redisQuerier == nil || redisQuerier.RedisClient == nil {
		log.Error("nil redis querier or client", "err", utils.NilPtrDerefErr)
		return ""
	}

	value, err := redisQuerier.RedisClient.Get(context.Background(), expression).Result()
	if err != nil {
		log.Warn("fail to get value from redis", "key", expression, "err", err)
		return ""
	}
	return value
}

func (redisQuerier *RedisQuerier) Close() {
	if redisQuerier == nil {
		log.Error("nil redis querier", "err", utils.NilPtrDerefErr)
		return
	}
	redisQuerier.RedisClient = nil
}
