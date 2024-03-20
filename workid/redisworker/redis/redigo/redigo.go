package redigo

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	redisWorker "github.com/gosharedlib/idgenerator/workid/redisworker/redis"
)

// pool 连接池信息
type pool struct {
	// redis连接
	delegate *redis.Pool
}

// Get 获取redis连接
func (p *pool) Get(ctx context.Context) (redisWorker.Conn, error) {
	if ctx != nil {
		c, err := p.delegate.GetContext(ctx)
		if err != nil {
			return nil, err
		}
		return &conn{c}, nil
	}
	return &conn{p.delegate.Get()}, nil
}

// NewPool 新建连接池
func NewPool(delegate *redis.Pool) redisWorker.Pool {
	return &pool{delegate: delegate}
}

// conn 标准连接实现
type conn struct {
	delegate redis.Conn
}

func (c *conn) SetNX(key, value string, ttl time.Duration) (bool, error) {
	result, err := redis.String(c.delegate.Do("SET", key, value, "EX", int64(ttl/time.Second), "NX"))
	return result == "OK", noErrNil(err)
}

func (c *conn) Expire(key string, ttl time.Duration) (bool, error) {
	result, err := redis.Int(c.delegate.Do("EXPIRE", key, int64(ttl/time.Second)))
	return result == 1, noErrNil(err)
}

func (c *conn) Del(key string) (int64, error) {
	result, err := redis.Int64(c.delegate.Do("DEL", key))
	return result, noErrNil(err)
}

// Close close
func (c *conn) Close() error {
	err := c.delegate.Close()
	return noErrNil(err)
}

// noErrNil redis nil判断
func noErrNil(err error) error {
	if err == redis.ErrNil {
		return nil
	}
	return err
}
