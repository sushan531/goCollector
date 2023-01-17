package redisLib

import (
	"github.com/go-redis/redis"
	"loadHandler/modules"
	"log"
	"strconv"
	"syscall"
	"time"
)

func GetRedis(conn modules.RedisConn) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     conn.Host + ":" + conn.Port,
		Password: conn.Password,
		DB:       conn.Database,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func WriteToRedisStream(message string, redisConn *redis.Client) {
	_, _ = redisConn.XAdd(&redis.XAddArgs{
		Stream: "testStream",
		MaxLen: 0,
		ID:     "*",
		//Values: message,
		Values: map[string]interface{}{
			"message":    message,
			"timestamp":  time.Now(),
			"ticketData": "some ticket data",
		},
	}).Result()

}

func WriteToRedisString(value int, redisConn *redis.Client) {
	_, _ = redisConn.Set("Benchmark"+strconv.Itoa(syscall.Getpid()), value, 0).Result()
}
