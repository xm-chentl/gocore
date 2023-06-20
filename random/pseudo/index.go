package pseudo

import (
	"math/rand"
	"time"

	"github.com/xm-chentl/gocore/random"
)

type impl struct {
	random.WeightFunc

	desc string
}

func (i impl) Rand(max int) int {
	//  + int64(0.0123*float64(time.Second))
	//  + (155 * time.Hour.Nanoseconds())
	// rand 在并发情况会有大量sync.(*Mutex.LockSlow)消耗，重新实例一个处理
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(int(max))
}

func (i impl) RandWeight(items ...random.Item) *random.Item {
	return i.WeightFunc(i, items...)
}

func (i impl) Desc() string {
	return i.desc
}

// New 伪随机
func New() random.IRandom {
	return &impl{
		WeightFunc: random.RandWeightFunc,
		desc:       "伪随机算法",
	}
}
