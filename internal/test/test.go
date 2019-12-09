// Package test is a helper for automatic tests
package test

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//const defaultLogLevel = log.DebugLevel
const defaultLogLevel = log.WarnLevel

// GetLogLevel returns the default log level for tests
func GetLogLevel() string {
	return defaultLogLevel.String()
}

// JSONMarshal marshals with indent and returns string (error is sinked)
func JSONMarshal(v interface{}) string {
	bytes, _ := json.MarshalIndent(&v, "", "  ") //nolint:errcheck

	return string(bytes)
}

// GetHTTPHeaderJSON returns "Content-Type: application/json"
func GetHTTPHeaderJSON() http.Header {
	return http.Header{
		"Content-Type": {"application/json"},
	}
}

// GetHTTPHeaderJSONToken returns:
// "Content-Type: application/json"
// "Authorization:Bearer $TOKEN"
func GetHTTPHeaderJSONToken(token string) http.Header {
	return http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token},
	}
}
