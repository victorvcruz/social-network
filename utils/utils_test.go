package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStringEnvOrElse(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		os.Setenv("TEST", "value")
		value := GetStringEnvOrElse("TEST", "not used")
		assert.Equal(t, "value", value)
	})
	t.Run("not found", func(t *testing.T) {
		os.Unsetenv("TEST")
		value := GetStringEnvOrElse("TEST", "value")
		assert.Equal(t, "value", value)
	})
}

func TestGetIntEnvOrElse(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		os.Setenv("TEST", "1")
		value := GetIntEnvOrElse("TEST", 1)
		assert.Equal(t, 1, value)
	})
	t.Run("not found", func(t *testing.T) {
		os.Unsetenv("TEST")
		value := GetIntEnvOrElse("TEST", 1)
		assert.Equal(t, 1, value)
	})
}
