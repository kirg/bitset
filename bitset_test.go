package bitset_test

import (
	"bitset"

	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var configs = [][]uint{
	[]uint{16},
	[]uint{7, 5},
	[]uint{4, 0, 0, 0, 0},
	[]uint{8, 4, 2, 1},
	[]uint{1, 0, 1, 0, 1, 0},
	[]uint{2, 2, 2, 2, 2, 2, 2},
}

func TestBitsetNewCount(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := bitset.New(cfg)

			assert.EqualValues(t, b.Count(), 0, "new bitset should have count 0")
		})
	}
}

func TestBitsetMaxIndex(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := bitset.New(cfg)

			var bits uint = 0
			for _, t := range cfg {
				bits += t
			}

			max := (1 << bits) - 1
			assert.EqualValues(t, max, b.Max())
		})
	}
}

func TestBitsetSetClear(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := bitset.New(cfg)

			i := b.Max() / 2

			assert.NotNil(t, b.Set(i))
			assert.EqualValues(t, b.Count(), 1)
			assert.EqualValues(t, b.Test(i), true)

			assert.NotNil(t, b.Set(i))
			assert.EqualValues(t, b.Count(), 1)
			assert.EqualValues(t, b.Test(i), true)

			assert.NotNil(t, b.Clear(i))
			assert.EqualValues(t, b.Count(), 0)
			assert.EqualValues(t, b.Test(i), false)

			assert.NotNil(t, b.Clear(i))
			assert.EqualValues(t, b.Count(), 0)
			assert.EqualValues(t, b.Test(i), false)

			i = b.Max()
			assert.EqualValues(t, b.Test(i), false)
			assert.NotNil(t, b.Clear(i))
			assert.EqualValues(t, b.Test(i), false)
			assert.NotNil(t, b.Set(i))
			assert.EqualValues(t, b.Test(i), true)
		})
	}
}

func TestBitsetNextSet(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := bitset.New(cfg)

			idx, found := b.NextSet(0)
			assert.EqualValues(t, found, false)

			idx, found = b.NextSet(math.MaxUint64)
			assert.EqualValues(t, found, false)

			for i := uint64(0); i <= b.Max(); i++ {

				b.Set(i)
				idx, found = b.NextSet(0)
				assert.EqualValues(t, true, found)
				assert.EqualValues(t, i, idx)

				idx, found = b.NextSet(i)
				assert.EqualValues(t, true, found)
				assert.EqualValues(t, i, idx)

				idx, found = b.NextSet(i + 1)
				assert.EqualValues(t, false, found)

				b.Clear(i)
				idx, found = b.NextSet(0)
				assert.EqualValues(t, false, found)
			}
		})
	}
}

func TestBitsetPrevSet(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := bitset.New(cfg)

			idx, found := b.PrevSet(0)
			assert.EqualValues(t, false, found)

			idx, found = b.PrevSet(b.Max())
			assert.EqualValues(t, false, found)

			idx, found = b.PrevSet(math.MaxUint64)
			assert.EqualValues(t, false, found)

			for i := uint64(0); i <= b.Max(); i++ {

				b.Set(i)
				idx, found = b.PrevSet(i)
				assert.EqualValues(t, true, found)
				assert.EqualValues(t, i, idx)

				idx, found = b.PrevSet(i + 1)
				assert.EqualValues(t, true, found)
				assert.EqualValues(t, i, idx)

				idx, found = b.PrevSet(b.Max())
				assert.EqualValues(t, true, found)
				assert.EqualValues(t, i, idx)

				b.Clear(i)
				idx, found = b.PrevSet(b.Max())
				assert.EqualValues(t, false, found)
			}
		})
	}
}

func TestBitsetNextClear(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := bitset.New(cfg)

			// set all
			b.ForEachClear(func(idx uint64) {
				b.Set(idx)
			})

			idx, found := b.NextClear(0)
			assert.EqualValues(t, found, false)

			for i := uint64(0); i <= b.Max(); i++ {

				b.Clear(i)
				idx, found = b.NextClear(0)
				assert.EqualValues(t, true, found)
				assert.EqualValues(t, i, idx)

				idx, found = b.NextClear(i)
				assert.EqualValues(t, true, found)
				assert.EqualValues(t, i, idx)

				idx, found = b.NextClear(i + 1)
				assert.EqualValues(t, false, found)

				b.Set(i)
				idx, found = b.NextClear(0)
				assert.EqualValues(t, false, found)
			}
		})
	}
}

func TestBitsetPrevClear(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := bitset.New(cfg)

			// set all
			b.ForEachClear(func(idx uint64) {
				b.Set(idx)
			})

			idx, found := b.PrevClear(0)
			assert.EqualValues(t, found, false)

			idx, found = b.PrevClear(math.MaxUint64)
			assert.EqualValues(t, found, false)

			for i := uint64(0); i <= b.Max(); i++ {

				b.Clear(i)
				idx, found = b.PrevClear(i)
				assert.EqualValues(t, true, found)
				assert.EqualValues(t, i, idx)

				idx, found = b.PrevClear(i + 1)
				assert.EqualValues(t, true, found)
				assert.EqualValues(t, i, idx)

				idx, found = b.PrevClear(b.Max())
				assert.EqualValues(t, true, found)
				assert.EqualValues(t, i, idx)

				b.Set(i)
				idx, found = b.PrevClear(math.MaxUint64)
				assert.EqualValues(t, false, found)
			}
		})
	}
}
