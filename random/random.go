package random

type IRandom interface {
	// Rand 随机
	Rand(int) int
	// RandWeight 指定权重随机
	RandWeight(items ...Item) *Item
	// Desc 描述
	Desc() string
}
