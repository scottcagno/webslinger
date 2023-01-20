package random

import (
	"reflect"
	"sync"
)

var singletonCache sync.Map

// GetSingleton creates a "singleton"--it ensures that it can only
// ever be created one time, and all access is to the same version.
func GetSingleton[T any]() (t *T) {
	hash := reflect.TypeOf(t)
	instance, hasInstance := singletonCache.Load(hash)
	if !hasInstance {
		instance = new(T)
		instance, _ = singletonCache.LoadOrStore(hash, instance)
	}
	return instance.(*T)
}
