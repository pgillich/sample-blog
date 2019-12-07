// Package configs contains global configs
package configs

const (
	// OptLogLevel is the log level
	OptLogLevel = "log-level"
	// DefaultLogLevel is the default to OptLogLevel
	DefaultLogLevel = "DEBUG"

	// OptServiceHostPort is the host:port listening on
	OptServiceHostPort = "listen"
	// DefaultServiceHostPort is default value to OptServiceHostPort
	DefaultServiceHostPort = ":8088"

	// OptDbDsn is the DB connection info
	OptDbDsn = "db-dsn"
	// DefaultDbDsn is default value to OptDbDsn
	DefaultDbDsn = ":memory:"
)
