package framework

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDebugLogger_Log(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := DebugLogger{
			out:     out,
			enabled: true,
		}

		logger.Log("test %s", "message")
		assert.Equal(t, "test message\n", out.String())
	})

	t.Run("disabled", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := DebugLogger{
			out:     out,
			enabled: false,
		}

		logger.Log("test %s", "message")
		assert.Empty(t, out.String())
	})
}
