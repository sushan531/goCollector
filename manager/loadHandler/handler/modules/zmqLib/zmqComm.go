package zmqLib

import (
	"github.com/pebbe/zmq4"
	"loadHandler/modules"
	"log"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetSocketListen(pull modules.ZmqPull) *zmq4.Socket {
	server, _ := zmq4.NewSocket(zmq4.PULL)
	err := server.Connect("tcp://127.0.0.1:5556")
	CheckErr(err)
	return server
}

func GetMessage(socket *zmq4.Socket) []string {
	message, _ := socket.RecvMessage(0)
	return message
}
