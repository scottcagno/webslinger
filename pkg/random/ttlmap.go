package random

import (
	"sync"
	"sync/atomic"
	"time"
)

type entry[T any] struct {
	data    T
	expires time.Time
}

type TimeoutMap[K comparable, V any] struct {
	mu sync.Mutex
	m  map[K]*entry[V]

	interval           time.Duration
	cleaningTicker     *time.Ticker
	cleaningTickerStop chan bool
	stoppedTicker      bool
	zeroValue          V
}

func NewTimeoutMap[K comparable, V any](interval time.Duration) *TimeoutMap[K, V] {
	var zeroValue V
	tm := &TimeoutMap[K, V]{
		m:                  make(map[K]*entry[V]),
		interval:           interval,
		cleaningTicker:     time.NewTicker(interval),
		cleaningTickerStop: make(chan bool),
		stoppedTicker:      true,
		zeroValue:          zeroValue,
	}
	tm.StartCleaner()
	return tm
}

func (tm *TimeoutMap[K, V]) clean() {
	// lock
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// skip cleaning if the map is empty
	if len(tm.m) == 0 {
		return
	}
	// iterate over the map (deleting expired values)
	for k, e := range tm.m {
		if time.Until(e.expires) < 1 {
			delete(tm.m, k)
		}
	}
}

func (tm *TimeoutMap[K, V]) StartCleaner() {
	// already running
	if !tm.stoppedTicker {
		return
	}
	// spawn a new goroutine
	go func() {
		for {
			select {
			case <-tm.cleaningTicker.C: // a new tick
				tm.clean()
			case <-tm.cleaningTickerStop: // stopping the ticker
				break
			}
		}
	}()
}

func (tm *TimeoutMap[K, V]) StopCleaner() {
	// already stopped
	if tm.stoppedTicker {
		return
	}

	// stop the ticker
	tm.cleaningTicker.Stop()

	// stop the cleaner
	go func() {
		tm.cleaningTickerStop <- true
		return
	}()
}

func (tm *TimeoutMap[K, V]) RestartCleanerWithInterval(interval time.Duration) {
	// stop the cleaner
	tm.StopCleaner()

	// set the new interval
	tm.interval = interval

	if tm.cleaningTicker == nil {
		// set new ticker
		tm.cleaningTicker = time.NewTicker(tm.interval)
	} else {
		// otherwise, reset existing ticker
		tm.cleaningTicker.Reset(tm.interval)
	}
	// restart the cleaner
	tm.StartCleaner()
}

func (tm *TimeoutMap[K, V]) CleanNow() {
	tm.clean()
}

func (tm *TimeoutMap[K, V]) Put(k K, v V, expires time.Duration) {
	// lock
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// check to see if the entry exists, in which case we can just update
	// it or if we need to create a new entry instance first.
	e, exists := tm.m[k]
	if !exists {
		// add a new entry
		e = new(entry[V])
	}
	// update entry
	e.data = v
	e.expires = time.Now().Add(expires)
	tm.m[k] = e
}

func (tm *TimeoutMap[K, V]) Get(k K) (V, bool) {
	// lock
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// check to see if the entry exists
	e, exists := tm.m[k]
	if !exists {
		// return empty value and false
		return tm.zeroValue, false
	}
	// otherwise, return correct value and true
	return e.data, true
}

func (tm *TimeoutMap[K, V]) Del(k K) {
	// lock
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// check to see if the entry exists
	_, exists := tm.m[k]
	if !exists {
		// do nothing, it is not there
		return
	}
	// otherwise, remove the entry
	delete(tm.m, k)
}

func incr(n *uint64) {
	atomic.AddUint64(n, 1)
}

func decr(n *uint64) {
	atomic.AddUint64(n, ^uint64(0))
}

func check(n *uint64) uint64 {
	return atomic.LoadUint64(n)
}
