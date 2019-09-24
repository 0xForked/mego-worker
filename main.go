package main

import (
	"github.com/aasumitro/mego-worker/service"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// do subscribe to messaging queue
	service.SubscribeMessage()
}
