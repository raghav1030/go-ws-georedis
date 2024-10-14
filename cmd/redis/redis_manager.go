package redis_manager

import (
	"log"
	"sync"

	"github.com/gomodule/redigo/redis"
)

// Singleton RedisManager
type RedisManager struct {
	pool *redis.Pool
}

var instance *RedisManager
var once sync.Once

// GetRedisManager returns the RedisManager singleton
func GetRedisManager() *RedisManager {
	once.Do(func() {
		instance = &RedisManager{
			pool: &redis.Pool{
				Dial: func() (redis.Conn, error) {
					return redis.Dial("tcp", "redis-geospatial:6379")
				},
			},
		}
	})
	return instance
}

// AddUserLocation stores the user's location in Redis
func (rm *RedisManager) AddUserLocation(userID string, lat, lon float64) error {
	conn := rm.pool.Get()
	defer conn.Close()

	reply , err := conn.Do("GEOADD", "users:locations", lon, lat, userID)
	log.Println(reply)
	return err
}

// GetNearbyUsers retrieves users near a given user within a specified radius
func (rm *RedisManager) GetNearbyUsers(userID string, radius float64) ([]string, error) {
	conn := rm.pool.Get()
	defer conn.Close()

	users, err := redis.Strings(conn.Do("GEORADIUSBYMEMBER", "users:locations", userID, radius, "km"))
	if err != nil {
		return nil, err
	}
	return users, nil
}
