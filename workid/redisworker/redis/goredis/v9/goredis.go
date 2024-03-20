package goredis

import (
	"context"
	"time"

	redisWorker "github.com/gosharedlib/idgenerator/workid/redisworker/redis"
	"github.com/redis/go-redis/v9"
)

// pool 连接池信息
type pool struct {
	// redis连接
	delegate *redis.Client
}

// Get 获取redis连接
func (p *pool) Get(ctx context.Context) (redisWorker.Conn, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return &conn{p.delegate, ctx}, nil
}

// NewPool 新建连接池
func NewPool(delegate *redis.Client) redisWorker.Pool {
	return &pool{delegate: delegate}
}

// conn 标准连接实现
type conn struct {
	delegate *redis.Client
	ctx      context.Context
}

func (c *conn) SetNX(key, value string, ttl time.Duration) (bool, error) {
	result, err := c.delegate.SetNX(c.ctx, key, value, ttl).Result()
	return result, noErrNil(err)
}

func (c *conn) Expire(key string, ttl time.Duration) (bool, error) {
	result, err := c.delegate.Expire(c.ctx, key, ttl).Result()
	return result, noErrNil(err)
}

func (c *conn) Del(key string) (int64, error) {
	result, err := c.delegate.Del(c.ctx, key).Result()
	return result, noErrNil(err)
}

// Close close
func (c *conn) Close() error {
	// Not needed for this library
	return nil
}

// noErrNil redis nil判断
func noErrNil(err error) error {
	if err == redis.Nil {
		return nil
	}
	return err
}
