package random

type Item struct {
	Value  interface{}
	Weight int
}

type SortArray []Item

func (s SortArray) Len() int {
	return len(s)
}

func (s SortArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortArray) Less(i, j int) bool {
	return s[j].Weight < s[i].Weight
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
