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

func MakeThing() *Thing {
	return &Thing{}
}

func TestGetSingletonNoConstructor(t *testing.T) {
	// get instance, and modify something
	v1 := Singleton[Thing](nil)
	v1.ID = 25
	// modify something

	// get another instance to check to see if it
	// is the same one
	v2 := Singleton[Thing](nil)
	if v2.ID != v1.ID {
		t.Error("singleton did not produce a singleton")
	}
	// modify again...
	v2.Name = "Dr Manhattan"

	// get another instance to check to see if is
	// still the same one
	v3 := Singleton[Thing](nil)
	if v3.Name != v2.Name {
		t.Error("singleton did not produce a singleton")
	}
	if v3.ID != v2.ID {
		t.Error("singleton did not produce a singleton")
	}
	if v1.Name != v3.Name || v1.Name != v2.Name {
		t.Error("singleton did not produce a singleton")
	}
	if v3.ID != v1.ID {
		t.Error("singleton did not produce a singleton")
	}
	fmt.Printf("v1: %p %+v\n", v1, v1)
	fmt.Printf("v2: %p %+v\n", v2, v2)
	fmt.Printf("v3: %p %+v\n", v3, v3)
}

func TestGetSingletonWithConstructor(t *testing.T) {
	// get instance, and modify something
	v1 := Singleton(MakeThing)
	v1.ID = 25
	// modify something

	// get another instance to check to see if it
	// is the same one
	v2 := Singleton(MakeThing)
	if v2.ID != v1.ID {
		t.Error("singleton did not produce a singleton")
	}
	// modify again...
	(*v2).Name = "Dr Manhattan"

	// get another instance to check to see if is
	// still the same one
	v3 := Singleton(MakeThing)
	if v3.Name != v2.Name {
		t.Error("singleton did not produce a singleton")
	}
	if v3.ID != v2.ID {
		t.Error("singleton did not produce a singleton")
	}
	if v1.Name != v3.Name || v1.Name != v2.Name {
		t.Error("singleton did not produce a singleton")
	}
	if v3.ID != v1.ID {
		t.Error("singleton did not produce a singleton")
	}
	fmt.Printf("v1: %p %+v\n", v1, v1)
	fmt.Printf("v2: %p %+v\n", v2, v2)
	fmt.Printf("v3: %p %+v\n", v3, v3)
}

var result any

func BenchmarkGetSingletonNoConstructor(b *testing.B) {
	var v *Thing
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			v = Singleton[Thing](nil)
			if v.ID != j {
				b.Error("singleton did not produce a singleton")
			}
			v.ID++
		}
		v.ID = 0
	}
	result = v
}

func BenchmarkGetSingletonWithConstructor(b *testing.B) {
	var v *Thing
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			v = Singleton(MakeThing)
			if v.ID != j {
				b.Error("singleton did not produce a singleton")
			}
			v.ID++
		}
		v.ID = 0
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
