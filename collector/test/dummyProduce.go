package test

import (
	"collector/modules"
	"fmt"
	"github.com/pebbe/zmq4"
	"strconv"
	"time"
)

// This is a dummy message generator code
func createJsonMsg(count int) string {
	msg := map[string]string{
		"msg":   "This is message" + strconv.Itoa(count),
		"count": strconv.Itoa(count),
	}
	return fmt.Sprintf("%v", msg)
}
func SendMsg() {
	client, _ := zmq4.NewSocket(zmq4.PUSH)
	err := client.Connect("tcp://127.0.0.1:5555")
	modules.CheckErr(err)
	var a = 0
	for {
		_, _ = client.SendMessage(createJsonMsg(a))
		a += 1
		time.Sleep(time.Second)
	}
}
