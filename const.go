package async

import (
	"time"

	"errors"
)

const (
	defaultTimeout = 5 * time.Minute
	defaultPoolNum = 100
)

var TimeoutErr = errors.New("timeout")
