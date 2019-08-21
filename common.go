package async

import (
	"log"
	"runtime/debug"
	"time"
)

func handlePanic() {
	if err := recover(); err != nil {
		log.Println("panic: ", err)
		debug.PrintStack()
	}
}

func handle(timeout time.Duration, f func() error) error {
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		eC := make(chan error)
		go func() {
			eC <- f()
		}()
		select {
		case <-timer.C:
			return TimeoutErr
		case err := <-eC:
			return err
		}
	}

	return f()
}
