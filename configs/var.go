package configs

import (
	"time"
)

// TimeNowFunc returns time.Now, hacked by automatic tests
var TimeNowFunc = func() func() time.Time { //nolint:gochecknoglobals
	return time.Now
}
