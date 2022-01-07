package main

import "testing"

func BenchmarkAtomic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CountAtomic()
	}
}

func BenchmarkMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CountMutex()
	}
}
