package main

//STRING TO JSON https://stackoverflow.com/questions/28859941/whats-the-golang-equivalent-of-converting-any-json-to-standard-dict-in-python
import (
	"github.com/pebbe/zmq4"
	//"github.com/spf13/viper"
	"log"
	//"os"
	"syscall"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getSocketListen() *zmq4.Socket {
	server, _ := zmq4.NewSocket(zmq4.PULL)
	err := server.Connect("tcp://127.0.0.1:5556")
	CheckErr(err)
	return server
}

//func readConfig(filePath string) modules.Configuration {
//	var configuration modules.Configuration
//
//	viper.AddConfigPath(filePath)
//	viper.SetConfigName("config")
//	viper.SetConfigType("yml")
//
//	err := viper.ReadInConfig()
//	modules.CheckErr(err)
//	err = viper.Unmarshal(&configuration)
//	modules.CheckErr(err)
//	return configuration
//}

func main() {
	//filePath := os.Args[1]
	//configuration := readConfig(filePath)

	sock := getSocketListen(
	//configuration.ZmqConfig
	)
	for {
		message, _ := sock.RecvMessage(0)
		println(syscall.Getpid(), message[0])
	}
}
