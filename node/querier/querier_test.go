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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisQuerier(t *testing.T) {
	// todos before testing:
	// $ docker pull redis:latest
	// $ docker run -d -p 6379:6379 redis:latest
	redisQuerier := &RedisQuerier{}
	redisQuerier.Init("localhost:6379")
	ctx := context.Background()
	err := redisQuerier.RedisClient.Set(ctx, "k1", "v1", 0).Err()
	require.Nil(t, err)

	value := redisQuerier.Do("k1")
	assert.Equal(t, "v1", value)
	value = redisQuerier.Do("k2")
	assert.Equal(t, "", value)

	err = redisQuerier.RedisClient.Del(ctx, "k1").Err()
	require.Nil(t, err)
	redisQuerier.Close()
}
