package healthcheck

import (
	"sync"
)

func WaitForCheck(f func() bool, wg *sync.WaitGroup) {
	wg.Done()
	f()
}
