package db

import (
	"net/url"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/revel/revel"
)

type (
	Redis interface {
		Float64(key string) (float64, error)
		Get(key string) (interface{}, error)
		GetString(key string) (string, error)
		Int64(key string) (int64, error)
		IsErrNil(err error) bool
		Set(key string, val interface{}) (interface{}, error)
		SetEx(key string, ttl time.Duration, val interface{}) (interface{}, error)
		String(key string) (string, error)
	}

	RedisConfig struct {
		IdleTimeout time.Duration
		MaxActive   int
		MaxIdle     int
	}

	AppRedis struct {
		pool *redis.Pool
	}
)

func NewRedisProvider(
	config *RedisConfig,
) *AppRedis {
	return NewRedisProviderWithParams(
		revel.Config.StringDefault("redis.url", ""),
		config,
	)
}

func NewRedisProviderWithParams(
	urlStr string,
	config *RedisConfig,
) *AppRedis {
	idleTimeout := 2 * time.Minute
	maxActive := 200
	maxIdle := 5

	if config != nil {
		if int64(config.IdleTimeout) != 0 {
			idleTimeout = config.IdleTimeout
		}

		if config.MaxActive != 0 {
			maxActive = config.MaxActive
		}

		if config.MaxIdle != 0 {
			maxIdle = config.MaxIdle
		}
	}

	redisPool := &redis.Pool{
		IdleTimeout: idleTimeout,
		MaxActive:   maxActive,
		MaxIdle:     maxIdle,
		Dial: func() (redis.Conn, error) {
			user, err := url.Parse(urlStr)
			if err != nil {
				panic(err.Error())
			}
			if user.User == nil {
				panic("redis user is nil")
			}
			password, ok := user.User.Password()
			if !ok {
				panic("could not get redis password")
			}

			c, err := redis.Dial(
				"tcp",
				user.Host,
			)
			if err != nil {
				panic(err.Error())
			}

			_, err = c.Do("AUTH", password)
			if err != nil {
				panic(err.Error())
			}

			return c, nil
		},
	}

	return &AppRedis{
		pool: redisPool,
	}
}

func (p *AppRedis) Float64(
	key string,
) (float64, error) {
	return redis.Float64(p.Get(key))
}

func (p *AppRedis) Get(
	key string,
) (interface{}, error) {
	return p.do("GET", key)
}

func (p *AppRedis) GetString(key string) (string, error) {
	return redis.String(p.Get(key))
}

func (p *AppRedis) Int64(
	key string,
) (int64, error) {
	return redis.Int64(p.Get(key))
}

func (p *AppRedis) IsErrNil(
	err error,
) bool {
	return err == redis.ErrNil
}

func (p *AppRedis) RedisPool() *redis.Pool {
	return p.pool
}

func (p *AppRedis) Set(
	key string,
	val interface{},
) (interface{}, error) {
	return p.do("SET", key, val)
}

func (p *AppRedis) SetEx(
	key string,
	ttl time.Duration,
	val interface{},
) (interface{}, error) {
	return p.do("SETEX", key, ttl.Seconds(), val)
}

func (p *AppRedis) String(
	key string,
) (string, error) {
	return redis.String(p.Get(key))
}

func (p *AppRedis) do(
	commandName string,
	args ...interface{},
) (interface{}, error) {
	conn := p.pool.Get()
	defer conn.Close()

	return conn.Do(commandName, args...)
}
