package utils

import (
	"flag"
	"log"
	"os"
)

const (
	FlagConfigName = "config"
	EnvConfigName  = "CONFIG_PATH"
)

func ConfigPath() (path string) {
	flag.StringVar(&path, FlagConfigName, "", "config path")
	flag.Parse()

	if path == "" {
		path = os.Getenv(EnvConfigName)
	}

	if path == "" {
		log.Fatal("config path is required")
	}

	return
}
