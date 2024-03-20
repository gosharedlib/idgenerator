package redis

import (
	"context"
	"time"
)

// Pool redis连接池
type Pool interface {
	// Get 获取连接方法
	Get(ctx context.Context) (Conn, error)
}

// Conn 标准方法
type Conn interface {
	// SetNX set
	SetNX(key, value string, ttl time.Duration) (bool, error)
	// Expire expire
	Expire(key string, ttl time.Duration) (bool, error)
	// Del del
	Del(key string) (int64, error)
	// Close 关闭连接
	Close() error
}
