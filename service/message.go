package service

import (
	"bytes"
	"fmt"
	"github.com/aasumitro/mego-worker/data"
	"github.com/aasumitro/mego-worker/delivery"
	"github.com/aasumitro/mego-worker/helper"
	"log"
	"time"
)

const (
	IMPORTANT = "important"
	INJECT    = "inject"
	SCHEMA    = "schema"
)

type App struct {
	Config *helper.Config
}

func (a *App) SubscribeMessage(config *helper.Config) {
	// load config
	a.Config = config
	// make connection to rabbit mq
	mq, _ := data.MQConnection(a.Config.RabbitMQ)
	// defer the close till after the main function has finished
	// executing
	defer mq.Close()
	// create messaging queue chanel
	channel, err := mq.Channel()
	// if there is an error with chanel, handle it
	helper.CheckError(err, "Failed to open a channel")
	// defer the close till after the main function has finished
	defer channel.Close()
	// QueueDeclare declares a queue to hold messages and deliver to consumers.
	// Declaring creates a queue if it doesn't already exist, or ensures that an
	// existing queue matches the same parameters.
	queue, err := channel.QueueDeclare(
		"sms_notify", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	// if there is an error, handle it
	helper.CheckError(err, "Failed to declare a queue")
	// Qos controls how many messages or how many bytes the server will try to keep on
	// the network for consumers before receiving delivery acks.  The intent of Qos is
	// to make sure the network buffers stay full between the server and client.
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	// if there is an error, handle it
	helper.CheckError(err, "Failed to set QoS")
	// Begin receiving on the returned chan Delivery before any other operation on the
	// Connection or Channel.
	msg, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	// if there is an error, handle it
	helper.CheckError(err, "Failed to register a consumer")
	// create an unbuffered channel for bool types.
	// Type is not important but we have to give one anyway.
	forever := make(chan bool)
	// fire up a goroutine that hooks onto message channel and reads
	// anything that pops into it. This essentially is a thread of
	// execution within the main thread. message is a channel constructed by
	// previous code.
	go func() {
		for d := range msg {
			// show log if new message is received
			helper.ShowMessage(fmt.Sprintf("Received a message: %s", d.Body))
			// make it happen
			validateAction(d.Body, a.Config.Service.Delivery)

			// -----------
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			// show finish message
			helper.ShowMessage("Done")
			// Ack delegates an acknowledgement through the Acknowledger interface that the
			// client or server has finished work on a delivery.
			d.Ack(false)
		}
	}()
	// show waiting message
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// We need to block the main thread so that the above thread stays
	// on reading from msg channel. To do that just try to read in from
	// the forever channel. As long as no one writes to it we will wait here.
	// Since we are the only ones that know of it it is guaranteed that
	// nothing gets written in it. We could also do a busy wait here but
	// that would waste CPU cycles for no good reason.
	<-forever
}

func validateAction(d []byte, deliveryMode string) {
	// convert/deserialize data from queue
	msg, err := helper.Deserialize(d)
	// if there is an error, handle it
	helper.CheckError(err, "Failed deserialize message")
	// convert to outbox model
	outbox := data.Outbox{
		DestinationNumber: fmt.Sprint(msg["phone"]),
		TextDecoded:       fmt.Sprint(msg["message"]),
	}
	// validate how to delivery the data
	switch deliveryMode {
	case INJECT:
		// this command is can be use if gammu-smsd is not activated
		// handle important message
		if msg["status"] == IMPORTANT {
			// sent message to device
			delivery.SendToDevice(outbox)
		} else {
			// handle message with gammu-smsd
			// handle message as queue
			delivery.SendToQueueTable(outbox)
		}
	case SCHEMA:
		// store to queue table
		delivery.StoreOutbox(outbox)
	default:
		// store to queue table
		delivery.StoreOutbox(outbox)
	}

}
