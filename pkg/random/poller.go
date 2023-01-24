package random

import (
	"errors"
	"log"
	"time"
)

var ErrStop = errors.New("stopping poller")

type PollerFunc func(t time.Time) error

type Poller struct {
	fn         PollerFunc
	interval   time.Duration
	ticker     *time.Ticker
	tickerStop chan bool
	isRunning  bool
}

func NewPoller(fn PollerFunc, interval time.Duration) *Poller {
	p := &Poller{
		fn:         fn,
		interval:   interval,
		ticker:     time.NewTicker(interval),
		tickerStop: make(chan bool),
	}
	p.startPolling()
	return p
}

func (p *Poller) runPollingFunc(t time.Time) {
	if err := p.fn(t); err == ErrStop {
		log.Println("[Poller] got an error, stopping...")
		p.StopPolling()
	}
}

func (p *Poller) startPolling() {
	// already running
	if p.isRunning {
		return
	}
	// spawn a new goroutine
	go func() {
		for {
			select {
			case t := <-p.ticker.C: // a new tick
				p.runPollingFunc(t)
			case <-p.tickerStop: // stopping the ticker
				break
			}
		}
	}()
	p.isRunning = true
}

func (p *Poller) StopPolling() {
	// already stopped
	if !p.isRunning {
		return
	}
	// stop the ticker
	p.ticker.Stop()
	// stop the poller
	go func() {
		p.tickerStop <- true
		return
	}()
	p.isRunning = false
}

func (p *Poller) RestartPolling() {
	// stop the poller
	p.StopPolling()

	if p.ticker == nil {
		// set new ticker
		p.ticker = time.NewTicker(p.interval)
	} else {
		// otherwise, reset existing ticker
		p.ticker.Reset(p.interval)
	}
	// restart the poller
	p.startPolling()
}

func (p *Poller) RestartPollerWithInterval(interval time.Duration) {
	// stop the poller
	p.StopPolling()

	// set the new interval
	p.interval = interval

	if p.ticker == nil {
		// set new ticker
		p.ticker = time.NewTicker(p.interval)
	} else {
		// otherwise, reset existing ticker
		p.ticker.Reset(p.interval)
	}
	// restart the poller
	p.startPolling()
}
