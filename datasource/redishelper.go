package datasource

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"lottery/conf"
	"sync"
	"time"
)

var redisLock sync.Mutex
var cacheInstance *RedisConn

type RedisConn struct {
	pool      *redis.Pool
	showDebug bool
}

func (redis *RedisConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {

	conn := redis.pool.Get()
	defer conn.Close()

	t1 := time.Now().UnixNano()

	reply, err = conn.Do(commandName, args...)
	if err != nil {
		e := conn.Err()
		if e != nil {
			log.Println("redishelper.Do() error", err, e)
		}
	}

	t2 := time.Now().UnixNano()
	if redis.showDebug {
		fmt.Printf("[redis] [info] [%dus] command=%s, error=%s, args=%s, reply=%s\n",
			(t2-t1)/1000, commandName, err, args, reply)
	}

	return reply, err
}

func (redis *RedisConn) ShowDebug(b bool) {

	redis.showDebug = b
}

func InstanceCache() *RedisConn {

	if cacheInstance != nil {
		return cacheInstance
	}

	redisLock.Lock()
	defer redisLock.Unlock()

	if cacheInstance != nil {
		return cacheInstance
	}

	return NewCache()
}

func NewCache() *RedisConn {

	pool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", conf.RedisCache.Host, conf.RedisCache.Port))
			if err != nil {
				log.Fatal("redishelper.NewCache().Dial() error", err)
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         10000,
		MaxActive:       10000,
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}

	instance := &RedisConn{
		pool: &pool,
	}
	cacheInstance = instance
	instance.ShowDebug(true)

	return instance
}
