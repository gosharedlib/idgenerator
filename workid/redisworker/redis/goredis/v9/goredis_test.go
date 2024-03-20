package goredis

import (
	"context"
	"reflect"
	"testing"

	goRedis "github.com/redis/go-redis/v9"
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
