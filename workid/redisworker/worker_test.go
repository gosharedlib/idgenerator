package redisworker

import (
	"context"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	goRedis "github.com/go-redis/redis"
	rediGo "github.com/gomodule/redigo/redis"
	"github.com/gosharedlib/idgenerator/workid/redisworker/redis"
	"github.com/gosharedlib/idgenerator/workid/redisworker/redis/goredis"
	"github.com/gosharedlib/idgenerator/workid/redisworker/redis/redigo"
)

func TestConfig_GetWorkID(t *testing.T) {
	type fields struct {
		AppName   string
		ModName   string
		pool      redis.Pool
		Heartbeat time.Duration // 心跳时间
	}
	type args struct {
		ctx context.Context
	}
	goRedisOpt := &goRedis.Options{
		Network:  "tcp",
		Addr:     "192.168.0.128:6379",
		Password: "yourpassword",
		DB:       0,
	}
	pool := goredis.NewPool(goRedis.NewClient(goRedisOpt))
	rediGoPool := redigo.NewPool(
		&rediGo.Pool{
			MaxIdle:     1024,
			MaxActive:   60000,
			IdleTimeout: time.Minute,
			Dial: func() (rediGo.Conn, error) {
				conn, err := rediGo.Dial(
					"tcp", "192.168.0.128:6379",
					rediGo.DialConnectTimeout(time.Millisecond*200),
					rediGo.DialReadTimeout(time.Millisecond*500),
					rediGo.DialWriteTimeout(time.Millisecond*500),
					rediGo.DialPassword("yourpassword"),
					rediGo.DialDatabase(0),
				)
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
		},
	)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test_01",
			fields: fields{
				AppName:   "qw-scrm",
				ModName:   "cus",
				pool:      pool,
				Heartbeat: time.Second * 10,
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "test_02",
			fields: fields{
				AppName:   "qw-scrm",
				ModName:   "company",
				pool:      rediGoPool,
				Heartbeat: time.Second * 10,
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &redisWorker{
					AppName:   tt.fields.AppName,
					ModName:   tt.fields.ModName,
					pool:      tt.fields.pool,
					Heartbeat: tt.fields.Heartbeat,
				}
				gotWorkID, err := c.Get(context.TODO()).GetWorkID(tt.args.ctx)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetWorkID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotWorkID < 0 || gotWorkID >= maxWorkID {
					t.Errorf("GetWorkID() gotWorkID = %v", gotWorkID)
				}
			},
		)
	}
}

func TestConfig_heartbeat(t *testing.T) {
	type fields struct {
		ID        int64 // ID
		AppName   string
		ModName   string
		pool      redis.Pool
		Heartbeat time.Duration // 心跳时间
		timerOnce *sync.Once
	}
	type args struct {
		ctx context.Context
	}
	goRedisOpt := &goRedis.Options{
		Network:  "tcp",
		Addr:     "192.168.0.128:6379",
		Password: "yourpassword",
		DB:       0,
	}
	pool := goredis.NewPool(goRedis.NewClient(goRedisOpt))
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test_01",
			fields: fields{
				AppName:   "qw-scrm",
				ModName:   "cus",
				pool:      pool,
				Heartbeat: time.Second * 1,
				timerOnce: new(sync.Once),
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &redisWorker{
					AppName:   tt.fields.AppName,
					ModName:   tt.fields.ModName,
					pool:      tt.fields.pool,
					Heartbeat: tt.fields.Heartbeat,
				}
				gotWorkID, err := c.Get(context.TODO()).GetWorkID(tt.args.ctx)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetWorkID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotWorkID < 0 || gotWorkID >= maxWorkID {
					t.Errorf("GetWorkID() gotWorkID = %v", gotWorkID)
				}
				time.Sleep(time.Second * 5)
			},
		)
	}
}

