// Package logger provides central logging
package logger

import (
	"github.com/pgillich/errfmt"
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger // nolint:gochecknoglobals

// Get returns a global logger
func Get() *log.Logger {
	return logger
}

// Init sets the global logger
func Init(logLevel string) {
	errfmt.AddSkipPackageFromStackTrace("github.com/pgillich/sample-blog")

	if logger == nil {
		flags := errfmt.FlagCallStackOnConsole | errfmt.FlagCallStackInHTTPProblem | errfmt.FlagPrintStructFieldNames
		callStackSkipLast := 3

		logLevelValue, err := log.ParseLevel(logLevel)
		if err != nil {
			logger.Panic("invalid log level, ", err)
		}

		logger = errfmt.NewTextLogger(logLevelValue, flags, callStackSkipLast)
	}
}
