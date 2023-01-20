package random

import (
	"testing"
	"time"
)

var tmap = NewTimedMap(time.Second * 10)
var expiresAt = time.Duration(time.Now().Add(time.Minute).Unix())

func BenchmarkTimedMap_SetTemporary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tmap.Set("some key", "some value", expiresAt)
	}
}

func BenchmarkTimedMap_Get(b *testing.B) {
	tmap.Set("some key", "some value", expiresAt)
	for i := 0; i < b.N; i++ {
		_, ok := tmap.Get("some key")
		if !ok {
			b.Fail()
		}
	}
}
