// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
