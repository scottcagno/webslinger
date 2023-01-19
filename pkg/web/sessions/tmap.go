package sessions

import (
	"sync"
	"sync/atomic"
	"time"
)

type entry struct {
	ttl   time.Duration
	value any
}

type TimedMap struct {
	m sync.Map

	interval           time.Duration
	cleaningTicker     *time.Ticker
	cleaningTickerStop chan bool
	stoppedTicker      bool
	entries            uint64
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

func NewTimedMap(interval time.Duration) *TimedMap {
	tm := &TimedMap{
		interval:           interval,
		cleaningTicker:     time.NewTicker(interval),
		cleaningTickerStop: make(chan bool),
		stoppedTicker:      true,
	}
	tm.StartCleaner()
	return tm
}

func (tm *TimedMap) clean() {
	// skip cleaning if the map is empty
	if tm.Len() == 0 {
		return
	}
	// current time as unix timestamp
	now := time.Now().Unix()
	// iterate over the map (deleting expired values)
	tm.m.Range(func(k, v any) bool {
		if now >= int64(v.(entry).ttl) {
			_, exists := tm.m.LoadAndDelete(k)
			if exists {
				// key was in the map so, decrement entry count
				decr(&tm.entries)
			}
			// otherwise, existing entry was not found so, there is
			// no need to modify the entry count
		}
		return true
	})
}

func (tm *TimedMap) StartCleaner() {
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

func (tm *TimedMap) StopCleaner() {
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

func (tm *TimedMap) RestartCleanerWithInterval(interval time.Duration) {
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

func (tm *TimedMap) CleanNow() {
	tm.clean()
}

func (tm *TimedMap) Len() int {
	return int(check(&tm.entries))
}

func (tm *TimedMap) Set(k string, v any, expires time.Duration) {
	_, exists := tm.m.LoadOrStore(k, entry{expires, v})
	if !exists {
		// key was not in the map so, increment entry count
		incr(&tm.entries)
	}
	// otherwise, existing entry was not found so, there is
	// no need to modify the entry count
}

func (tm *TimedMap) Get(k string) (any, bool) {
	v, exists := tm.m.Load(k)
	return v, exists
}

func (tm *TimedMap) Del(k string) {
	_, exists := tm.m.LoadAndDelete(k)
	if exists {
		// key was in the map so, decrement entry count
		decr(&tm.entries)
	}
	// otherwise, existing entry was not found so, there is
	// no need to modify the entry count
}
