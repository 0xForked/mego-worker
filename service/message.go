package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aasumitro/mego-worker/data"
	"github.com/aasumitro/mego-worker/delivery"
	"github.com/aasumitro/mego-worker/helper"
	"log"
	"time"
)

// declare string map to handle queue message
type Message map[string]interface{}

func SubscribeMessage() {
	// load config
	config := helper.GetDataConfig()
	// make connection to rabbit mq
	mq, _ := data.MQConnection(config.RabbitMQ)
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
		"sms_notify", 	// name
		true,         	// durable
		false,         // delete when unused
		false,        	// exclusive
		false,        	// no-wait
		nil,          		// arguments
	)
	// if there is an error, handle it
	helper.CheckError(err, "Failed to declare a queue")
	// Qos controls how many messages or how many bytes the server will try to keep on
	// the network for consumers before receiving delivery acks.  The intent of Qos is
	// to make sure the network buffers stay full between the server and client.
	err = channel.Qos(
		1,     	// prefetch count
		0,     	// prefetch size
		false, 		// global
	)
	// if there is an error, handle it
	helper.CheckError(err, "Failed to set QoS")
	// Begin receiving on the returned chan Delivery before any other operation on the
	// Connection or Channel.
	msg, err := channel.Consume(
		queue.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
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
			log.Printf("Received a message: %s", d.Body)
			// doing action there is a new message
			validateAction(d.Body)


			// -----------
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			// show finish message
			log.Printf("Done")
			// Ack delegates an acknowledgement through the Acknowledger interface that the
			// client or server has finished work on a delivery.
			d.Ack(false)
		}
	}()
	// show waiting message
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// We need to block the main thread so that the above thread stays
	// on reading from msgs channel. To do that just try to read in from
	// the forever channel. As long as no one writes to it we will wait here.
	// Since we are the only ones that know of it it is guaranteed that
	// nothing gets written in it. We could also do a busy wait here but
	// that would waste CPU cycles for no good reason.
	<-forever
}

func validateAction(b []byte) {
	// convert/deserialize data from queue
	convert, err := deserialize(b)
	// if there is an error, handle it
	helper.CheckError(err, "Failed deserialize message")
	// convert to outbox model
	outbox := data.Outbox{
		DestinationNumber: fmt.Sprint(convert["phone"]),
		TextDecoded:       fmt.Sprint(convert["message"]),
		CreatorID:         "worker_one",
		Class:             "-1",
	}
	// handle important message
	if convert["status"] == "important" {
		delivery.SendToDevice(outbox)
	}
	// handle message as queue
	if convert["status"] == "queue" {
		delivery.SendToQueueTable(outbox)
	}
}

func deserialize(b []byte) (Message, error) {
	var msg Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}