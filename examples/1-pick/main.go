package main

import (
	"flag"
	"fmt"

	"github.com/moorara/konfig"
)

// Global is the single source of truth for all configurations
var Global = struct {
	Enabled   bool
	Port      int
	LogLevel  string
	Endpoints []string
}{
	// Default Values
	Enabled:   true,
	Port:      8080,
	LogLevel:  "Warn",
	Endpoints: []string{"localhost:8080"},
}

func main() {
	konfig.Pick(&Global)
	flag.Parse()

	fmt.Printf("\nConfigurations: %+v\n\n", Global)
}
