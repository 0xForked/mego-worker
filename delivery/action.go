package delivery

import (
	"github.com/aasumitro/mego-worker/data"
	"github.com/aasumitro/mego-worker/helper"
	"log"
)

func SendToDevice(outbox data.Outbox) {
	log.Printf("send to device")
}

func SendToQueueTable(outbox data.Outbox) {
	text := len(outbox.TextDecoded)
	if text > 160 {
		log.Printf("Handle Multipart Message")
	}

	if text < 160 {
		callback := data.Store(outbox)
		if callback != nil {
			helper.CheckError(callback, callback.Error())
		}
	}

}
