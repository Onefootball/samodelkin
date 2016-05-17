// Package redis provides an abstraction
// for handling redis connection logic
package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

type (
	// RedisConfig is a struct type
	// which holds redis connector
	// pool configurations
	RedisConfig struct {
		Host           string
		Port           int
		MaxIdleConns   int
		MaxActiveConns int
		IdleTimeout    time.Duration
		ConnectTimeout time.Duration
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
	}
)

// Address returns redis connection string
func (rc *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", rc.Host, rc.Port)
}

type (
	// RedisConnector is an interface type
	// which describes methods for retrieving
	// a redis connection
	RedisConnector interface {
		Connect() redis.Conn
		PingConnect() (redis.Conn, error)
	}

	// DBConnector is a struct type
	// which implements a RedisConnector interface
	//
	// DBConnector can hold multiple redis connection pools
	// this provides a possibility to operate with multiple
	// redis slave from one connector
	DBConnector struct {
		pools         []*redis.Pool
		activePoolIdx uint
	}
)

// NewDBConnector inits and returns a pointer to DBConnector instance
func NewDBConnector(cfgs []RedisConfig) (*DBConnector, error) {
	if len(cfgs) == 0 {
		return nil, errors.New("invalid argument: cfgs")
	}

	pools := make([]*redis.Pool, 0, len(cfgs))
	for _, cfg := range cfgs {
		pools = append(pools, newPool(cfg))
	}

	return &DBConnector{pools: pools}, nil
}

// Connect returns an awailable redis connection
// from a random pool.
// If connection is by any reason broken, error
// will be returned on first attempt to use the connection
// In cases when it's needed to check if the connection
// is alive, use PingConnect method
func (c *DBConnector) Connect() redis.Conn {
	return c.pickPool().Get()
}

// PingConnect returns an awailable redis connection
// from a random pool.
// Second arguments is an error generated by a "PING"
// method sent to a retrieved connection
func (c *DBConnector) PingConnect() (redis.Conn, error) {
	redisConn := c.pickPool().Get()
	return redisConn, pingConn(redisConn, time.Now())
}

// pickPool returns the next available pool from a RedisConnector
// pools are picked up based on round-robin pattern - one by one.
// There is NO logic involved for picking up a pool based on balancing
// awailable connections in the pool
func (c *DBConnector) pickPool() *redis.Pool {
	if len(c.pools) == 1 {
		return c.pools[0]
	}

	defer func() {
		if uint(len(c.pools)-1) > c.activePoolIdx {
			c.activePoolIdx++
		} else {
			c.activePoolIdx = 0
		}
	}()

	return c.pools[c.activePoolIdx]
}

// PingConn pings redis connection and returns an error
// for a failed connection
func pingConn(c redis.Conn, _ time.Time) error {
	_, err := c.Do("PING")
	return err
}

// newPool returns a pointer to redis pool instance
func newPool(cfg RedisConfig) *redis.Pool {
	return &redis.Pool{
		Dial:         dialFunc(cfg),
		TestOnBorrow: pingConn,
		MaxIdle:      cfg.MaxIdleConns,
		MaxActive:    cfg.MaxActiveConns,
		IdleTimeout:  cfg.IdleTimeout,
	}
}

// dialFunc returns a func which handles establishing a
// redis connection
func dialFunc(cfg RedisConfig) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		return redis.DialTimeout(
			"tcp",
			cfg.Address(),
			cfg.ConnectTimeout,
			cfg.ReadTimeout,
			cfg.WriteTimeout,
		)
	}
}