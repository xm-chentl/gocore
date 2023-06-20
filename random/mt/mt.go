package mt

import (
	"math/rand"
	"time"

	"github.com/xm-chentl/gocore/random"

	"github.com/seehuhn/mt19937"
)

// impl 梅森旋转算法（Mersenne Twister）
// 效率快了四位
type impl struct {
	random.WeightFunc

	desc string
}

func (i impl) Rand(max int) int {
	mt := rand.New(mt19937.New())
	mt.Seed(time.Now().UnixNano())
	return mt.Intn(int(max))
}

func (i impl) RandWeight(items ...random.Item) *random.Item {
	return i.WeightFunc(i, items...)
}

func (i impl) Desc() string {
	return i.desc
}

func New() random.IRandom {
	return &impl{
		WeightFunc: random.RandWeightFunc,
		desc:       "梅森旋转算法(Mersenne Twister), 比go自带的随机快将近四倍;",
	}
}
