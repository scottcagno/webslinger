package tests

import (
	"testing"
)

type MatchTest struct {
	pattern, s string
	match      bool
	err        error
}

var matchTests = []MatchTest{
	{"h?llo", "hello", true, nil},
	{"h?llo", "hallo", true, nil},
	{"h?llo", "hxllo", true, nil},
	{"h*llo", "hllo", true, nil},
	{"h*llo", "heeeeeello", true, nil},
	// {"h[ae]llo", "hello", true, nil},
	// {"h[ae]llo", "hallo", true, nil},
	// {"h[ae]llo", "hillo", false, nil},
	// {"h[^e]llo", "hallo", true, nil},
	// {"h[^e]llo", "hbllo", true, nil},
	// {"h[^e]llo", "hxllo", true, nil},
	// {"h[^e]llo", "hello", false, nil},
	// {"h[a-b]llo", "hallo", true, nil},
	// {"h[a-b]llo", "hbllo", true, nil},
	// {"h[a-b]llo", "hello", false, nil},
	{"*name*", "firstname", true, nil},
	{"*name*", "lastname", true, nil},
	{"a??", "firstname", false, nil},
	{"a??", "age", true, nil},
	{"abc", "abc", true, nil},
	{"*", "abc", true, nil},
	{"*c", "abc", true, nil},
	{"a*", "a", true, nil},
	{"a*", "abc", true, nil},
	{"a*b?c*x", "abxbbxdbxebxczzx", true, nil},
	{"a*b?c*x", "abxbbxdbxebxczzy", false, nil},
	// {"ab[c]", "abc", true, nil},
	// {"ab[b-d]", "abc", true, nil},
	// {"ab[e-g]", "abc", false, nil},
	// {"ab[^c]", "abc", false, nil},
	// {"ab[^b-d]", "abc", false, nil},
	// {"ab[^e-g]", "abc", true, nil},
	// {"a?b", "a☺b", true, nil},
	// {"a[^a]b", "a☺b", true, nil},
	// {"a???b", "a☺b", false, nil},
	// {"a[^a][^a][^a]b", "a☺b", false, nil},
	// {"[a-ζ]*", "α", true, nil},
	// {"*[a-ζ]", "A", false, nil},
	{"a?b", "a_b", true, nil},
	{"a*b", "a_b", true, nil},
	{"*x", "xxx", true, nil},
}

func TestMatchV1(t *testing.T) {
	for _, tt := range matchTests {
		ok, err := MatchV1(tt.pattern, tt.s)
		if tt.match != ok || err != tt.err {
			t.Errorf("Match(%#q, %#q) = %v, %v want %v, %v", tt.pattern, tt.s, ok, err, tt.match, tt.err)
		}
	}
}

func TestMatchV2(t *testing.T) {
	for _, tt := range matchTests {
		ok, err := MatchV2(tt.pattern, tt.s)
		if tt.match != ok || err != tt.err {
			t.Errorf("Match(%#q, %#q) = %v, %v want %v, %v", tt.pattern, tt.s, ok, err, tt.match, tt.err)
		}
	}
}

func TestMatchV3(t *testing.T) {
	for _, tt := range matchTests {
		ok, err := MatchV3(tt.pattern, tt.s)
		if tt.match != ok || err != tt.err {
			t.Errorf("Match(%#q, %#q) = %v, %v want %v, %v", tt.pattern, tt.s, ok, err, tt.match, tt.err)
		}
	}
}

var benchmarks = []struct {
	name string
	fn   func(p, s string) (bool, error)
}{
	{"MatchV1", MatchV1},
	{"MatchV2", MatchV2},
	{"MatchV3", MatchV3},
}

func BenchmarkMatchers(b *testing.B) {
	for _, tt := range benchmarks {
		b.Run(
			tt.name, func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					for _, test := range matchTests {
						ok, err := tt.fn(test.pattern, test.s)
						if test.match != ok || err != test.err {
							b.Errorf(
								"Match(%#q, %#q) = %v, %v want %v, %v",
								test.pattern, test.s, ok, err, test.match, test.err,
							)
						}
					}
				}
			},
		)
	}
}
