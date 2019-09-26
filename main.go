package main

import (
	"github.com/aasumitro/mego-worker/helper"
	"github.com/aasumitro/mego-worker/service"
)

func main() {
	config := helper.GetConfig()
	app := service.App{}
	app.SubscribeMessage(config)
}
