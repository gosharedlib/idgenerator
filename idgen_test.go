package idgenerator

import (
	"context"
	"reflect"
	"testing"

	goRedis "github.com/go-redis/redis"
	"github.com/gosharedlib/idgenerator/snowflake"
	"github.com/gosharedlib/idgenerator/workid"
	"github.com/gosharedlib/idgenerator/workid/redisworker"
	"github.com/gosharedlib/idgenerator/workid/redisworker/redis"
	"github.com/gosharedlib/idgenerator/workid/redisworker/redis/goredis"
)

func TestNewRedisWorker(t *testing.T) {
	type args struct {
		appName string
		pool    redis.Pool
	}
	goRedisOpt := &goRedis.Options{
		Network:  "tcp",
		Addr:     "192.168.0.128:6379",
		Password: "yourpassword",
		DB:       0,
	}
	pool := goredis.NewPool(goRedis.NewClient(goRedisOpt))
	tests := []struct {
		name string
		args args
		want workid.Worker
	}{
		{
			name: "test01",
			args: args{
				appName: "test01",
				pool:    pool,
			},
			want: redisworker.NewRedisWorker("test01", pool),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := NewRedisWorker(tt.args.appName, tt.args.pool)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewRedisWorker() = %v, want %v", got, tt.want)
				}
				t.Log(got.Get(context.TODO()).GetWorkID(context.TODO()))
			},
		)
	}
}

func TestNewSnowflakeGenerator(t *testing.T) {
	type args struct {
		worker workid.Conn
		epoch  []int64
	}
	goRedisOpt := &goRedis.Options{
		Network:  "tcp",
		Addr:     "192.168.0.128:6379",
		Password: "yourpassword",
		DB:       0,
	}
	pool := goredis.NewPool(goRedis.NewClient(goRedisOpt))
	worker := NewRedisWorker("test01", pool).Get(context.TODO())
	worker1 := NewRedisWorker("test02", pool).Get(context.TODO())
	tests := []struct {
		name string
		args args
		want snowflake.Generator
	}{
		{
			name: "test01",
			args: args{
				worker: worker,
				epoch:  []int64{},
			},
			want: snowflake.NewSnowflakeGenerator(worker1),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := NewSnowflakeGenerator(tt.args.worker, tt.args.epoch...)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewSnowflakeGenerator() = %v, want %v", got, tt.want)
				}
				t.Log(got.GenID())
				t.Log(got.GenIntID())
			},
		)
	}
}
