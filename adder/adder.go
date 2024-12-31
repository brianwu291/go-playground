package adder

import (
	"fmt"
	"sync"
)

type (
	AdderConfig struct {
		Size int
	}

	addEvent struct {
		Value int
	}

	valueControl struct {
		value int
		mu    sync.RWMutex
	}
)

type Adder struct {
	Config       AdderConfig
	events       chan *addEvent
	valueControl valueControl
	closeOnce    sync.Once
	done         chan struct{}
}

const (
	maxSize     = 30
	defaultSize = 10
)

func NewAdder(config *AdderConfig) *Adder {
	if config.Size <= 0 {
		config.Size = defaultSize
	}
	if config.Size > maxSize {
		config.Size = maxSize
	}

	addedEvents := make(chan *addEvent, config.Size)

	adder := &Adder{
		Config:       *config,
		events:       addedEvents,
		valueControl: valueControl{},
		done:         make(chan struct{}),
	}

	adder.processEvents()

	return adder
}

func (ad *Adder) Inc(value int) {
	if value <= 0 {
		fmt.Printf("value should bigger than zero. current is %d", value)
		return
	}

	select {
	case <-ad.done:
		fmt.Printf("adder has done\n")
		return
	default:
		ad.events <- &addEvent{Value: value}
		fmt.Printf("sending events...\n")
	}
}

func (ad *Adder) GetCurrentValue() int {
	ad.valueControl.mu.RLock()
	defer ad.valueControl.mu.RUnlock()
	cur := ad.valueControl.value
	return cur
}

func (ad *Adder) Close() {
	ad.closeOnce.Do(func() {
		close(ad.done)
		close(ad.events)
	})
}

func (ad *Adder) processEvents() {
	go func() {
		for event := range ad.events {
			ad.valueControl.mu.Lock()
			ad.valueControl.value += event.Value
			ad.valueControl.mu.Unlock()
		}
	}()
}
