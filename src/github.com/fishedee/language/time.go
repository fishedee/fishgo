package language

import (
	"time"
)

var (
	ZERO_TIME time.Time
)

func init() {
	ZERO_TIME = time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
}
