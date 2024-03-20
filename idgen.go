package idgenerator

import (
	guid "github.com/gofrs/uuid"
	"github.com/gosharedlib/idgenerator/md5"
	"github.com/gosharedlib/idgenerator/snowflake"
	"github.com/gosharedlib/idgenerator/uuid"
	"github.com/gosharedlib/idgenerator/workid"
	"github.com/gosharedlib/idgenerator/workid/redisworker"
	"github.com/gosharedlib/idgenerator/workid/redisworker/redis"
)

var (
	global = newDefaultGenerator()
)

// IDGenerator 代表 Key 生成器.
type IDGenerator interface {
	// NewRedisWorker 基于redis的workerID生成器
	NewRedisWorker(appName string, pool redis.Pool) workid.Worker
	// NewSnowflakeGenerator 雪花算法生成器
	NewSnowflakeGenerator(worker workid.Conn, epoch ...int64) snowflake.Generator
	// NewUUIDV1Generator UUID V1
	NewUUIDV1Generator() uuid.Generator
	// NewUUIDV2Generator UUID V2，由于安全缺陷，上游依赖已移除 V2 实现
	// NewUUIDV2Generator(domain byte) uuid.Generator
	// NewUUIDV3Generator UUID V3
	NewUUIDV3Generator(ns guid.UUID, name string) uuid.Generator
	// NewUUIDV4Generator UUID V4
	NewUUIDV4Generator() uuid.Generator
	// NewUUIDV5Generator UUID V5
	NewUUIDV5Generator(ns guid.UUID, name string) uuid.Generator
	// NewMD5Generator MD5
	NewMD5Generator(str string) md5.Generator
}

type idGenerator struct {
}

func newDefaultGenerator() IDGenerator {
	return &idGenerator{}
}

func NewRedisWorker(appName string, pool redis.Pool) workid.Worker {
	return global.NewRedisWorker(appName, pool)
}

func NewSnowflakeGenerator(worker workid.Conn, epoch ...int64) snowflake.Generator {
	return global.NewSnowflakeGenerator(worker, epoch...)
}

func NewUUIDV1Generator() uuid.Generator {
	return global.NewUUIDV1Generator()
}

// 由于安全缺陷，上游依赖已移除 V2 实现
// func NewUUIDV2Generator(domain byte) uuid.Generator {
// 	return global.NewUUIDV2Generator(domain)
// }

func NewUUIDV3Generator(ns guid.UUID, name string) uuid.Generator {
	return global.NewUUIDV3Generator(ns, name)
}

func NewUUIDV4Generator() uuid.Generator {
	return global.NewUUIDV4Generator()
}

func NewUUIDV5Generator(ns guid.UUID, name string) uuid.Generator {
	return global.NewUUIDV5Generator(ns, name)
}

func NewMD5Generator(str string) md5.Generator {
	return global.NewMD5Generator(str)
}

func (g *idGenerator) NewRedisWorker(appName string, pool redis.Pool) workid.Worker {
	return redisworker.NewRedisWorker(appName, pool)
}

func (g *idGenerator) NewSnowflakeGenerator(worker workid.Conn, epoch ...int64) snowflake.Generator {
	return snowflake.NewSnowflakeGenerator(worker, epoch...)
}

func (g *idGenerator) NewUUIDV1Generator() uuid.Generator {
	return uuid.NewV1Generator()
}

// 由于安全缺陷，上游依赖已移除 V2 实现
// func (g *idGenerator) NewUUIDV2Generator(domain byte) uuid.Generator {
// 	return uuid.NewV2Generator(domain)
// }

func (g *idGenerator) NewUUIDV3Generator(ns guid.UUID, name string) uuid.Generator {
	return uuid.NewV3Generator(ns, name)
}

func (g *idGenerator) NewUUIDV4Generator() uuid.Generator {
	return uuid.NewV4Generator()
}

func (g *idGenerator) NewUUIDV5Generator(ns guid.UUID, name string) uuid.Generator {
	return uuid.NewV5Generator(ns, name)
}

func (g *idGenerator) NewMD5Generator(str string) md5.Generator {
	return md5.NewMD5Generator(str)
}
