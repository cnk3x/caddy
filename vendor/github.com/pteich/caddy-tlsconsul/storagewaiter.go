package storageconsul

import (
	"sync"
	"time"
)

// Implementation of certmagic.Waiter
type consulStorageWaiter struct {
	key          string
	waitDuration time.Duration
	wg           *sync.WaitGroup
}

func (csw *consulStorageWaiter) Wait() {
	csw.wg.Add(1)
	go time.AfterFunc(csw.waitDuration, func() {
		csw.wg.Done()
	})
	csw.wg.Wait()
}
