package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogging(level string) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Level(lvl)
}
