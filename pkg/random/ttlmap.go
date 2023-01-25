package random

import (
	"runtime"
	"sync"
	"time"
)

// NeverExpire is a constant that can be used if
// you want to add any entries that never expire
const NeverExpire time.Duration = -1

// entry is a TimeoutMap entry.
type entry[V any] struct {
	data    V
	expires int64 // unix timestamp
}

// TimeoutMap is a map that automatically collects
// and removes entries at a specified interval.
type TimeoutMap[K comparable, V any] struct {
	mu sync.Mutex
	m  map[K]*entry[V]

	interval       time.Duration
	ticker         *time.Ticker
	tickerStopChan chan bool
	isRunning      bool
	zeroValue      V
}

// NewTimeoutMap initializes and returns a new TimeoutMap instance
// setup to clean at the interval supplied in the constructor.
func NewTimeoutMap[K comparable, V any](interval time.Duration) *TimeoutMap[K, V] {
	var zeroValue V
	tm := &TimeoutMap[K, V]{
		m:              make(map[K]*entry[V]),
		interval:       interval,
		ticker:         time.NewTicker(interval),
		tickerStopChan: make(chan bool),
		zeroValue:      zeroValue,
	}
	tm.startCleaner()
	return tm
}

// clean is the internal method that iterates through
// the map and cleans up expired entries.
func (tm *TimeoutMap[K, V]) clean() {
	// lock
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// skip cleaning if the map is empty
	if len(tm.m) == 0 {
		return
	}
	// get the current time
	now := time.Now().UTC().Unix()
	// iterate over the map (deleting expired values)
	for k, e := range tm.m {
		if e.expires == int64(NeverExpire) {
			continue
		}
		if now > e.expires {
			delete(tm.m, k)
		}
	}
}

// startCleaner is the internal method to initially
// start up the cleaner. This is not exported and
// should only be called by the other methods.
func (tm *TimeoutMap[K, V]) startCleaner() {
	// already running
	if tm.isRunning {
		return
	}
	// spawn a new goroutine
	go func() {
		for {
			select {
			case <-tm.ticker.C: // a new tick
				tm.clean()
			case <-tm.tickerStopChan: // stopping the ticker
				break
			}
		}
	}()
	// now that it's running, set to true
	tm.isRunning = true
}

// StopCleaner stops the cleaner from running. It will
// remain "off" until either RestartCleanerWithInterval
// or RestartCleaner is called.
func (tm *TimeoutMap[K, V]) StopCleaner() {
	// already stopped
	if !tm.isRunning {
		return
	}
	// stop the ticker
	tm.ticker.Stop()
	// stop the cleaner
	go func() {
		tm.tickerStopChan <- true
		return
	}()
	// make sure to set isRunning to false
	// now that we stopped the ticker
	tm.isRunning = false
}

// RestartCleaner simply restarts the cleaner using the
// same interval. It can be thought of like a "StartCleaner"
// type of method, but it ensures that the cleaner is actually
// stopped before restarting.
func (tm *TimeoutMap[K, V]) RestartCleaner() {
	tm.RestartCleanerWithInterval(tm.interval)
}

// RestartCleanerWithInterval stops the cleaner if it is
// currently running, updates the interval and then starts
// the cleaner again at the provided interval.
func (tm *TimeoutMap[K, V]) RestartCleanerWithInterval(interval time.Duration) {
	// stop the cleaner
	tm.StopCleaner()

	// set the new interval
	tm.interval = interval

	if tm.ticker == nil {
		// set new ticker
		tm.ticker = time.NewTicker(tm.interval)
	} else {
		// otherwise, reset existing ticker
		tm.ticker.Reset(tm.interval)
	}
	// restart the cleaner
	tm.startCleaner()
}

// CleanNow calls the cleaner manually, now.
func (tm *TimeoutMap[K, V]) CleanNow() {
	tm.clean()
}

func getTime(expiry time.Duration) int64 {
	if expiry == NeverExpire {
		return -1
	}
	return time.Now().Add(expiry).UTC().Unix()
}

// Put writes the key and value to the map overwriting any
// existing entry that might be in there already. It takes
// an expiry time and sets the entry to automatically expire
// and be cleaned up (removed) by the cleaner at that time.
// If you wish to add an entry that never expires, set the
// expiry time to -1 or use the built-in NeverExpire.
func (tm *TimeoutMap[K, V]) Put(k K, v V, expiry time.Duration) {
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
	e.expires = getTime(expiry)
	tm.m[k] = e
}

// Get attempts to return the value and a found
// boolean from the map using the map key.
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

// Del removes the entry with the matching key
// from the map immediately.
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

// LazyDel marks the entry for deletion during
// the next cleaning cycle.
func (tm *TimeoutMap[K, V]) LazyDel(k K) {
	// lock
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// check to see if the entry exists
	e, exists := tm.m[k]
	if !exists {
		// do nothing, it is not there
		return
	}
	// set the entry to expire during the
	// next cleanup cycle
	e.expires = time.Now().UTC().Unix()
	tm.m[k] = e
}

// IsCleanerRunning returns the status of weather or
// not the cleaner is currently running.
func (tm *TimeoutMap[K, V]) IsCleanerRunning() bool {
	// lock
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.isRunning
}

// Clear clears out all entries in the map.
func (tm *TimeoutMap[K, V]) Clear() {
	// lock
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// clear the map
	tm.m = nil
	runtime.GC()
	tm.m = make(map[K]*entry[V])
	tm.RestartCleaner()
}

func (tm *TimeoutMap[K, V]) Range(fn func(k K, v V, remaining time.Duration) bool) {
	// lock
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// range
	for k, e := range tm.m {
		if !fn(k, e.data, time.Until(time.Unix(e.expires, 0))) {
			break
		}
	}
}
