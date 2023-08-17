package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-vela/vela-k6/plugin"
	"github.com/go-vela/vela-k6/version"
	"github.com/sirupsen/logrus"
)

func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	cfg, err := plugin.ConfigFromEnv()
	if err != nil {
		log.Fatalf("FATAL: %s\n", err)
	}

	err = plugin.RunSetupScript(cfg)
	if err != nil {
		log.Fatalf("FATAL: %s\n", err)
	}

	err = plugin.RunPerfTests(cfg)
	if err != nil {
		log.Fatalf("FATAL: %s\n", err)
	}
}
