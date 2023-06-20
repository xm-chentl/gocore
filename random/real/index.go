package real

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/xm-chentl/gocore/random"
)

type impl struct {
	random.WeightFunc
	seed []byte
	desc string
}

func (i impl) Rand(max int) int {
	rand.Read(i.seed[:])
	return int(binary.LittleEndian.Uint32(i.seed[:])) % max
}

func (i impl) Desc() string {
	return i.desc
}

func (i impl) RandWeight(items ...random.Item) *random.Item {
	return i.WeightFunc(i, items...)
}

// New 真随机
func New(len int) random.IRandom {
	if len == 0 {
		len = 64
	}
	return &impl{
		WeightFunc: random.RandWeightFunc,
		seed:       make([]byte, len),
		desc:       "真随机算法",
	}
}
