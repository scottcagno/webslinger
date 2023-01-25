package main

import (
	"fmt"
	"runtime"
	"time"
)

const valueSize = 255

var m = make(map[int][]byte)

func addKeys() {
	for i := 0; i < 100000; i++ {
		m[i] = make([]byte, valueSize)
	}
}

func main() {
	addKeys()
	fmt.Println("added keys to map")
	runtime.GC()
	time.Sleep(10 * time.Second)

	// fmt.Println("deleting keys")
	// for key := range m {
	//	delete(m, key)
	// }

	// runtime.GC()
	// time.Sleep(time.Second)

	fmt.Println("setting m = nil")
	m = nil

	time.Sleep(10 * time.Second)
	runtime.GC()
	time.Sleep(time.Second)
	runtime.GC()
	time.Sleep(time.Second)
}
