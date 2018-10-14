package redisc

import (
	"encoding/json"
	"time"

	redis "github.com/gomodule/redigo/redis"
	redisc "github.com/mna/redisc"
)

var (
	redisCfg *Config
)

// Cluster represents cluster client
type Cluster struct {
	clusterPool   *redisc.Cluster
	retryCount    int
	retryDuration time.Duration
}

func genDialOption() []redis.DialOption {
	options := []redis.DialOption{}

	if redisCfg.DialConnectTimeout != 0 {
		options = append(options, redis.DialConnectTimeout(time.Duration(redisCfg.DialConnectTimeout)*time.Second))
	}

	if redisCfg.DialWriteTimeout != 0 {
		options = append(options, redis.DialWriteTimeout(time.Duration(redisCfg.DialWriteTimeout)*time.Second))
	}

	if redisCfg.DialReadTimeout != 0 {
		options = append(options, redis.DialReadTimeout(time.Duration(redisCfg.DialReadTimeout)*time.Second))
	}

	if redisCfg.DialDatabase != 0 {
		options = append(options, redis.DialDatabase(redisCfg.DialDatabase))
	}

	if redisCfg.DialKeepAlive != 0 {
		options = append(options, redis.DialKeepAlive(time.Duration(redisCfg.DialKeepAlive)*time.Second))
	}

	if redisCfg.DialPassword != "" {
		options = append(options, redis.DialPassword(redisCfg.DialPassword))
	}

	return options
}

func createPool(address string, opts ...redis.DialOption) (*redis.Pool, error) {
	rpool := &redis.Pool{
		MaxActive:   redisCfg.MaxActive,
		MaxIdle:     redisCfg.MaxIdle,
		IdleTimeout: time.Duration(redisCfg.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address, opts...)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Second {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	if _, err := rpool.Dial(); err != nil {
		rpool.Close()
		return nil, err
	}

	return rpool, nil
}

// InitCluster initialize cluster client
func InitCluster(cfg *Config) (*Cluster, error) {
	redisCfg = cfg

	pool := redisc.Cluster{
		StartupNodes: []string{redisCfg.Host},
		DialOptions:  genDialOption(),
		CreatePool:   createPool,
	}

	if err := pool.Refresh(); err != nil {
		return nil, err
	}

	cluster := &Cluster{
		clusterPool:   &pool,
		retryCount:    redisCfg.RetryCount,
		retryDuration: time.Duration(redisCfg.RetryDuration),
	}

	return cluster, nil
}

// ClusterStatus get cluster status
func (c *Cluster) ClusterStatus() ([]byte, error) {
	return json.MarshalIndent(c.clusterPool.Stats(), "", "    ")
}
