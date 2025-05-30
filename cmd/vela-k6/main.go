// SPDX-License-Identifier: Apache-2.0

// Package main is the entry point for the Vela K6 plugin.
// It captures the version information, configures the plugin from environment variables,
// runs the setup script, and executes performance tests.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-vela/vela-k6/plugin"
	"github.com/go-vela/vela-k6/version"
)

func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	var bytes []byte

	var err error

	if bytes, err = json.MarshalIndent(v, "", "  "); err != nil {
		log.Fatal(err)
	}

	// output the version information to stdout
	_, _ = fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	p := plugin.New()
	if err = p.ConfigFromEnv(); err != nil {
		log.Fatalf("FATAL: %s\n", err)
	}

	if err = p.RunSetupScript(); err != nil {
		log.Fatalf("FATAL: %s\n", err)
	}

	if err = p.RunPerfTests(); err != nil {
		log.Fatalf("FATAL: %s\n", err)
	}
}
