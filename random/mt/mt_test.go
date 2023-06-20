package mt

import (
	"testing"
)

func BenchmarkRand(b *testing.B) {
	m := make(map[int]int)
	for i := 0; i < b.N; i++ {
		impl := New()
		v := impl.Rand(100)
		m[v] = +1
	}
}
