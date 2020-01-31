package main

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/moorara/konfig"
)

var config = struct {
	Port      uint16
	LogLevel  string
	Timeout   time.Duration
	Endpoints []url.URL
}{
	Port:     8080,            // default port
	LogLevel: "info",          // default log level
	Timeout:  2 * time.Minute, // default API call timeout
}

func main() {
	konfig.Pick(&config)
	flag.Parse()

	fmt.Printf("\nConfigurations: %+v\n\n", config)
}
