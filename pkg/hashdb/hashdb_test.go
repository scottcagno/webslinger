package hashdb

import (
	"bytes"
	"fmt"
	"testing"
)

var testData = []struct {
	n   int
	rec *record
	dat []byte
}{
	{
		-1,
		&record{
			[]byte("rec-001"),
			[]byte("this is the first record"),
		},
		nil,
	},
	{
		-1,
		&record{
			[]byte("2"),
			[]byte("second record"),
		},
		nil,
	},
	{
		-1,
		&record{
			[]byte("003"),
			[]byte("record number three"),
		},
		nil,
	},
}

func init() {
	for i := range testData {
		b, n := encRecord(testData[i].rec.key, testData[i].rec.val)
		testData[i].n = n
		testData[i].dat = b
	}
}

func matchRecDat(r1, r2 []byte) bool {
	return bytes.Equal(r1, r2)
}

func matchRecPtr(r1, r2 *record) bool {
	return bytes.Equal(r1.key, r2.key) && bytes.Equal(r1.val, r2.val)
}

func TestRecord(t *testing.T) {

	// encode record #1
	r, n := encRecord([]byte("rec-001"), []byte("this is the first record"))
	if !matchRecDat(r, testData[0].dat) {
		t.Errorf("match fail: wanted=%v, got=%v\n", testData[0].dat, r)
	}
	fmt.Printf("record %d: size=%d, data=%v\n", 1, n, toString(r))

	// encode record #2
	r, n = encRecord([]byte("2"), []byte("second record"))
	if !matchRecDat(r, testData[1].dat) {
		t.Errorf("match fail: wanted=%v, got=%v\n", testData[0].dat, r)
	}
	fmt.Printf("record %d: size=%d, data=%v\n", 2, n, toString(r))

	// encode record #3
	r, n = encRecord([]byte("003"), []byte("record number three"))
	if !matchRecDat(r, testData[2].dat) {
		t.Errorf("match fail: wanted=%v, got=%v\n", testData[0].dat, r)
	}
	fmt.Printf("record %d: size=%d, data=%v\n", 3, n, toString(r))
}

func TestCBDHash(t *testing.T) {
	k1 := []byte("rec-001")
	h1 := cbdHash(k1)
	fmt.Printf("%q, %d\n", k1, h1&0xff)

	k2 := []byte("2")
	h2 := cbdHash(k2)
	fmt.Printf("%q, %d\n", k2, h2&0xff)

	k3 := []byte("003")
	h3 := cbdHash(k3)
	fmt.Printf("%q, %d\n", k3, h3&0xff)
}

func TestKeyOffset(t *testing.T) {
	keys := [][]byte{
		[]byte("some record"),
		[]byte("baz"),
		[]byte("key-001"),
		[]byte("3"),
		[]byte("key-099"),
		[]byte("user:12"),
		[]byte("bar"),
		[]byte("key:559199"),
		[]byte("6"),
		[]byte("id:991"),
		[]byte("keys"),
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
		[]byte("e"),
		[]byte("f"),
		[]byte("g"),
		[]byte("h"),
		[]byte("i"),
		[]byte("j"),
		[]byte("k"),
		[]byte("l"),
		[]byte("m"),
		[]byte("n"),
		[]byte("o"),
		[]byte("p"),
		[]byte("q"),
		[]byte("r"),
		[]byte("s"),
		[]byte("t"),
		[]byte("u"),
		[]byte("v"),
		[]byte("w"),
		[]byte("x"),
		[]byte("y"),
		[]byte("z"),
		[]byte("0"),
		[]byte("1"),
		[]byte("2"),
		[]byte("3"),
		[]byte("4"),
		[]byte("5"),
		[]byte("6"),
		[]byte("7"),
		[]byte("8"),
		[]byte("9"),
	}

	for i := range keys {
		fmt.Printf("key=%q, offset=%d\n", keys[i], getKeyHash(cbdHash, keys[i]))
	}
}
