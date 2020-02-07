package main

import (
	"net/http"
	"sync"

	"github.com/moorara/konfig"
	"github.com/moorara/observe/log"
)

var config = struct {
	sync.Mutex
	LogLevel string
}{
	// default value
	LogLevel: "Info",
}

func main() {
	// Create the logger
	lo := log.Options{Name: "server"}
	logger := log.NewLogger(lo)

	// Listening for any update to configurations
	ch := make(chan konfig.Update)
	go func() {
		for update := range ch {
			if update.Name == "LogLevel" {
				config.Lock()
				lo.Level = config.LogLevel
				config.Unlock()
				logger.SetOptions(lo)
			}
		}
	}()

	// Watching for configurations
	close, _ := konfig.Watch(&config, []chan konfig.Update{
		ch,
	})

	defer close()

	// HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.InfoKV("message", "new request received")
		w.WriteHeader(http.StatusOK)
	})

	// Starting the HTTP server
	logger.InfoKV("message", "starting http server ...")
	http.ListenAndServe(":8080", nil)
}
