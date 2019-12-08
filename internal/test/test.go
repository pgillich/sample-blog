// Package test is a helper for automatic tests
package test

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

//const defaultLogLevel = log.WarnLevel
const defaultLogLevel = log.DebugLevel

// GetLogLevel returns the default log level for tests
func GetLogLevel() string {
	return defaultLogLevel.String()
}

// JSONMarshal marshals with indent and returns string (error is sinked)
func JSONMarshal(v interface{}) string {
	bytes, _ := json.MarshalIndent(&v, "", "  ") //nolint:errcheck

	return string(bytes)
}
