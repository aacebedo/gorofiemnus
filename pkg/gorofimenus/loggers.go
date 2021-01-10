package gorofimenus

import (
	"github.com/sirupsen/logrus"
)

var (
	//nolint:gochecknoglobals // Intends to set the loggers for the whole library
	log *logrus.Logger = logrus.New()

	// MainLogger The main logger for the library.
	//nolint:gochecknoglobals // Intends to set the loggers for the whole library
	MainLogger *logrus.Entry = log.WithFields(logrus.Fields{"emitter": "gorofimenus"})

	// VerboseLogger The logger used for verbose stuff.
	//nolint:gochecknoglobals // Intends to set the loggers for the whole library
	VerboseLogger *logrus.Entry = log.WithFields(logrus.Fields{"emitter": "gorofimenus"})
)

// SetLogLevel Set the log level used in the library.
func SetLogLevel(lvl logrus.Level, verbosity uint32) {
	log.SetLevel(lvl)
}
