package jwt

import (
	"crypto"
	"testing"
)

func benchmarkValidateToken(b *testing.B, validator validator, raw RawToken) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				_, err := validator.ValidateRawToken(raw, validator.Method)
				if err != nil {
					b.Fatal(err)
				}
			}
		},
	)
}

func benchmarkParseRawToken(b *testing.B, raw RawToken) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				_, err := ParseRawToken(raw)
				if err != nil {
					b.Fatal(err)
				}
			}
		},
	)
}

// Helper method for benchmarking various signing methods
func benchmarkSigning(b *testing.B, method SigningMethod, key crypto.PrivateKey) {
	b.Helper()
	t, err := NewToken(method, nil, key)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				_, err := method.Sign(t.SigningSection(), key)
				if err != nil {
					b.Fatal(err)
				}
			}
		},
	)
}
