package chanlock

import (
	"fmt"
	"sync"
	"time"
)

type DemoLockWithChan struct {
	semaphore chan struct{}
}

func NewLockWithChan(cap int) *DemoLockWithChan {
	return &DemoLockWithChan{
		semaphore: make(chan struct{}, cap),
	}
}

func (dlwc *DemoLockWithChan) ExecuteWithLock(do func() error) error {
	timer := time.After(60 * time.Second)
	var err error
	select {
	case <-timer:
		err = fmt.Errorf("timeout on acquire")
	case dlwc.semaphore <- struct{}{}:
		defer dlwc.Release()
		err = do()
	}
	return err
}

func (dlwc *DemoLockWithChan) Release() {
	select {
	case <-dlwc.semaphore:
	default:
	}
}

func Demo() {
	cap := 5
	lockDemo := NewLockWithChan(cap)
	howManyTask := 5
	results := make(chan int, howManyTask)
	errs := make(chan error, howManyTask)
	makeDo := func(sec time.Duration) func() error {
		res := func() error {
			time.Sleep(sec * time.Second)
			// fmt.Printf("do for %d seconds\n", int(sec))
			// fmt.Printf("val %d \n", int(sec) % 2)
			random := (int(sec) % 2) == 0
			if random {
				return fmt.Errorf("error when do: " + fmt.Sprintf("%d", int(sec)))
			}
			return nil
		}
		return res
	}
	var wg sync.WaitGroup
	for i := 0; i < howManyTask; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sec := time.Duration(idx)
			err := lockDemo.ExecuteWithLock(makeDo(sec))
			if err != nil {
				errs <- err
			} else {
				results <- idx
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(errs)
		close(results)
	}()

	var resultDone, errDone bool
	for !resultDone || !errDone {
		select {
		case r, ok := <-results:
			if ok {
				fmt.Printf("result: %d\n", r)
			} else {
				resultDone = true
			}
		case e, ok := <-errs:
			if ok {
				fmt.Printf("error: %+v\n", e)
			} else {
				errDone = true
			}
		}
	}
}
