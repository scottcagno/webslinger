package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func main() {

}

type TimedMap struct {
	m sync.Map

	interval           time.Duration
	cleaningTicker     *time.Ticker
	cleaningTickerStop chan bool
	stoppedTicker      bool
	entries            uint64
}

type Poller struct {
	timeout  time.Duration
	sessions *sync.Map

	ticker *time.Ticker
	ctx    context.Context
	cancel context.CancelFunc
}

func (ss *Poller) cleanUpRoutine() {
	// When we receive a "tick", we should loop through the
	// sessions, checking to see if any of them are expired.
	// If we find any that are expired, we should remove them.
	for {
		select {
		case t := <-ss.ticker.C:
			// Clean up expired sessions
			log.Printf("Checking for expired sessions: %v\n", t)
			ss.sessions.Range(
				func(sid, session any) bool {
					if session.(Session).ExpiresIn() < 1 {
						ss.sessions.Delete(sid)
					}
					return true
				},
			)
		case <-ss.ctx.Done():
			ss.ticker.Stop()
			return
		}
	}
}

func (ss *Poller) close() {
	log.Printf("*sessionStore.Close has been called.\n")
	// stop the ticker and free any other
	// resources.
	ss.ticker.Stop()
	ss.cancel()
}
