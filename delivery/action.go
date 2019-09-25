package delivery

import (
	"fmt"
	"github.com/aasumitro/mego-worker/data"
	"github.com/aasumitro/mego-worker/helper"
	"log"
	"sync"
)

// make sure your gammu is running well
// and your modem/device are connected
func SendToDevice(outbox data.Outbox) {
	log.Printf("send to device")
	// A WaitGroup waits for a collection of goroutines to finish.
	wg := new(sync.WaitGroup)
	// Add adds delta, which may be negative, to the WaitGroup counter.
	// If the counter becomes zero, all goroutines blocked on Wait are released.
	// If the counter goes negative, Add panics.
	wg.Add(1)
	// "echo '%s' | gammu --sendsms TEXT %s"
	// Command String
	cmd := fmt.Sprintf(
		"echo send text '%s' to %s",
		outbox.TextDecoded,
		outbox.DestinationNumber)
	// exec command with goroutine
	go helper.ExecCommand(cmd, wg)
	// Wait blocks until the WaitGroup counter is zero.
	wg.Wait()

	// TODO make function to store sentitems data
}

func SendToQueueTable(outbox data.Outbox) {
	// A WaitGroup waits for a collection of goroutines to finish.
	wg := new(sync.WaitGroup)
	// Add adds delta, which may be negative, to the WaitGroup counter.
	// If the counter becomes zero, all goroutines blocked on Wait are released.
	// If the counter goes negative, Add panics.
	wg.Add(1)
	// count text length
	text := len(outbox.TextDecoded)
	// handle when message with text under 160 chars
	if text < 160 {
		log.Printf("inject to device")
		// "echo '%s' | gammu-smsd-inject TEXT %s"
		// Command String
		cmd := fmt.Sprintf(
			"echo send text '%s' to %s",
			outbox.TextDecoded,
			outbox.DestinationNumber)
		// exec command with goroutine
		go helper.ExecCommand(cmd, wg)
	}
	// handle when message have more than 160 chars
	if text > 160 {
		log.Printf("Handle Multipart Message")
		// echo '%s' | gammu-smsd-inject TEXT %s -len 400
		cmd := fmt.Sprintf(
			"echo send text '%s' to %s",
			outbox.TextDecoded,
			outbox.DestinationNumber)
		// exec command with goroutine
		go helper.ExecCommand(cmd, wg)
	}
	// exec command with goroutine
	wg.Wait()
}
