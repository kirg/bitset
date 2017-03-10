package main

import (
	"bitset"
	"fmt"
)

func main() {

	b := bitset.New([]uint{2, 2, 2, 2})

	idx, found := b.FindClear(0)
	fmt.Printf("FindClear(%d): idx=%v found=%v\n", 0, idx, found)

	cap := b.Cap()
	fmt.Printf("Cap(): %d\n", cap)

	stats := b.Stats()
	fmt.Printf("Stats(): %v\n", stats)

	i := cap - 1
	if b.Set(i) != nil {
		fmt.Printf("Set(%x)\n", i)
	} else {
		fmt.Printf("Set(%x) failed\n", i)
	}

	idx, found = b.FindClear(i)
	fmt.Printf("FindClear(%x): idx=%x found=%v\n", i, idx, found)

	for i := uint64(0); i < b.Cap(); i++ {
		b.Set(i)
		stats = b.Stats()
		fmt.Printf("Set(%d) -> %v\n", i, stats)
	}
}
