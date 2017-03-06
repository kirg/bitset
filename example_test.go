package bitset

import (
	"fmt"
)

func ExampleBitset() {

	/*
		b := New([]uint{4, 8, 8, 2})

		var i uint64 = b.Cap() / 2

		v := b.Test(i)
		fmt.Printf("Test(%d) => %v\n", i, v)

		b.Set(i)
		fmt.Printf("Set(%d)\n", i)

		v = b.Test(i)
		fmt.Printf("Test(%d) => %v\n", i, v)

		b.Clear(i)
		fmt.Printf("Clear(%d)\n", i)

		v = b.Test(i)
		fmt.Printf("Test(%d) => %v\n", i, v)
	*/

	b := New([]uint{16, 8, 4})

	idx, found := b.FindClear(0)
	fmt.Printf("FindClear(%d): idx=%v found=%v\n", 0, idx, found)

	i := b.Cap() / 2

	b.Set(i)
	fmt.Printf("Set(%d)\n", i)

	idx, found = b.FindClear(i)
	fmt.Printf("FindClear(%d): idx=%v found=%v\n", i, idx, found)

	b.Clear(i)
	fmt.Printf("Clear(%d)\n", i)

	idx, found = b.FindClear(i)
	idx, found = b.FindClear(i)
	fmt.Printf("FindClear(%d): idx=%v found=%v\n", i, idx, found)
}
