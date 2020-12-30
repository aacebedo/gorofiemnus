package loggers

import (
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger = logrus.New() //nolint:gochecknoglobals // Intends to set the loggers for the whole library
	// MainLogger The main logger for the library.
	MainLogger *logrus.Entry = log.WithFields(logrus.Fields{"emitter": "main"}) //nolint:gochecknoglobals // Intends to set the loggers for the whole library
)

// SetLogLevel Set the log level used in the library.
func SetLogLevel(lvl logrus.Level) {
	log.SetLevel(lvl)
}
