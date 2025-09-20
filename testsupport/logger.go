package testsupport

import (
	"bytes"

	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Logger {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	return logger
}
