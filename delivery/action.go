package delivery

import (
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
	// "echo '%s' | gammu --sendsms TEXT %s" - use this command or
	// gammu sendsms TEXT 082271115593 -text "Halooo .." - this command
	// Command String gammu sendsms TEXT %s -text '%s'
	//cmd := fmt.Sprintf(
	//	`gammu-smsd-inject TEXT %s -text "%s"`,
	//	outbox.DestinationNumber,
	//	outbox.TextDecoded)
	cmd := "gammu-smsd-inject"
	args := []string{"TEXT", outbox.DestinationNumber, "-text", outbox.TextDecoded}
	// exec command with goroutine
	go helper.ExecCommand(cmd, args, wg)
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
		log.Printf("send message with gammu-smsd-inject")
		// gammu-smsd-inject TEXT %s -text '%s'
		// Command String
		//cmd := fmt.Sprintf(
		//	"gammu-smsd-inject TEXT %s -text '%s'",
		//	outbox.DestinationNumber,
		//	outbox.TextDecoded)
		cmd := "gammu-smsd-inject"
		args := []string{"TEXT", outbox.DestinationNumber, "-text", outbox.TextDecoded}
		// exec command with goroutine
		go helper.ExecCommand(cmd, args, wg)
	}
	// handle when message have more than 160 chars
	if text > 160 {
		log.Printf("Handle Multipart Message")
		// echo '%s' | gammu-smsd-inject TEXT %s -len 400
		// gammu-smsd-inject TEXT %s -len 400 -text '%s'
		//cmd := fmt.Sprintf(
		//	"gammu-smsd-inject TEXT %s -len 400 -text '%s'",
		//	outbox.DestinationNumber,
		//	outbox.TextDecoded)
		cmd := "gammu-smsd-inject"
		args := []string{"TEXT", "-len", "400", outbox.DestinationNumber, "-text", outbox.TextDecoded}
		// exec command with goroutine
		go helper.ExecCommand(cmd, args, wg)
	}
	// exec command with goroutine
	wg.Wait()
}
