package main

import (
	"github.com/CalebEWheeler/StateFlow/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	App := app.New()
	if err := App.Run(); err != nil {
		log.Error(err)
	}
}
