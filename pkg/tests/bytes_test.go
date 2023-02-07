package tests

import (
	"bytes"
	"testing"
)

var (
	inputA = []byte("eyJhbGciOiJub25lIn0")
	inputB = []byte("eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc")
	want   = []byte("eyJhbGciOiJub25lIn0.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc")
)

type JoinFunc func(a, b []byte) []byte

func bytesJoin(a, b []byte) []byte {
	return bytes.Join([][]byte{a, b}, []byte{'.'})
}

func joinUsingCopy(a, b []byte) []byte {
	buf := make([]byte, len(a)+len(b)+1)
	n := copy(buf, a)
	copy(buf[n+1:], b)
	buf[n] = '.'
	return buf
}

func joinUsingBytesBuffer(a, b []byte) []byte {
	var buf bytes.Buffer
	buf.Grow(len(a) + len(b) + 1)
	buf.Write(a)
	buf.WriteByte('.')
	buf.Write(b)
	return buf.Bytes()
}

func joinUsingLoops(a, b []byte) []byte {
	c := make([]byte, len(a)+len(b)+1)
	var i int
	for _, ch := range a {
		c[i] = ch
		i++
	}
	c[i] = '.'
	i++
	for _, ch := range b {
		c[i] = ch
		i++
	}
	return c
}

func joinUsingAppend(a, b []byte) []byte {
	return append(a, append([]byte{'.'}, b...)...)
}

func BenchmarkBytesJoinFuncCompare(b *testing.B) {

	tests := []struct {
		Name string
		Func JoinFunc
	}{
		{
			"bytes.Join",
			bytesJoin,
		},
		{
			"joinUsingCopy",
			joinUsingCopy,
		},
		{
			"joinUsingBytesBuffer",
			joinUsingBytesBuffer,
		},
		{
			"joinUsingLoops",
			joinUsingLoops,
		},
		{
			"joinUsingAppend",
			joinUsingAppend,
		},
	}

	for _, tc := range tests {
		b.Run(
			tc.Name, func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					got := tc.Func(inputA, inputB)
					if !bytes.Equal(got, want) {
						b.Fatalf("got: %q, want: %q\n", got, want)
					}
				}
			},
		)
	}
}

// growSlice grows b by n, preserving the original content of b.
// If the allocation fails, it panics with ErrTooLarge.
func growSlice(b []byte, n int) []byte {
	defer func() {
		if recover() != nil {
			panic(bytes.ErrTooLarge)
		}
	}()
	// TODO(http://golang.org/issue/51462): We should rely on the append-make
	// pattern so that the compiler can call runtime.growslice. For example:
	//	return append(b, make([]byte, n)...)
	// This avoids unnecessary zero-ing of the first len(b) bytes of the
	// allocated slice, but this pattern causes b to escape onto the heap.
	//
	// Instead use the append-make pattern with a nil slice to ensure that
	// we allocate buffers rounded up to the closest size class.
	c := len(b) + n // ensure enough space for n elements
	if c < 2*cap(b) {
		// The growth rate has historically always been 2x. In the future,
		// we could rely purely on append to determine the growth rate.
		c = 2 * cap(b)
	}
	b2 := append([]byte(nil), make([]byte, c)...)
	copy(b2, b)
	return b2[:len(b)]
}
