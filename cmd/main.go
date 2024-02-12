package main

import (
	"github.com/cukhoaimon/SimpleBank/internal/app"
	"github.com/cukhoaimon/SimpleBank/utils"
	"log"
)

func main() {
	// load config
	config, err := utils.LoadConfig(".")
	if err != nil {

		log.Fatal(err.Error())
	}

	app.Run(config)
}
