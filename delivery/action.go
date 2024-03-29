package delivery

import (
	"github.com/aasumitro/mego-worker/data"
	"github.com/aasumitro/mego-worker/helper"
	"log"
	"sync"
)

// run this command if gammu-smsd is not activate
// make sure your gammu is running well
// and your modem/device are connected
func SendToDevice(outbox data.Outbox) {
	// A WaitGroup waits for a collection of goroutines to finish.
	wg := new(sync.WaitGroup)
	// Add adds delta, which may be negative, to the WaitGroup counter.
	// If the counter becomes zero, all goroutines blocked on Wait are released.
	// If the counter goes negative, Add panics.
	wg.Add(1)
	// gammu command
	cmd := "gammu"
	// command argument
	args := []string{"sendsms", "TEXT", outbox.DestinationNumber, "-text", outbox.TextDecoded}
	// exec command with goroutine
	go helper.ExecCommand(cmd, args, wg)
	// Wait blocks until the WaitGroup counter is zero.
	wg.Wait()
	// TODO make function to store sentitems data
	//data.CreateSentMessage()
}

// this function is for send message using gammu-smsd-inject to
// gammu database schema as queue table
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
		// gammu-smsd-inject command
		cmd := "gammu-smsd-inject"
		// gammu-smsd-inject command argument
		args := []string{"TEXT", outbox.DestinationNumber, "-text", outbox.TextDecoded}
		// exec command with goroutine
		go helper.ExecCommand(cmd, args, wg)
	}
	// handle when message have more than 160 chars
	if text > 160 {
		// gammu-smsd-inject command
		cmd := "gammu-smsd-inject"
		// gammu-smsd-inject command argument
		args := []string{"TEXT", "-len", "400", outbox.DestinationNumber, "-text", outbox.TextDecoded}
		// exec command with goroutine
		go helper.ExecCommand(cmd, args, wg)
	}
	// exec command with goroutine
	wg.Wait()
}

// this function is inject direct to specified table in database schema
func StoreOutbox(outbox data.Outbox) {
	log.Printf("Add new record to outbox queue")
	// store queue from MQ to queue DB
	err := data.CreateOutboxMessage(outbox)
	// check if query is have an error
	helper.CheckError(err, "Failed Store Message")
}
