package helper

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func ExecCommand(cmd string, wg *sync.WaitGroup) {
	fmt.Println("Execution command: ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]
	// exec command
	out, err := exec.Command(head, parts...).Output()
	// handle error
	CheckError(err, "Failed execution cli command")
	// show callback
	fmt.Printf("%s", out)
	// Need to signal to wait group that this goroutine is done
	wg.Done()
}
