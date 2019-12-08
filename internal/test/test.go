// Package test is a helper for automatic tests
package test

import (
	log "github.com/sirupsen/logrus"
)

//const defaultLogLevel = log.WarnLevel
const defaultLogLevel = log.DebugLevel

// GetLogLevel returns the default log level for tests
func GetLogLevel() string {
	return defaultLogLevel.String()
}
