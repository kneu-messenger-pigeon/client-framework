package framework

import (
	"fmt"
	"io"
)

type DebugLogger struct {
	out     io.Writer
	enabled bool
}

func (logger *DebugLogger) Log(format string, a ...any) {
	if !logger.enabled || logger.out == nil {
		return
	}
	if format[len(format)-1] != '\n' {
		format += "\n"
	}

	_, _ = fmt.Fprintf(logger.out, format, a...)
}
