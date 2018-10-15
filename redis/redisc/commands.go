package redisc

import (
	"fmt"
	"strings"
	"time"

	redis "github.com/gomodule/redigo/redis"
	redisc "github.com/mna/redisc"
)

// Del do DEL command to redis
// Cluster does not support multi command, will invesitage later
func (c *Cluster) Del(key string) error {
	client := c.clusterPool.Get()
	defer client.Close()

	if client.Err() != nil {
		return client.Err()
	}

	rc, err := redisc.RetryConn(client, c.retryCount, c.retryDuration*time.Millisecond)
	if err != nil {
		return err
	}

	resp, err := redis.Int(rc.Do("DEL", key))
	if err != nil {
		return err
	}

	if resp < 0 {
		return fmt.Errorf("Unexpected redis response %d", resp)
	}

	return nil
}

// Get do GET command to redis
func (c *Cluster) Get(key string) (string, error) {
	client := c.clusterPool.Get()
	defer client.Close()

	if client.Err() != nil {
		return "", client.Err()
	}

	rc, err := redisc.RetryConn(client, c.retryCount, c.retryDuration*time.Millisecond)
	if err != nil {
		return "", err
	}

	resp, err := redis.String(rc.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return "", err
	}

	return resp, nil
}

// Setex do SETEX command to redis
func (c *Cluster) Setex(key string, value string, ttl int) error {
	client := c.clusterPool.Get()
	defer client.Close()

	if client.Err() != nil {
		return client.Err()
	}

	rc, err := redisc.RetryConn(client, c.retryCount, c.retryDuration*time.Millisecond)
	if err != nil {
		return err
	}

	// set default TTL
	if ttl == 0 {
		ttl = 3600
	}

	resp, err := redis.String(rc.Do("SET", key, value, "EX", ttl))
	if err != nil {
		return err
	}

	if !strings.EqualFold("OK", resp) {
		return fmt.Errorf("Unexpected redis response %s", resp)
	}

	return nil
}

// HGet do HGET command to redis
func (c *Cluster) HGet(key, field string) (string, error) {
	client := c.clusterPool.Get()
	defer client.Close()

	if client.Err() != nil {
		return "", client.Err()
	}

	rc, err := redisc.RetryConn(client, c.retryCount, c.retryDuration*time.Millisecond)
	if err != nil {
		return "", err
	}

	resp, err := redis.String(rc.Do("HGET", key, field))
	if err != nil && err != redis.ErrNil {
		return "", err
	}

	return resp, nil
}

// HSet do HSET command to redis
func (c *Cluster) HSet(key string, field string, value string, ttl int) error {
	client := c.clusterPool.Get()
	defer client.Close()

	if client.Err() != nil {
		return client.Err()
	}

	rc, err := redisc.RetryConn(client, c.retryCount, c.retryDuration*time.Millisecond)
	if err != nil {
		return err
	}

	// do HSET
	resp, err := redis.Int(rc.Do("HSET", key, field, value))
	if err != nil {
		return err
	}

	if resp < 0 || resp > 1 {
		return fmt.Errorf("Unexpected redis response %d", resp)
	}

	// set default TTL
	if ttl == 0 {
		ttl = 3600
	}

	// do EXPIRE
	respExpire, err := redis.Int(rc.Do("EXPIRE", key, ttl))
	if err != nil {
		return err
	}

	if respExpire < 0 || respExpire > 1 {
		return fmt.Errorf("Unexpected redis response %d", respExpire)
	}

	return nil
}

// HMGet do HMGET command to redis
func (c *Cluster) HMGet(key string, fields []string) ([]string, error) {
	client := c.clusterPool.Get()
	defer client.Close()

	if client.Err() != nil {
		return []string{}, client.Err()
	}

	rc, err := redisc.RetryConn(client, c.retryCount, c.retryDuration*time.Millisecond)
	if err != nil {
		return []string{}, err
	}

	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, v := range fields {
		args[i+1] = v
	}

	resp, err := redis.Strings(rc.Do("HMGET", args...))
	if err != nil && err != redis.ErrNil {
		return []string{}, err
	}

	return resp, nil
}

// HMSet do HMSET command to redis
func (c *Cluster) HMSet(key string, ttl int, m map[string]string) error {
	client := c.clusterPool.Get()
	defer client.Close()

	if client.Err() != nil {
		return client.Err()
	}

	rc, err := redisc.RetryConn(client, c.retryCount, c.retryDuration*time.Millisecond)
	if err != nil {
		return err
	}

	// do HMSET
	resp, err := redis.String(rc.Do("HMSET", redis.Args{}.Add(key).AddFlat(m)...))
	if err != nil {
		return err
	}

	if !strings.EqualFold("OK", resp) {
		return fmt.Errorf("Unexpected redis response %s", resp)
	}

	// set default TTL
	if ttl == 0 {
		ttl = 3600
	}

	// do EXPIRE
	respExpire, err := redis.Int(rc.Do("EXPIRE", key, ttl))
	if err != nil {
		return err
	}

	if respExpire < 0 || respExpire > 1 {
		return fmt.Errorf("Unexpected redis response %d", respExpire)
	}

	return nil
}
