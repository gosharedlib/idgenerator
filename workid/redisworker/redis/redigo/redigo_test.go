package redigo

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func Test_conn_Close(t *testing.T) {
	type fields struct {
		delegate redis.Conn
	}
	rediGoPool := redis.Pool{
		MaxIdle:     1024,
		MaxActive:   60000,
		IdleTimeout: time.Minute,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp", "192.168.0.128:6379",
				redis.DialConnectTimeout(time.Millisecond*200),
				redis.DialReadTimeout(time.Millisecond*500),
				redis.DialWriteTimeout(time.Millisecond*500),
				redis.DialPassword("yourpassword"),
				redis.DialDatabase(0),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "test_01",
			fields:  fields{delegate: rediGoPool.Get()},
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
			args:    args{err: redis.ErrNil},
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
		delegate *redis.Pool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "text_01",
			fields:  fields{delegate: &redis.Pool{}},
			args:    args{ctx: context.TODO()},
			wantErr: true,
		},
		{
			name:    "text_02",
			fields:  fields{delegate: &redis.Pool{}},
			args:    args{ctx: nil},
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
					t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			},
		)
	}
}

func TestNewPool(t *testing.T) {
	type args struct {
		delegate *redis.Pool
	}
	rediGoPool := redis.Pool{
		MaxIdle:     1024,
		MaxActive:   60000,
		IdleTimeout: time.Minute,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp", "192.168.0.128:6379",
				redis.DialConnectTimeout(time.Millisecond*200),
				redis.DialReadTimeout(time.Millisecond*500),
				redis.DialWriteTimeout(time.Millisecond*500),
				redis.DialPassword("yourpassword"),
				redis.DialDatabase(0),
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
			name: "test01",
			args: args{
				delegate: &rediGoPool,
			},
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
		delegate redis.Conn
	}
	rediGoConn, _ := redis.Dial(
		"tcp", "192.168.0.128:6379",
		redis.DialConnectTimeout(time.Millisecond*200),
		redis.DialReadTimeout(time.Millisecond*500),
		redis.DialWriteTimeout(time.Millisecond*500),
		redis.DialPassword("yourpassword"),
		redis.DialDatabase(0),
	)

	type args struct {
		key   string
		value string
		ttl   time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "test01",
			fields: fields{delegate: rediGoConn},
			args: args{
				key:   "test:redigo:nx",
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
		delegate redis.Conn
	}
	rediGoConn, _ := redis.Dial(
		"tcp", "192.168.0.128:6379",
		redis.DialConnectTimeout(time.Millisecond*200),
		redis.DialReadTimeout(time.Millisecond*500),
		redis.DialWriteTimeout(time.Millisecond*500),
		redis.DialPassword("yourpassword"),
		redis.DialDatabase(0),
	)
	type args struct {
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
			name: "test01",
			fields: fields{
				delegate: rediGoConn,
			},
			args: args{
				key: "test:redigo:expire",
				ttl: time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &conn{
					delegate: tt.fields.delegate,
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
		delegate redis.Conn
	}
	type args struct {
		key string
	}
	rediGoConn, _ := redis.Dial(
		"tcp", "192.168.0.128:6379",
		redis.DialConnectTimeout(time.Millisecond*200),
		redis.DialReadTimeout(time.Millisecond*500),
		redis.DialWriteTimeout(time.Millisecond*500),
		redis.DialPassword("yourpassword"),
		redis.DialDatabase(0),
	)
	_, _ = rediGoConn.Do("SET", "test_del_key1", "1", "EX", 1)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "test01",
			fields:  fields{delegate: rediGoConn},
			args:    args{key: "test_del_key"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "test01",
			fields:  fields{delegate: rediGoConn},
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
