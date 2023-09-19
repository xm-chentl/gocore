package random

import (
	"sort"
	"testing"
)

func Test_Sort(t *testing.T) {
	data := SortArray{
		Item{
			Value:  "1",
			Weight: 1000,
		},
		{
			Value:  "2",
			Weight: 888,
		},
		{
			Value:  "3",
			Weight: 777,
		},
		{
			Value:  "4",
			Weight: 666,
		},
		{
			Value:  "5",
			Weight: 333,
		},
		{
			Value:  "6",
			Weight: 666,
		},
	}
	sort.Sort(data)
	t.Fatal(data)
}
