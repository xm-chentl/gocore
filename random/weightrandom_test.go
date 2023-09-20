package random

import (
	"math/rand"
	"testing"
	"time"
)

type impl struct {
	WeightFunc

	desc string
}

func (i impl) Rand(max int) int {
	//  + int64(0.0123*float64(time.Second))
	//  + (155 * time.Hour.Nanoseconds())
	// rand 在并发情况会有大量sync.(*Mutex.LockSlow)消耗，重新实例一个处理
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(int(max))
}

func (i impl) RandWeight(items ...Item) *Item {
	return i.WeightFunc(i, items...)
}

func (i impl) Desc() string {
	return i.desc
}

// New 伪随机
func New() IRandom {
	return &impl{
		WeightFunc: RandWeightFunc,
		desc:       "伪随机算法",
	}
}

func Test_Sort(t *testing.T) {
	data := SortArray{{Value: 1001, Weight: 150000}, {Value: 1002, Weight: 130000}, {Value: 1003, Weight: 100000}, {Value: 1004, Weight: 70000}, {Value: 1005, Weight: 60000}, {Value: 1006, Weight: 50000}, {Value: 1007, Weight: 40000}, {Value: 1008, Weight: 40000}, {Value: 1009, Weight: 30000}, {Value: 1010, Weight: 60000}, {Value: 1011, Weight: 20000}, {Value: 1011, Weight: 20000}, {Value: 1012, Weight: 1280}, {Value: 1013, Weight: 640}, {Value: 1014, Weight: 320}, {Value: 1015, Weight: 160}, {Value: 1016, Weight: 80}, {Value: 1017, Weight: 40}, {Value: 1018, Weight: 20}, {Value: 1019, Weight: 10}, {Value: 1020, Weight: 5}}
	d := make(map[int]int)
	r := &impl{}
	for i := 0; i < 1000000; i++ {
		randItem := RandWeightFunc(r, data...)
		id := randItem.Value.(int)
		d[id] += 1
	}
	t.Fatal(d)
}
