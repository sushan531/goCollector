package main

import (
	"flag"
	"fmt"
	"github.com/pebbe/zmq4"
	"loadBalancer/modules"
	"loadBalancer/modules/endpoints"
	"loadBalancer/modules/handlers"
	"log"
	"net/http"
	"os"
)

func managerEntryPointServer(config modules.WebSocketConfig, socket *zmq4.Socket) {
	hostAndPort := config.Host + ":" + config.Port
	fmt.Println("Running server on ... " + hostAndPort)
	handler := endpoints.LogHandler(socket)
	http.Handle(config.Path, &handler)
	var addr = flag.String("addr", hostAndPort, "http service address")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func bindZmqSocket(config modules.ZmqConfig) *zmq4.Socket {
	zmqConnStr := config.Network + "://" + config.Host + ":" + config.Port
	fmt.Println("Binding a ZMQ socket Type: PUSH, " + "CONN: " + zmqConnStr)
	zmqSocket, _ := zmq4.NewSocket(zmq4.PUSH)
	_ = zmqSocket.Bind(zmqConnStr)
	return zmqSocket
}

func main() {
	filePath := os.Args[1]
	config := modules.ReadConfig(filePath)

	handlers.SpawnLoadHandlers(config.MinHandlers)
	zmqSocket := bindZmqSocket(config.ZmqConfig)
	managerEntryPointServer(config.WebSocketConfig, zmqSocket)

}
