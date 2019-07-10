package konfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingsString(t *testing.T) {
	tests := []struct {
		name           string
		settings       settings
		expectedString string
	}{
		{
			"NoOption",
			settings{},
			"",
		},
		{
			"WithTelepresence",
			settings{
				checkForTelepresence: true,
			},
			"Telepresence",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.settings.String())
		})
	}
}

func TestTelepresence(t *testing.T) {
	tests := []struct {
		name             string
		settings         *settings
		expectedSettings *settings
	}{
		{
			"OK",
			&settings{},
			&settings{
				checkForTelepresence: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opt := Telepresence()
			opt.apply(tc.settings)

			assert.Equal(t, tc.expectedSettings, tc.settings)
		})
	}
}
