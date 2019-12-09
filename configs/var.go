package configs

import (
	"time"
)

var (
	// TimeNowFunc returns time.Now, hacked by automatic tests
	TimeNowFunc = func() func() time.Time { //nolint:gochecknoglobals
		return time.Now
	}

	// BuildTag is set at build time
	BuildTag string //nolint:gochecknoglobals
	// BuildCommit is set at build time
	BuildCommit string //nolint:gochecknoglobals
	// BuildBranch is set at build time
	BuildBranch string //nolint:gochecknoglobals
	// BuildTime is set at build time
	BuildTime string //nolint:gochecknoglobals
)
