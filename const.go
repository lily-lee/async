package async

import (
	"time"

	"errors"
)

const defaultTimeout = 5 * time.Minute

var TimeoutErr = errors.New("timeout")
