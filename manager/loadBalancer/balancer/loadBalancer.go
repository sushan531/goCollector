package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/docker/docker/client"
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

func graceFullShutdown(signals []syscall.Signal, loadHandlerList []string, socket *zmq4.Socket, cli *client.Client, ctx context.Context) {
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
		println("Killing Docker containers")
		for _, loadHandlerId := range loadHandlerList {
			handlers.KillDockerContainer(cli, loadHandlerId, ctx)
		}
		handlers.RemoveDockerContainers(cli, loadHandlerList, ctx)
		os.Exit(0)

	}()
}

func main() {
	filePath := os.Args[1]
	config := modules.ReadConfig(filePath)
	cli, ctx := handlers.GetDockerClient()
	handlersList := handlers.SpawnLoadHandlers(config.MinHandlers, cli, ctx)
	zmqSocket := bindZmqSocket(config.ZmqConfig)
	signalsToHandle := []syscall.Signal{
		syscall.SIGSTOP, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	}
	graceFullShutdown(signalsToHandle, handlersList, zmqSocket, cli, ctx)
	managerEntryPointServer(config.WebSocketConfig, zmqSocket)
}
