package bitset_example

import (
	"bitset"

	"fmt"
)

func ExampleBitset() {

	b := bitset.New([]uint{16, 8, 4})

	idx, found := b.NextClear(0)
	fmt.Printf("NextClear(%d): idx=%v found=%v\n", 0, idx, found)

	i := b.Cap() / 2

	b.Set(i)
	fmt.Printf("Set(%d)\n", i)

	idx, found = b.NextClear(i)
	fmt.Printf("NextClear(%d): idx=%v found=%v\n", i, idx, found)

	b.Clear(i)
	fmt.Printf("Clear(%d)\n", i)

	idx, found = b.NextClear(i)
	idx, found = b.NextClear(i)
	fmt.Printf("NextClear(%d): idx=%v found=%v\n", i, idx, found)
}
