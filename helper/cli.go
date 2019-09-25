package helper

import (
	"fmt"
	"os/exec"
	"sync"
)

func ExecCommand(cmd string, args []string, wg *sync.WaitGroup) {
	fmt.Println("Execution command: ", cmd)
	fmt.Println("Command argument: ", args)
	// splitting head => g++ parts => rest of the command
	//parts := strings.Fields(cmd)
	//head := parts[0]
	//fmt.Println(head)
	//parts = parts[1:len(parts)]
	//fmt.Println(parts)
	// exec command
	out, err := exec.Command(cmd, args...).Output()
	// show callback
	fmt.Printf("%s", out)
	// handle error
	CheckError(err, "Failed execution cli command")
	// Need to signal to wait group that this goroutine is done
	wg.Done()
}
