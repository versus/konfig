package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/moorara/konfig"
	"github.com/moorara/observe/log"
)

var config = struct {
	sync.Mutex
	LogLevel      string
	ServerAddress string
}{
	// default values
	LogLevel:      "info",
	ServerAddress: "http://localhost:8080",
}

func main() {
	// Create the logger
	lo := log.Options{Name: "client"}
	logger := log.NewLogger(lo)

	// Server address
	endpoint := "/"
	url := fmt.Sprintf("%s%s", config.ServerAddress, endpoint)

	// Listening for any update to configurations
	ch := make(chan konfig.Update)
	go func() {
		for update := range ch {
			switch update.Name {
			case "LogLevel":
				config.Lock()
				lo.Level = config.LogLevel
				config.Unlock()
				logger.SetOptions(lo)
			case "ServerAddress":
				config.Lock()
				url = fmt.Sprintf("%s%s", config.ServerAddress, endpoint)
				config.Unlock()
			}
		}
	}()

	// Watching for configurations
	close, _ := konfig.Watch(&config, []chan konfig.Update{
		ch,
	})

	defer close()

	// Sending requests to server
	logger.InfoKV("message", "start sending requests ...")

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: &http.Transport{},
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.ErrorKV("error", err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			logger.ErrorKV("error", err)
			continue
		}

		logger.InfoKV("message", "response received from server", "http.statusCode", resp.StatusCode)
	}
}
