package main

//STRING TO JSON https://stackoverflow.com/questions/28859941/whats-the-golang-equivalent-of-converting-any-json-to-standard-dict-in-python
import (
	"github.com/go-redis/redis"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
	"loadHandler/modules"
	"loadHandler/modules/redisLib"
	"loadHandler/modules/zmqLib"
	"os"
	"time"
)

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

func handlerLoop(socket *zmq4.Socket, client *redis.Client) {
	count := 0
	t1 := time.Now()
	t2 := t1.Add(time.Second * 10)
	for {
		message := zmqLib.GetMessage(socket)
		redisLib.WriteToRedisStream(message[0], client)
		count += 1
		if time.Now() == t2 {
			redisLib.WriteToRedisString(count/10, client)
			count = 0
			t2 = t2.Add(time.Second * 10)
		}

	}
}

func main() {
	filePath := os.Args[1]
	configuration := readConfig(filePath)
	sock := zmqLib.GetSocketListen(configuration.ZmqPull)
	client := redisLib.GetRedis(configuration.RedisConn)
	handlerLoop(sock, client)
}
