package random

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var tmap = NewTimeoutMap[string, int](time.Second * 10)
var expiresAt = time.Duration(time.Now().Add(time.Minute).Unix())

func newMap() *TimeoutMap[string, int] {
	return NewTimeoutMap[string, int](time.Second)
}

func addEntries(tm *TimeoutMap[string, int]) {
	// Add entries...
	// 5 that will expire in 10 seconds
	// 5 that will expire in 20 seconds
	// 5 that will expire in 30 seconds
	// 5 that will expire in 40 seconds
	// 5 that will expire in 50 seconds
	// 5 that will never expire
	for i := 0; i < 30; i++ {
		switch {
		case i > 0 && i < 5:
			tm.Put(strconv.Itoa(i), i, 10*time.Second)
		case i > 4 && i < 10:
			tm.Put(strconv.Itoa(i), i, 20*time.Second)
		case i > 9 && i < 15:
			tm.Put(strconv.Itoa(i), i, 30*time.Second)
		case i > 14 && i < 20:
			tm.Put(strconv.Itoa(i), i, 40*time.Second)
		case i > 19 && i < 25:
			tm.Put(strconv.Itoa(i), i, 50*time.Second)
		case i > 24 && i < 30:
			tm.Put(strconv.Itoa(i), i, NeverExpire)
		}
	}
}

func printMap(tm *TimeoutMap[string, int]) {
	var count int
	tm.Range(func(k string, v int, remaining time.Duration) bool {
		if k != "" {
			count++
		}
		return true
	})
	fmt.Printf(">> %d entries remaining...\n", count)
}

func TestTimeoutMap(t *testing.T) {
	tm := newMap()
	addEntries(tm)

	for i := 0; i < 10; i++ {
		printMap(tm)
		if i == 3 {
			fmt.Println(">> Stopping the cleaner...")
			tm.StopCleaner()
		}
		if i == 6 {
			fmt.Println(">> Restarting the cleaner...")
			tm.RestartCleaner()
		}
		fmt.Println("\nSleeping for 10 seconds...")
		time.Sleep(10 * time.Second)
	}
}

func BenchmarkTimedMap_SetTemporary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tmap.Put("some key", 1, expiresAt)
	}
}

func BenchmarkTimedMap_Get(b *testing.B) {
	tmap.Put("some key", 1, expiresAt)
	for i := 0; i < b.N; i++ {
		_, ok := tmap.Get("some key")
		if !ok {
			b.Fail()
		}
	}
}
