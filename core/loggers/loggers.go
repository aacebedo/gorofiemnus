package loggers

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// MainLogger The main logger for the library.
	MainLogger *logrus.Entry  //nolint:gochecknoglobals // Intends to set the loggers for the whole library
	log        *logrus.Logger //nolint:gochecknoglobals // Intends to set the loggers for the whole library
)

// InitLoggers Initialize the loggers for the whole library.
func InitLoggers() {
	log = logrus.New()
	log.Out = os.Stdout
	MainLogger = log.WithFields(logrus.Fields{"emitter": "main"})
}

// SetLogLevel Set the log level used in the library.
func SetLogLevel(lvl logrus.Level) {
	log.SetLevel(lvl)
}
