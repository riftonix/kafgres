package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Setup configures the global logger.
func Setup() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}
