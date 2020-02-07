package main

import (
	"sync"
	"time"

	"github.com/moorara/konfig"
)

// config is the single source of truth for all configurations
var config = struct {
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
			logger.Debug("Debugging ...")
		}
	}()

	go func() {
		t2 := time.NewTicker(1 * time.Second)
		for range t2.C {
			logger.Info("Informing ...")
		}
	}()

	go func() {
		t3 := time.NewTicker(2 * time.Second)
		for range t3.C {
			logger.Warn("Warning ...")
		}
	}()

	go func() {
		t4 := time.NewTicker(4 * time.Second)
		for range t4.C {
			logger.Error("Erroring ...")
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
				config.Lock()
				logger.SetLevel(config.LogLevel)
				config.Unlock()
			}
		}
	}()

	// Start watching for configurations values
	close, _ := konfig.Watch(&config, []chan konfig.Update{ch})
	defer close()

	// Simulate logging
	startLogging(logger)
}
