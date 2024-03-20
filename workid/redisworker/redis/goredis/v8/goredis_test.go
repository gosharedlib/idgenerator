package goredis

import (
	"context"
	"reflect"
	"testing"
	"time"

	goRedis "github.com/go-redis/redis/v8"
)

func Test_conn_Close(t *testing.T) {
	type fields struct {
		delegate *goRedis.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "test_01",
			fields:  fields{delegate: &goRedis.Client{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &conn{
					delegate: tt.fields.delegate,
				}
				if err := c.Close(); (err != nil) != tt.wantErr {
					t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func Test_noErrNil(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test_01",
			args:    args{err: goRedis.Nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := noErrNil(tt.args.err); (err != nil) != tt.wantErr {
					t.Errorf("noErrNil() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func Test_pool_Get(t *testing.T) {
	type fields struct {
		delegate *goRedis.Client
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
	client := goRedis.NewClient(goRedisOpt)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test01",
			fields:  fields{delegate: client},
			args:    args{ctx: context.TODO()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				p := &pool{
					delegate: tt.fields.delegate,
				}
				_, err := p.Get(tt.args.ctx)
				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			},
		)
	}
}

func TestNewPool(t *testing.T) {
	type args struct {
		delegate *goRedis.Client
	}
	goRedisOpt := &goRedis.Options{
		Network:  "tcp",
		Addr:     "192.168.0.128:6379",
		Password: "yourpassword",
		DB:       0,
	}
	client := goRedis.NewClient(goRedisOpt)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test01",
			args:    args{delegate: client},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewPool(tt.args.delegate); !reflect.DeepEqual(got, nil) == tt.wantErr {
					t.Errorf("NewPool() = %v", got)
				}
			},
		)
	}
}

func Test_conn_SetNX(t *testing.T) {
	type fields struct {
		delegate *goRedis.Client
		ctx      context.Context
	}
	type args struct {
		key   string
		value string
		ttl   time.Duration
	}
	ctx := context.TODO()
	goRedisOpt := &goRedis.Options{
		Network:  "tcp",
		Addr:     "192.168.0.128:6379",
		Password: "yourpassword",
		DB:       0,
	}
	client := goRedis.NewClient(goRedisOpt)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "test01",
			fields: fields{delegate: client, ctx: ctx},
			args: args{
				key:   "test:goredis:nx",
				value: "1",
				ttl:   time.Second,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &conn{
					delegate: tt.fields.delegate,
					ctx:      tt.fields.ctx,
				}
				got, err := c.SetNX(tt.args.key, tt.args.value, tt.args.ttl)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetNX() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("SetNX() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_conn_Expire(t *testing.T) {
	type fields struct {
		delegate *goRedis.Client
		ctx      context.Context
	}
	type args struct {
		key string
		ttl time.Duration
	}
	ctx := context.TODO()
	goRedisOpt := &goRedis.Options{
		Network:  "tcp",
		Addr:     "192.168.0.128:6379",
		Password: "yourpassword",
		DB:       0,
	}
	client := goRedis.NewClient(goRedisOpt)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "test01",
			fields: fields{delegate: client, ctx: ctx},
			args: args{
				key: "test:goredis:expire",
				ttl: time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &conn{
					delegate: tt.fields.delegate,
					ctx:      tt.fields.ctx,
				}
				got, err := c.Expire(tt.args.key, tt.args.ttl)
				if (err != nil) != tt.wantErr {
					t.Errorf("Expire() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("Expire() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_conn_Del(t *testing.T) {
	type fields struct {
		delegate *goRedis.Client
		ctx      context.Context
	}
	type args struct {
		key string
	}
	ctx := context.TODO()
	goRedisOpt := &goRedis.Options{
		Network:  "tcp",
		Addr:     "192.168.0.128:6379",
		Password: "yourpassword",
		DB:       0,
	}
	client := goRedis.NewClient(goRedisOpt)
	client.Set(ctx, "test_del_key1", "1", time.Second)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "test01",
			fields:  fields{delegate: client, ctx: ctx},
			args:    args{key: "test_del_key"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "test02",
			fields:  fields{delegate: client, ctx: ctx},
			args:    args{key: "test_del_key1"},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &conn{
					delegate: tt.fields.delegate,
					ctx:      tt.fields.ctx,
				}
				got, err := c.Del(tt.args.key)
				if (err != nil) != tt.wantErr {
					t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("Del() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
