/* Package configs contains global configs
It should not import any local packages.
*/
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

	// OptDbDialect is the DB dialect (Gorm driver name)
	OptDbDialect = "db-dialect"
	// DefaultDbDialect is default value to OptDbDialect
	DefaultDbDialect = "sqlite3"

	// OptDbDsn is the DB connection info
	OptDbDsn = "db-dsn"
	// DefaultDbDsn is default value to OptDbDsn
	DefaultDbDsn = ":memory:"

	// OptDbSample enables filling DB by sampe data
	OptDbSample = "db-sample"
	// DefaultDbSample is default value to OptDbSample
	DefaultDbSample = true
)
