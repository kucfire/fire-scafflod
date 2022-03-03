package cache

import (
	"math/rand"
	"strings"
	"time"

	"fire-scaffold/conf"
	"fire-scaffold/pkg/trace"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type Option func(*option)

type Trace = trace.T

type option struct {
	Trace *trace.Trace
	Redis *trace.Redis
}

func newOption() *option {
	return &option{}
}

var client *redis.Client

type Cache interface {
	i()
	Set(key, value string, ttl time.Duration, options ...Option) error
	Get(key string, options ...Option) (string, error)
	TTL(key string) (time.Duration, error)
	Expire(key string, ttl time.Duration) bool
	ExpireAt(key string, ttl time.Time) bool
	Del(key string, options ...Option) bool
	Exists(keys ...string) bool
	Incr(key string, options ...Option) int64
	Version() string
}

type cache struct {
	client *redis.Client
	rand   *rand.Rand
}

func InitRedis(conf conf.Redis) error {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Pass,
		DB:           conf.Db,
		MaxRetries:   conf.MaxRetries,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConns,
	})

	if err := client.Ping().Err(); err != nil {
		return errors.Wrap(err, "ping redis err")
	}

	return nil
}

func Close() error {
	return client.Close()
}

func New() Cache {
	return &cache{
		client: client,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *cache) i()

func (c *cache) Set(key, value string, ttl time.Duration, options ...Option) error {
	now := time.Now()

	opt := newOption()

	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = time.Now().Unix()
			opt.Redis.Optional = trace.SetOptional
			opt.Redis.Key = key
			opt.Redis.Value = value
			opt.Redis.TTL = ttl.Seconds()
			opt.Redis.CostSeconds = time.Since(now).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	if err := c.client.Set(key, value, ttl).Err(); err != nil {
		return errors.Wrapf(err, "redis set key: %s err", key)
	}

	return nil
}

func (c *cache) Get(key string, options ...Option) (string, error) {
	now := time.Now()
	opt := newOption()

	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = time.Now().Unix()
			opt.Redis.Optional = trace.GetOptional
			opt.Redis.Key = key
			opt.Redis.CostSeconds = time.Since(now).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	result, err := c.client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *cache) TTL(key string) (time.Duration, error) {
	ttl, err := c.client.TTL(key).Result()
	if err != nil {
		return -1, errors.Wrapf(err, "redis get key: %s error", key)
	}

	return ttl, nil
}

func (c *cache) Expire(key string, ttl time.Duration) bool {
	ok, _ := c.client.Expire(key, ttl).Result()

	return ok
}

func (c *cache) ExpireAt(key string, ttl time.Time) bool {
	ok, _ := c.client.ExpireAt(key, ttl).Result()

	return ok
}

func (c *cache) Del(key string, options ...Option) bool {
	now := time.Now()
	opt := newOption()

	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = time.Now().Unix()
			opt.Redis.Optional = trace.DeleteOptional
			opt.Redis.Key = key
			opt.Redis.CostSeconds = time.Since(now).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	value, _ := c.client.Del(key).Result()

	return value > 0
}

func (c *cache) Exists(keys ...string) bool {
	value, _ := c.client.Exists(keys...).Result()

	return value > 0
}

func (c *cache) Incr(key string, options ...Option) int64 {
	now := time.Now()
	opt := newOption()

	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = time.Now().Unix()
			opt.Redis.Optional = trace.IncreaseOptional
			opt.Redis.Key = key
			opt.Redis.CostSeconds = time.Since(now).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	value, _ := c.client.Incr(key).Result()

	return value
}

func WithTrace(t Trace) Option {
	return func(opt *option) {
		if t != nil {
			opt.Trace = t.(*trace.Trace)
			opt.Redis = new(trace.Redis)
		}
	}
}

func (c *cache) Version() string {
	server := c.client.Info("server").Val()
	spl1 := strings.Split(server, "# Server")
	spl2 := strings.Split(spl1[1], "redis_version:")
	spl3 := strings.Split(spl2[1], "redis_git_sha1:")
	return spl3[0]
}
