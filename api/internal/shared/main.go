package shared

import "time"

type LocalNowTimeProvider = func() time.Time
type UTCNowTimeProvider = func() time.Time
