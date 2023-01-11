package main

//STRING TO JSON https://stackoverflow.com/questions/28859941/whats-the-golang-equivalent-of-converting-any-json-to-standard-dict-in-python
import (
	"collector/modules"
	"collector/test"
	"github.com/go-redis/redis"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"os"
)

func getRedis(config modules.RedisConfig) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func getSocketListen(config modules.ZmqConfig) *zmq4.Socket {
	server, _ := zmq4.NewSocket(zmq4.PULL)
	err := server.Bind(config.ConnType + "://" + config.Host + ":" + config.Port)
	modules.CheckErr(err)
	return server
}

func writeToRedis(message string, redisConn *redis.Client, redisStream modules.RedisStream) {
	_, err := redisConn.XAdd(&redis.XAddArgs{
		Stream: redisStream.Stream,
		MaxLen: 0,
		ID:     "*",
		Values: map[string]interface{}{
			"message":    message,
			"ticketID":   int(rand.Intn(100000000)),
			"ticketData": string("some ticket data"),
		},
	}).Result()
	modules.CheckErr(err)
}

func main() {
	file_path := os.Args[1]

	viper.AddConfigPath(file_path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	var configuration modules.Configuration
	err := viper.ReadInConfig()
	modules.CheckErr(err)
	err = viper.Unmarshal(&configuration)
	modules.CheckErr(err)
	sock := getSocketListen(configuration.ZmqConfig)
	redisConn := getRedis(configuration.RedisConfig)
	println("We are here")
	if configuration.RedisStream.Dummy == true {
		go test.SendMsg()
	}
	for {
		message, _ := sock.RecvMessage(0)
		println(message[0])
		go writeToRedis(message[0], redisConn, configuration.RedisStream)
	}
}
