package random

import (
	"fmt"
	"reflect"
	"testing"
)

type Thing struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestGetSingletonNonPtrNoConstructor(t *testing.T) {
	for i := 0; i < 10; i++ {
		t1 := GetSingleton[Thing](nil)
		addr(t1)
	}
}

func TestGetSingletonPtrToNoConstructor(t *testing.T) {
	for i := 0; i < 10; i++ {
		t1 := GetSingleton[*Thing](nil)
		t1.ID = 5
		addr(t1)
	}
}

var result any

func BenchmarkGetSingletonNonPtrNoConstructor(b *testing.B) {
	var v Thing
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			v = GetSingleton[Thing](nil)
		}
	}
	result = v
}

func BenchmarkGetSingletonPtrToNoConstructor(b *testing.B) {
	var v *Thing
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			v = GetSingleton[*Thing](nil)
		}
	}
	result = v
}

func comp(t1, t2 any) {
	fmt.Printf("t1: %p\nt2: %p\nsame: %v\n", t1, t2, reflect.DeepEqual(t1, t2))
}

func addr(v any) {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		fmt.Printf("address=%p\n", &v)
		return
	}
	fmt.Printf("address=%p\n", v)
}

func info(v any) {
	tof := reflect.TypeOf(v).Kind()
	fmt.Printf("type=%v\n", tof)
	switch tof {
	case reflect.Struct:
		fmt.Println("got struct")
	case reflect.Ptr:
		fmt.Println("got pointer")
	}
}
