package main

import (
	"sync"
	"time"

	"github.com/moorara/konfig"
)

// Global is the single source of truth for all configurations
var Global = struct {
	sync.Mutex
	LogLevel string
}{
	LogLevel: "info",
}

func startLogging(logger *Logger) {
	wait := make(chan struct{}, 4)

	go func() {
		t1 := time.NewTicker(500 * time.Millisecond)
		for range t1.C {
			logger.Debug("debugging this service!")
		}
	}()

	go func() {
		t2 := time.NewTicker(1 * time.Second)
		for range t2.C {
			logger.Info("Just info level logging!")
		}
	}()

	go func() {
		t3 := time.NewTicker(2 * time.Second)
		for range t3.C {
			logger.Warn("Warning!")
		}
	}()

	go func() {
		t4 := time.NewTicker(4 * time.Second)
		for range t4.C {
			logger.Error("Error happened!")
		}
	}()

	<-wait
}

func main() {
	logger := &Logger{}

	// Listening for configuration values and acting on them
	ch := make(chan konfig.Update, 1)
	go func() {
		for update := range ch {
			if update.Name == "LogLevel" {
				Global.Lock()
				logger.SetLevel(Global.LogLevel)
				Global.Unlock()
			}
		}
	}()

	// Start watching for configurations values
	stop, _ := konfig.Watch(&Global, []chan konfig.Update{ch}, konfig.WatchInterval(5*time.Second))
	defer stop()

	// Simulate logging
	startLogging(logger)
}
