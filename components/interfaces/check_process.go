package interfaces

import (
	l "github.com/faelmori/logz"
)

// ICheckProcess defines the interface for managing and monitoring a process check.
type ICheckProcess interface {
	// WatchResults monitors and retrieves the results of the process.
	// Returns:
	//   IResult: The result of the process.
	WatchResults() IResult

	// WatchErrors monitors and retrieves any errors that occur during the process.
	// Returns:
	//   error: An error encountered during the process, if any.
	WatchErrors() error

	// WatchDone monitors and checks if the process is completed.
	// Returns:
	//   bool: True if the process is done, false otherwise.
	WatchDone() bool

	// GetWorker retrieves the worker associated with the process.
	// Returns:
	//   IWorker: The worker instance.
	GetWorker() IWorker

	// GetPackages retrieves the list of packages associated with the process.
	// Returns:
	//   []string: A slice of package names.
	GetPackages() []string

	// GetConfig retrieves the configuration associated with the process.
	// Returns:
	//   IConfig: The configuration instance.
	GetConfig() any //IConfig

	// GetChanResult retrieves the channel used to communicate process results.
	// Returns:
	//   chan IResult: The result channel.
	GetChanResult() chan IResult

	// GetChanError retrieves the channel used to communicate errors.
	// Returns:
	//   chan error: The error channel.
	GetChanError() chan error

	// GetChanDone retrieves the channel used to signal process completion.
	// Returns:
	//   chan bool: The done channel.
	GetChanDone() chan bool

	// GetLogger retrieves the logger instance used by the process.
	// Returns:
	//   l.Logger: The logger instance.
	GetLogger() l.Logger

	// SetWorker sets the worker associated with the process.
	// Parameters:
	//   IWorker: The worker instance to set.
	SetWorker(IWorker)

	// SetPackages sets the list of packages associated with the process.
	// Parameters:
	//   []string: A slice of package names to set.
	SetPackages([]string)

	// SetConfig sets the configuration associated with the process.
	// Parameters:
	//   IConfig: The configuration instance to set.
	SetConfig(any)

	// SetChanResult sets the channel used to communicate process results.
	// Parameters:
	//   chan IResult: The result channel to set.
	SetChanResult(chan IResult)

	// SetChanError sets the channel used to communicate errors.
	// Parameters:
	//   chan error: The error channel to set.
	SetChanError(chan error)

	// SetChanDone sets the channel used to signal process completion.
	// Parameters:
	//   chan bool: The done channel to set.
	SetChanDone(chan bool)

	// SetLogger sets the logger instance used by the process.
	// Parameters:
	//   l.Logger: The logger instance to set.
	SetLogger(l.Logger)
}
