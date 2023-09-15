package random

type Item struct {
	Value  interface{}
	Weight int
}

type WeightFunc func(IRandom, ...Item) *Item

func RandWeightFunc(r IRandom, items ...Item) (res *Item) {
	totalWeight := 0
	for _, item := range items {
		totalWeight += item.Weight
	}

	rnd := r.Rand(totalWeight)
	for _, item := range items {
		rnd -= item.Weight
		if rnd < 0 {
			return &item
		}
	}

	return
}
