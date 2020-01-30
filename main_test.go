package main

import (
	"testing"
)

func BenchmarkContext(b *testing.B){
	for i := 0; i < b.N; i++ {
		Hasher(64)
	}
}

