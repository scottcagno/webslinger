package random

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestPoller(t *testing.T) {

	fn := func(t time.Time) error {
		fmt.Printf("%s\n", t.Format(time.Layout))
		return nil
	}

	log.Println(">> Starting at 2 second intervals... (sleep=10)")
	p := NewPoller(fn, 2*time.Second)
	time.Sleep(10 * time.Second)

	log.Println(">> Stopping... (sleep=10)")
	p.StopPolling()
	time.Sleep(10 * time.Second)

	log.Println(">> Re-Starting... (sleep=10)")
	p.RestartPolling()
	time.Sleep(5 * time.Second)

	time.Sleep(5 * time.Second)

	log.Println(">> Restarting at 4 second intervals... (sleep=20)")
	p.RestartPollerWithInterval(4 * time.Second)
	time.Sleep(20 * time.Second)

	log.Println(">> Stopping... (sleep=5)")
	p.StopPolling()
	time.Sleep(5 * time.Second)
}

func TestPollerWithErr(t *testing.T) {

	var n int

	fn := func(t time.Time) error {
		n++
		if n > 7 {
			return ErrStop
		}
		fmt.Printf("%s\n", t.Format(time.Layout))
		return nil
	}

	log.Println(">> Starting at 2 second intervals... (sleep=10)")
	p := NewPoller(fn, 2*time.Second)
	time.Sleep(10 * time.Second)

	log.Println(">> Stopping... (sleep=10)")
	p.StopPolling()
	time.Sleep(10 * time.Second)

	log.Println(">> Re-Starting... (sleep=10)")
	p.RestartPolling()
	time.Sleep(5 * time.Second)

	time.Sleep(5 * time.Second)

	log.Println(">> Restarting at 4 second intervals... (sleep=20)")
	p.RestartPollerWithInterval(4 * time.Second)
	time.Sleep(20 * time.Second)

	log.Println(">> Stopping... (sleep=5)")
	p.StopPolling()
	time.Sleep(5 * time.Second)
}
