package entity

import (
	"io"
	"log/slog"
)

func CloseBody(c io.Closer) {
	if err := c.Close(); err != nil {
		slog.Error("closing response body", err)
	}
}
