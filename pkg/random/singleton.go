package random

import (
	"reflect"
	"sync"
)

// GetSingleton creates a "singleton"--it ensures that it can only
// ever be created one time, and all access is to the same version.
// It accepts an optional constructor function (or nil) and in the
// case where a constructor function is provided it will use it.
var singletons sync.Map

// Singleton returns a singleton of T.
func Singleton[T any]() (t *T) {
	hash := reflect.TypeOf(t)
	v, ok := singletons.Load(hash)
	if ok {
		return v.(*T)
	}
	v = new(T)
	v, _ = singletons.LoadOrStore(hash, v)
	return v.(*T)
}
