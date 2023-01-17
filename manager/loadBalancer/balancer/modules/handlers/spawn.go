package handlers

import (
	"fmt"
	"log"
	"os"
)

var LoadHandlers []*os.Process

func SpawnLoadHandlers(minHandlers int) {
	for i := 0; i < minHandlers; i++ {
		var procAttr os.ProcAttr
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr} // todo make it write to file
		// todo make it dynamic !!!
		process, err := os.StartProcess(
			"/Users/spt/GolandProjects/collector/manager/loadHandler/handler/loadHandler",
			[]string{"", "/Users/spt/GolandProjects/collector/manager/loadHandler/"},
			&procAttr,
		)
		if err != nil {
			log.Fatal(err)
		}
		LoadHandlers = append(LoadHandlers, process)
	}
	for _, value := range LoadHandlers {
		fmt.Println("Process Id : ", value)
	}

}
