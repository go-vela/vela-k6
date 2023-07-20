package main

import (
	"log"

	"github.com/go-vela/vela-k6/plugin"
)

func main() {
	cfg, err := plugin.ConfigFromEnv()
	if err != nil {
		log.Fatalf("FATAL: %s\n", err)
	}

	err = plugin.RunPerfTests(cfg)
	if err != nil {
		log.Fatalf("FATAL: %s\n", err)
	}
}
