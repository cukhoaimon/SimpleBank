package main

import (
	"github.com/cukhoaimon/SimpleBank/internal/app"
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// load config
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Err(err)
	}

	app.Run(config)
}
