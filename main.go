package main

import (
	"github.com/aasumitro/mego-worker/service"
)

func main() {
	// do subscribe to messaging queue
	service.SubscribeMessage()
}
