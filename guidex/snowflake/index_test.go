package snowflake

import "testing"

func Benchmark_String(b *testing.B) {
	inst := snowflake{}
	pool := make(map[string]string)
	var v string
	for i := 0; i < b.N; i++ {
		v = inst.String()
		if _, ok := pool[v]; ok {
			b.Fatal("err")
		}
		pool[v] = v
	}
}

func Test_String(t *testing.T) {
	impl := New()
	t.Fatal(impl.String())
}
