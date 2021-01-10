package gorofimenus

import "github.com/rotisserie/eris"

var (
	// InvalidArgumentError Sentinel error for invalid argument.
	//nolint:gochecknoglobals // Intends to be a sentinel value
	InvalidArgumentError error = eris.New("Invalid argument")

	// AlreadyExistsError Sentinel error when a value already exists in a collection.
	//nolint:gochecknoglobals // Intends to be a sentinel value
	AlreadyExistsError error = eris.New("Element already exists")

	// NotExistsError Sentinel error when a value does not exist in a collection.
	//nolint:gochecknoglobals // Intends to be a sentinel value
	NotExistsError error = eris.New("Element does not exist")

	// InternalError Sentinel error when a invalid argument is given.
	//nolint:gochecknoglobals // Intends to be a sentinel value
	InternalError error = eris.New("Internal error")

	// ConfigurationLoadingError Sentinel error when an invalid configuration is loaded.
	//nolint:gochecknoglobals // Intends to be a sentinel value
	ConfigurationLoadingError error = eris.New("Unable to load the configuration")
)
