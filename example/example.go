package main

import (
	"bitset"
	"fmt"
)

func main() {

	b := bitset.New([]uint{5, 4, 3, 2, 1, 0})

	fmt.Printf("init -> max=%d count=%d stats=%v\nfinal -> count=%d stats=%v\n",
		b.MaxIndex(), b.Count(), b.Stats(),
		b.ForEachClear(func(i uint64) {
			b.Set(i)
			fmt.Printf("Set(%d) -> %v\n", i, b.Stats())
		}).ForEachSet(func(i uint64) {
			b.Clear(i)
			fmt.Printf("Clear(%d) -> %v\n", i, b.Stats())
		}).Count(), b.Stats())

	/*
		idx, found := b.NextClear(0)
		fmt.Printf("NextClear(%d): idx=%v found=%v\n", 0, idx, found)

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

		idx, found = b.NextClear(i)
		fmt.Printf("NextClear(%x): idx=%x found=%v\n", i, idx, found)
	*/

	idx, found := b.Set(b.MaxIndex()).NextSet(0)
	fmt.Printf("idx=%v found=%v\n", idx, found)

	/*
		for i := uint64(0); i <= b.Cap(); i++ {
			b.Set(i)
			fmt.Printf("Set(%d) -> %v\n", i, b.Stats())
		}
	*/

}