func TestConfig_SetModName(t *testing.T) {
	type fields struct {
		ID        int64 // ID
		AppName   string
		ModName   string
		pool      redis.Pool
		Heartbeat time.Duration // 心跳时间
	}
	type args struct {
		modName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "test_01",
			fields: fields{
				ModName: defaultModName,
			},
			args: args{
				modName: "",
			},
			want: defaultModName,
		},
		{
			name: "test_02",
			fields: fields{
				ModName: defaultModName,
			},
			args: args{
				modName: "cus",
			},
			want: "cus",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &redisWorker{
					AppName: tt.fields.AppName,
					ModName: tt.fields.ModName,
					pool:    tt.fields.pool,
				}
				c.SetModName(tt.args.modName)
				if !reflect.DeepEqual(c.ModName, tt.want) {
					t.Errorf("SetModName() = %v, want %v", c, tt.want)
				}
			},
		)
	}
}

func TestConfig_add(t *testing.T) {
	type fields struct {
		ID        int64 // ID
		AppName   string
		ModName   string
		pool      redis.Pool
		Heartbeat time.Duration // 心跳时间
		timerOnce *sync.Once
	}
	type args struct {
		ctx   context.Context
		key   string
		value string
	}
	goRedisOpt := &goRedis.Options{
		Network:  "tcp",
		Addr:     "192.168.0.128:6379",
		Password: "yourpassword",
		DB:       0,
	}
	pool := goredis.NewPool(goRedis.NewClient(goRedisOpt))
	rediGoPool := redigo.NewPool(
		&rediGo.Pool{
			MaxIdle:     1024,
			MaxActive:   60000,
			IdleTimeout: time.Minute,
			Dial: func() (rediGo.Conn, error) {
				conn, err := rediGo.Dial(
					"tcp", "192.168.0.128:6379",
					rediGo.DialConnectTimeout(time.Millisecond*200),
					rediGo.DialReadTimeout(time.Millisecond*500),
					rediGo.DialWriteTimeout(time.Millisecond*500),
					rediGo.DialPassword("yourpassword"),
					rediGo.DialDatabase(0),
				)
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
		},
	)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test_01",
			fields: fields{
				AppName:   "qw-scrm",
				ModName:   "cus",
				pool:      pool,
				Heartbeat: time.Second * 10,
				timerOnce: new(sync.Once),
			},
			args: args{
				ctx:   context.TODO(),
				key:   "qw-scrm:cus",
				value: "0",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test_02",
			fields: fields{
				AppName:   "qw-scrm",
				ModName:   "cus",
				pool:      rediGoPool,
				Heartbeat: time.Second * 10,
				timerOnce: new(sync.Once),
			},
			args: args{
				ctx:   context.TODO(),
				key:   "qw-scrm:cus",
				value: "1",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &redisConn{
					appName:   tt.fields.AppName,
					modName:   tt.fields.ModName,
					pool:      tt.fields.pool,
					timeout:   tt.fields.Heartbeat,
					timerOnce: tt.fields.timerOnce,
				}
				got, err := c.add(tt.args.ctx, tt.args.key, tt.args.value)
				if (err != nil) != tt.wantErr {
					t.Errorf("sAdd() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("sAdd() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestConfig_expire(t *testing.T) {
	type fields struct {
		ID        int64 // ID
		AppName   string
		ModName   string
		pool      redis.Pool
		Heartbeat time.Duration // 心跳时间
		timerOnce *sync.Once
	}
	type args struct {
		ctx context.Context
		key string
		ttl time.Duration
	}
	goRedisOpt := &goRedis.Options{
		Network:  "tcp",
		Addr:     "192.168.0.128:6379",
		Password: "yourpassword",
		DB:       0,
	}
	pool := goredis.NewPool(goRedis.NewClient(goRedisOpt))
	rediGoPool := redigo.NewPool(
		&rediGo.Pool{
			MaxIdle:     1024,
			MaxActive:   60000,
			IdleTimeout: time.Minute,
			Dial: func() (rediGo.Conn, error) {
				conn, err := rediGo.Dial(
					"tcp", "192.168.0.128:6379",
					rediGo.DialConnectTimeout(time.Millisecond*200),
					rediGo.DialReadTimeout(time.Millisecond*500),
					rediGo.DialWriteTimeout(time.Millisecond*500),
					rediGo.DialPassword("yourpassword"),
					rediGo.DialDatabase(0),
				)
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
		},
	)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test_01",
			fields: fields{
				AppName:   "qw-scrm",
				ModName:   "cus",
				pool:      pool,
				Heartbeat: time.Second * 10,
				timerOnce: new(sync.Once),
			},
			args: args{
				ctx: context.TODO(),
				key: "qw-scrm:cus",
				ttl: time.Second * 10,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test_02",
			fields: fields{
				AppName:   "qw-scrm",
				ModName:   "cus",
				pool:      rediGoPool,
				Heartbeat: time.Second * 10,
				timerOnce: new(sync.Once),
			},
			args: args{
				ctx: context.TODO(),
				key: "qw-scrm:cus",
				ttl: time.Second * 10,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &redisConn{
					appName:   tt.fields.AppName,
					modName:   tt.fields.ModName,
					pool:      tt.fields.pool,
					timeout:   tt.fields.Heartbeat,
					timerOnce: tt.fields.timerOnce,
				}
				_, _ = c.add(tt.args.ctx, tt.args.key, "1")
				got, err := c.expire(tt.args.ctx, tt.args.key, tt.args.ttl)
				if (err != nil) != tt.wantErr {
					t.Errorf("sRem() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("sRem() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestNewWorker(t *testing.T) {
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
	rediGoPool := &rediGo.Pool{
		MaxIdle:     1024,
		MaxActive:   60000,
		IdleTimeout: time.Minute,
		Dial: func() (rediGo.Conn, error) {
			conn, err := rediGo.Dial(
				"tcp", "192.168.0.128:6379",
				rediGo.DialConnectTimeout(time.Millisecond*200),
				rediGo.DialReadTimeout(time.Millisecond*500),
				rediGo.DialWriteTimeout(time.Millisecond*500),
				rediGo.DialPassword("yourpassword"),
				rediGo.DialDatabase(0),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test_01",
			args: args{
				appName: "qw-scrm",
				pool:    goredis.NewPool(goRedis.NewClient(goRedisOpt)),
			},
			wantErr: false,
		},
		{
			name: "test_02",
			args: args{
				appName: "qw-scrm",
				pool:    redigo.NewPool(rediGoPool),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewRedisWorker(tt.args.appName, tt.args.pool); !reflect.DeepEqual(got == nil, tt.wantErr) {
					t.Errorf("NewConfig() = %v, want %v", got, tt.wantErr)
				}
			},
		)
	}
}

func TestWorker_SetHeartbeat(t *testing.T) {
	type fields struct {
		AppName   string
		ModName   string
		Heartbeat time.Duration
		pool      redis.Pool
		timerOnce *sync.Once
	}
	type args struct {
		heartbeat time.Duration
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantHeartbeat time.Duration
	}{
		{
			name: "test_01",
			fields: fields{
				AppName:   "qw-scrm",
				Heartbeat: 0,
			},
			args: args{
				heartbeat: time.Second,
			},
			wantHeartbeat: time.Second,
		},
		{
			name: "test_02",
			fields: fields{
				AppName:   "qw-scrm",
				Heartbeat: defaultTTL,
			},
			args: args{
				heartbeat: time.Millisecond,
			},
			wantHeartbeat: defaultTTL,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &redisWorker{
					AppName:   tt.fields.AppName,
					ModName:   tt.fields.ModName,
					Heartbeat: tt.fields.Heartbeat,
					pool:      tt.fields.pool,
				}
				if got := c.SetHeartbeat(tt.args.heartbeat); !reflect.DeepEqual(got.Heartbeat, tt.wantHeartbeat) {
					t.Errorf("SetHeartbeat() = %v, want %v", got, tt.wantHeartbeat)
				}
			},
		)
	}
}

func TestWorker_add(t *testing.T) {
	type fields struct {
		ID        int
		AppName   string
		ModName   string
		Heartbeat time.Duration
		pool      redis.Pool
		timerOnce *sync.Once
	}
	type args struct {
		ctx   context.Context
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test_01",
			fields: fields{
				ID:        0,
				AppName:   "qw-scrm",
				ModName:   defaultModName,
				Heartbeat: 0,
				pool:      redigo.NewPool(&rediGo.Pool{}),
				timerOnce: nil,
			},
			args: args{
				ctx:   context.TODO(),
				key:   "qw-scrm:cus:1",
				value: "1",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &redisConn{
					appName:   tt.fields.AppName,
					modName:   tt.fields.ModName,
					timeout:   tt.fields.Heartbeat,
					pool:      tt.fields.pool,
					timerOnce: tt.fields.timerOnce,
				}
				got, err := c.add(tt.args.ctx, tt.args.key, tt.args.value)
				if (err != nil) != tt.wantErr {
					t.Errorf("add() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("add() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestWorker_expire(t *testing.T) {
	type fields struct {
		ID        int
		AppName   string
		ModName   string
		Heartbeat time.Duration
		pool      redis.Pool
		timerOnce *sync.Once
	}
	type args struct {
		ctx context.Context
		key string
		ttl time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test_01",
			fields: fields{
				ID:        0,
				AppName:   "qw-scrm",
				ModName:   defaultModName,
				Heartbeat: 0,
				pool:      redigo.NewPool(&rediGo.Pool{}),
				timerOnce: nil,
			},
			args: args{
				ctx: context.TODO(),
				key: "qw-scrm:cus:1",
				ttl: 100,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &redisConn{
					id:        tt.fields.ID,
					appName:   tt.fields.AppName,
					modName:   tt.fields.ModName,
					timeout:   tt.fields.Heartbeat,
					pool:      tt.fields.pool,
					timerOnce: tt.fields.timerOnce,
				}
				got, err := c.expire(tt.args.ctx, tt.args.key, tt.args.ttl)
				if (err != nil) != tt.wantErr {
					t.Errorf("expire() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("expire() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestWorker_heart(t *testing.T) {
	type fields struct {
		ID        int
		AppName   string
		ModName   string
		Heartbeat time.Duration
		pool      redis.Pool
		timerOnce *sync.Once
	}
	type args struct {
		ctx context.Context
		key string
		ttl time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test_01",
			fields: fields{
				ID:        0,
				AppName:   "qw-scrm",
				ModName:   defaultModName,
				Heartbeat: 0,
				pool:      redigo.NewPool(&rediGo.Pool{}),
				timerOnce: nil,
			},
			args: args{
				ctx: context.TODO(),
				key: "qw-scrm:cus:1",
				ttl: 100,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &redisConn{
					id:        tt.fields.ID,
					appName:   tt.fields.AppName,
					modName:   tt.fields.ModName,
					timeout:   tt.fields.Heartbeat,
					pool:      tt.fields.pool,
					timerOnce: tt.fields.timerOnce,
				}
				c.heartbeat(context.TODO())
			},
		)
	}
}

func TestConfig_GetWorkIDErr(t *testing.T) {
	type fields struct {
		ID        int64 // ID
		AppName   string
		ModName   string
		pool      redis.Pool
		Heartbeat time.Duration // 心跳时间
		timerOnce *sync.Once
	}
	type args struct {
		ctx context.Context
	}
	goRedisOpt := &goRedis.Options{
		Network:  "tcp",
		Addr:     "192.168.0.128:6379",
		Password: "yourpassword",
		DB:       0,
	}
	pool := goredis.NewPool(goRedis.NewClient(goRedisOpt))
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test_01",
			fields: fields{
				AppName:   "qw-scrm-1024",
				ModName:   "cus",
				pool:      pool,
				Heartbeat: time.Second * 200,
				timerOnce: new(sync.Once),
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &redisConn{
					appName:   tt.fields.AppName,
					modName:   tt.fields.ModName,
					pool:      tt.fields.pool,
					timeout:   tt.fields.Heartbeat,
					timerOnce: tt.fields.timerOnce,
				}

				for i := 0; i < maxWorkID; i++ {
					key := workIDKey + c.appName + ":" + c.modName + ":" + strconv.Itoa(i)
					_, _ = c.add(tt.args.ctx, key, "1")
				}

				gotWorkID, err := c.GetWorkID(tt.args.ctx)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetWorkID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotWorkID < 0 || gotWorkID >= maxWorkID {
					t.Errorf("GetWorkID() gotWorkID = %v", gotWorkID)
				}
			},
		)
	}
}
