package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
	"loadBalancer/modules"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{}
var LoadHandlers []*os.Process

type Handler struct {
	*zmq4.Socket
}

func spawnLoadHandlers(minHandlers int) {
	for i := 0; i < minHandlers; i++ {
		var procAttr os.ProcAttr
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr} // todo make it write to file
		process, err := os.StartProcess(
			"/Users/spt/GolandProjects/collector/loadBalancerSrc/loadHandler/loadHandler",
			[]string{"", "/Users/spt/GolandProjects/collector/loadBalancerSrc/loadHandler/"},
			&procAttr,
		)
		if err != nil {
			log.Fatal(err)
		}
		LoadHandlers = append(LoadHandlers, process)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		_, _ = h.SendMessage(message)

	}
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
	config := readConfig(filePath)
	spawnLoadHandlers(config.MinHandlers)

	for _, value := range LoadHandlers {
		fmt.Println("Process Id : ", value)
	}

	var handler Handler
	server, _ := zmq4.NewSocket(zmq4.PUSH)
	_ = server.Bind(config.ZmqConfig.Network + "://" + config.ZmqConfig.Host + ":" + config.ZmqConfig.Port)
	handler.Socket = server

	var addr = flag.String("addr", config.WebSocketConfig.Host+":"+config.WebSocketConfig.Port, "http service address")
	http.Handle(config.WebSocketConfig.Path, &handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
