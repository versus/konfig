package konfig

import (
	"strings"
)

// settings struct has all the settings regarding how configuration values are read
type settings struct {
	checkForTelepresence bool
}

func (s settings) String() string {
	opts := []string{}

	if s.checkForTelepresence {
		opts = append(opts, "Telepresence")
	}

	return strings.Join(opts, " + ")
}

// Option configures how configuration values are read
type Option interface {
	apply(*settings)
}

// funcOption implements Option interface
type funcOption struct {
	fn func(*settings)
}

func (fo *funcOption) apply(s *settings) {
	fo.fn(s)
}

func newFuncOption(fn func(*settings)) *funcOption {
	return &funcOption{
		fn: fn,
	}
}

// Telepresence is the option for reading files when running in a Telepresence shell.
// If the TELEPRESENCE_ROOT environment variable exist, files will be read from mounted volume.
// See https://telepresence.io/howto/volumes.html for details.
func Telepresence() Option {
	return newFuncOption(func(s *settings) {
		s.checkForTelepresence = true
	})
}
