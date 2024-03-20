package goredis

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	redisWorker "github.com/gosharedlib/idgenerator/workid/redisworker/redis"
)

// pool 连接池信息
type pool struct {
	// redis连接
	delegate *redis.Client
}

// Get 获取redis连接
func (p *pool) Get(ctx context.Context) (redisWorker.Conn, error) {
	c := p.delegate
	if ctx != nil {
		c = c.WithContext(ctx)
	}
	return &conn{delegate: c}, nil
}

// NewPool 新建连接池
func NewPool(delegate *redis.Client) redisWorker.Pool {
	return &pool{delegate}
}

// conn 标准连接实现
type conn struct {
	delegate *redis.Client
}

func (c *conn) SetNX(key, value string, ttl time.Duration) (bool, error) {
	result, err := c.delegate.SetNX(key, value, ttl).Result()
	return result, noErrNil(err)
}

func (c *conn) Expire(key string, ttl time.Duration) (bool, error) {
	result, err := c.delegate.Expire(key, ttl).Result()
	return result, noErrNil(err)
}

func (c *conn) Del(key string) (int64, error) {
	result, err := c.delegate.Del(key).Result()
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
