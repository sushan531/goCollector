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
	"os/signal"
	"syscall"
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

func graceFullShutdown(signals []syscall.Signal, handlers []*os.Process, socket *zmq4.Socket) {
	// https://www.developer.com/languages/os-signals-go/
	// Also refer above doc may be there is sth there.
	// Currently system don't handle kill -9 i.e SIGKILL
	var gracefulStop = make(chan os.Signal, 1)
	for _, sig := range signals {
		signal.Notify(gracefulStop, sig)
	}
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		_ = socket.Close()
		for _, handler := range handlers {
			_ = handler.Kill()
		}
		os.Exit(0)

	}()
}

func main() {
	filePath := os.Args[1]
	config := modules.ReadConfig(filePath)
	handlersList := handlers.SpawnLoadHandlers(config.MinHandlers)
	zmqSocket := bindZmqSocket(config.ZmqConfig)
	signalsToHandle := []syscall.Signal{
		syscall.SIGSTOP, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	}
	graceFullShutdown(signalsToHandle, handlersList, zmqSocket)
	managerEntryPointServer(config.WebSocketConfig, zmqSocket)
}
