package redisworker

import (
	"context"
	"github.com/gosharedlib/idgenerator/workid"
	"github.com/gosharedlib/idgenerator/workid/redisworker/redis"
	"github.com/pkg/errors"
	"log/slog"
	"strconv"
	"sync"
	"time"
)

const (
	workIDKey      = "workid:"        // 机器码key前缀
	maxWorkID      = 1024             // 机器码上限
	defaultModName = "default_mod"    // 默认模块名
	defaultTTL     = time.Second * 30 // 默认心跳时间
)

// redisWorker workId生成器配置
type redisWorker struct {
	AppName   string        // 服务名
	ModName   string        // 模块名
	Heartbeat time.Duration // 心跳时间
	pool      redis.Pool    // redis连接池
}

// NewRedisWorker 获取workID配置
func NewRedisWorker(appName string, pool redis.Pool) workid.Worker {
	return &redisWorker{
		AppName:   appName,
		ModName:   defaultModName,
		Heartbeat: defaultTTL,
		pool:      pool,
	}
}

func (c *redisWorker) Get(_ context.Context) workid.Conn {
	return &redisConn{
		appName:   c.AppName,
		modName:   c.ModName,
		timeout:   c.Heartbeat,
		pool:      c.pool,
		timerOnce: new(sync.Once),
	}
}

// SetAppName 设置模块名。如果一个服务里面多个业务需要各自的workId，则必须单独设置，否则不需要设置。默认值： default_mod
func (c *redisWorker) SetAppName(appName string) {
	if appName == "" {
		return
	}
	c.ModName = appName
}

// SetModName 设置模块名。如果一个服务里面多个业务需要各自的workId，则必须单独设置，否则不需要设置。默认值： default_mod
func (c *redisWorker) SetModName(modName string) {
	if modName == "" {
		return
	}
	c.ModName = modName
}

// SetHeartbeat 设置心跳时间，如果小于1s，用默认心跳时间
func (c *redisWorker) SetHeartbeat(heartbeat time.Duration) *redisWorker {
	if heartbeat < time.Second {
		return c
	}
	c.Heartbeat = heartbeat
	return c
}

// redisConn workID生成器配置
type redisConn struct {
	id        int           // id
	appName   string        // 服务名
	modName   string        // 模块名
	timeout   time.Duration // key过期时间
	pool      redis.Pool    // redis连接池
	timerOnce *sync.Once
}

// GetWorkID 获取workID
func (c *redisConn) GetWorkID(ctx context.Context) (workID int, err error) {
	workID, err = createWorkID(
		maxWorkID, func(n int) (bool, error) {
			key := workIDKey + c.appName + ":" + c.modName + ":" + strconv.Itoa(n)
			return c.add(ctx, key, "1")
		},
	)
	c.id = workID
	c.startTimer(ctx)
	return
}

func (c *redisConn) CleanWorkID(ctx context.Context) error {
	success, err := c.del(ctx)
	if err != nil {
		return err
	}

	if !success {
		err = errors.Errorf("del workid[%d] has fail", c.id)
	}

	return nil
}

// heartbeat 心跳
func (c *redisConn) heartbeat(ctx context.Context) {
	success, err := c.expire(ctx, c.getKey(), c.timeout*2+time.Second)
	if err != nil || !success {
		slog.WarnContext(ctx, "heartbeat", slog.Any("success", success), slog.Any("err", err))
	}
}

// add 新增workID
func (c *redisConn) add(ctx context.Context, key, value string) (bool, error) {
	var (
		conn    redis.Conn
		success bool
		err     error
	)
	conn, err = c.pool.Get(ctx)
	if err != nil {
		err = errors.WithStack(err)
		return success, err
	}
	success, err = conn.SetNX(key, value, c.timeout*2+time.Second)
	return success, errors.WithStack(err)
}

// del 删除workID
func (c *redisConn) del(ctx context.Context) (bool, error) {
	var (
		conn    redis.Conn
		success bool
		err     error
	)
	conn, err = c.pool.Get(ctx)
	if err != nil {
		err = errors.WithStack(err)
		return success, err
	}

	result, err := conn.Del(c.getKey())
	return result == 1, errors.WithStack(err)
}

func (c *redisConn) getKey() string {
	return workIDKey + c.appName + ":" + c.modName + ":" + strconv.Itoa(c.id)
}

// expire 设置workerID过期时间
func (c *redisConn) expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	var (
		conn    redis.Conn
		success bool
		err     error
	)
	conn, err = c.pool.Get(ctx)
	if err != nil {
		err = errors.WithStack(err)
		return success, err
	}
	success, err = conn.Expire(key, ttl)
	return success, errors.WithStack(err)
}

// startTimer 启动定时器
func (c *redisConn) startTimer(ctx context.Context) {
	c.timerOnce.Do(
		func() {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						slog.ErrorContext(ctx, "startTimer panic", slog.Any("panic", r))
					}
				}()

				for {
					time.Sleep(c.timeout)
					c.heartbeat(ctx)
				}
			}()
		},
	)
}

// createWorkID 创建workID
func createWorkID(times int, f func(n int) (bool, error)) (int, error) {
	var (
		workID  int
		success bool
		err     error
	)
	for i := 0; i < times; i++ {
		success, err = f(i)
		if success && err == nil {
			workID = i
			break
		}
	}
	if err != nil {
		err = errors.Wrap(err, "没有可用的workid")
		return workID, err
	}
	if !success {
		err = errors.New("没有可用的workid")
		return workID, err
	}

	return workID, err
}
