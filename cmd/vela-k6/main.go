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
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

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
