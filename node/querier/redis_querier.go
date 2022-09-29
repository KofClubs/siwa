/*
Copyright (c) 2022 Zhang Zhanpeng <zhangregister@outlook.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

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
