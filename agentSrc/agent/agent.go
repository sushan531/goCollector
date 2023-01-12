package main

//STRING TO JSON https://stackoverflow.com/questions/28859941/whats-the-golang-equivalent-of-converting-any-json-to-standard-dict-in-python
import (
	"agent/modules"
	"agent/test"
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

func readConfig(filePath string) modules.Configuration {
	var configuration modules.Configuration

	viper.AddConfigPath(filePath)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	modules.CheckErr(err)
	err = viper.Unmarshal(&configuration)
	modules.CheckErr(err)
	return configuration
}

func main() {
	filePath := os.Args[1]
	configuration := readConfig(filePath)
	redisConn := getRedis(configuration.RedisConfig)

	var sock *zmq4.Socket
	if configuration.RedisStream.Dummy == true {
		go test.SendMsg()
		sock = getSocketListen(configuration.ZmqConfig)
	} else {
		sock = getSocketListen(configuration.ZmqConfig)
	}
	for {
		message, _ := sock.RecvMessage(0)
		println(message[0])
		// todo : send to a rest api or a websocket
		// go sendToWebSocket()
		go writeToRedis(message[0], redisConn, configuration.RedisStream)
	}
}
