package entity

import (
	"io"

	"github.com/sirupsen/logrus"
)

func CloseBody(c io.Closer) {
	if err := c.Close(); err != nil {
		logrus.Errorf("closing response body: %v", err)
	}
}
