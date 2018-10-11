package redisc

import (
	"log"

	redis "github.com/gomodule/redigo/redis"
	redisc "github.com/mna/redisc"
)

var (
	redisCfg *Config
)

func genDialOption() []redis.DialOption {
	options := []redis.DialOption{}

	if rediscfg.DialConnectTimeout != 0 {
		options = append(options, redisCfg.DialConnectTimeout)
	}

	if rediscfg.DialWriteTimeout != 0 {
		options = append(options, redisCfg.DialWriteTimeout)
	}

	if rediscfg.DialReadTimeout != 0 {
		options = append(options, redisCfg.DialWriteTimeout)
	}

	if rediscfg.DialDatabase != 0 {
		options = append(options, redisCfg.DialDatabase)
	}

	if rediscfg.DialKeepAlive != 0 {
		options = append(options, redisCfg.DialKeepAlive)
	}

	if rediscfg.DialPassword != "" {
		options = append(options, redisCfg.DialPassword)
	}

	return options
}

func createPool(address string, opts ...redis.DialOption) (*redis.Pool, error) {
    rpool := &redis.Pool{
        MaxActive:   redisCfg.MaxActive,
        MaxIdle:     redisCfg.MaxIdle,
        IdleTimeout: redisCfg.IdleTimeout * time.Second,
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

func InitRedisCluster(cfg *Config) *redisc.Cluster, error {
	rediscfg = cfg

	rdsCluster := redisc.Cluster{
		StartupNodes: []string{rediscfg.Host},
		DialOptions:  genDialOption(),
		CreatePool:   createPool,
	}

	if err := rdsCluster.Refresh(); err != nil {
		return nil, err
    }
    
    return rdsCluster, nil
}
