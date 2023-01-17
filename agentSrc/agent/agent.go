////STRING TO JSON https://stackoverflow.com/questions/28859941/whats-the-golang-equivalent-of-converting-any-json-to-standard-dict-in-python

package main

import (
	"agent/modules"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net"
	"net/url"
	"os"
)

func getWebSocketConn(config modules.WebSocketConfig) *websocket.Conn {
	var addr = flag.String("addr", config.Host+":"+config.Port, "http service address")
	urlObj := url.URL{Scheme: config.Scheme, Host: *addr, Path: config.Path}
	log.Printf("Connecting to %s", urlObj.String())
	conn, _, err := websocket.DefaultDialer.Dial(urlObj.String(), nil)
	modules.CheckErr(err)
	return conn
}

func getSocketConn(config modules.SocketConfig) net.Listener {
	listener, err := net.Listen(config.Network, config.Host+":"+config.Port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on " + config.Host + ":" + config.Port)
	fmt.Println("Waiting for client...")
	return listener
}

func processSocketResponse(sConn net.Conn, wsConn *websocket.Conn) {
	buffer := make([]byte, 512)
	for {
		mLen, err := sConn.Read(buffer)
		if err != nil && err == io.EOF {
			_ = sConn.Close()
			return
		}
		err = wsConn.WriteMessage(websocket.TextMessage, buffer[:mLen])
		modules.CheckErr(err)
	}

}

func main() {
	filePath := os.Args[1]
	config := modules.ReadConfig(filePath)
	wsConn := getWebSocketConn(config.WebSocketConfig)
	listener := getSocketConn(config.SocketConfig)
	for {
		sConn, err := listener.Accept()
		modules.CheckErr(err)
		fmt.Println("client connected")
		modules.CheckErr(err)
		go processSocketResponse(sConn, wsConn)
	}
}
