package eslgo

import (
	"io"
	"log/slog"
)

// NewDiscardLogger returns a logger that discards all log messages.
// This is useful when you want to suppress all logging output.
func NewDiscardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
