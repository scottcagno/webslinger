package sessions

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrInvalidEntryType = errors.New("invalid Entry type")
var ErrNotFound = errors.New("not found")

type Entry interface {
	ExpiresAt() time.Time
}

type entry struct {
	data    any
	expires time.Time
}

func (e entry) ExpiresAt() time.Time {
	return e.expires
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
	// iterate over the map (deleting expired values)
	tm.m.Range(func(k, v any) bool {
		if e, ok := v.(Entry); ok {
			if time.Until(e.ExpiresAt()) < 1 {
				_, exists := tm.m.LoadAndDelete(k)
				if exists {
					// key was in the map so, decrement entry count
					decr(&tm.entries)
				}
				// otherwise, existing entry was not found so, there is
				// no need to modify the entry count
			}
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

func (tm *TimedMap) SetEntry(k string, e Entry) error {
	if _, ok := e.(Entry); !ok {
		return ErrInvalidEntryType
	}
	_, exists := tm.m.LoadOrStore(k, e)
	if !exists {
		// key was not in the map so, increment entry count
		incr(&tm.entries)
	}
	// otherwise, existing entry was not found so, there is
	// no need to modify the entry count
	return nil
}

func (tm *TimedMap) GetEntry(k string) (Entry, error) {
	v, exists := tm.m.Load(k)
	if !exists {
		return nil, ErrNotFound
	}
	e, ok := v.(Entry)
	if !ok {
		return nil, ErrInvalidEntryType
	}
	return e, nil
}

func (tm *TimedMap) DelEntry(k string) (Entry, error) {
	v, exists := tm.m.LoadAndDelete(k)
	if exists {
		// key was in the map so, decrement entry count
		decr(&tm.entries)
	}
	// otherwise, existing entry was not found so, there is
	// no need to modify the entry count
	e, ok := v.(Entry)
	if !ok {
		return nil, ErrInvalidEntryType
	}
	return e, nil
}

func (tm *TimedMap) Set(k string, v any, expires time.Duration) error {
	return tm.SetEntry(k, entry{
		data:    v,
		expires: time.Now().Add(expires),
	})
}

func (tm *TimedMap) Get(k string) (any, bool) {
	e, err := tm.GetEntry(k)
	if err != nil {
		return nil, false
	}
	return e, true
}

func (tm *TimedMap) Del(k string) {
	tm.DelEntry(k)
}
