package bitset

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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

			b := New(cfg)

			assert.EqualValues(t, b.Count(), 0, "new bitset should have count 0")
		})
	}
}

func TestBitsetMaxIndex(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := New(cfg)

			var bits uint = 0
			for _, t := range cfg {
				bits += t
			}

			max := (1 << bits) - 1
			assert.EqualValues(t, max, b.MaxIndex())
		})
	}
}

func TestBitsetSetClear(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := New(cfg)

			i := b.Cap() / 2

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

			i = b.Cap() - 1
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

			b := New(cfg)

			idx, found := b.NextSet(0)
			assert.EqualValues(t, found, false)

			for i := uint64(0); i < b.Cap(); i++ {

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

func TestBitsetNextClear(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := New(cfg)

			idx, found := b.NextClear(0)
			assert.EqualValues(t, found, true)
			assert.EqualValues(t, idx, 0)

			i := b.Cap() / 2

			b.Set(i)
			idx, found = b.NextClear(i)
			assert.EqualValues(t, true, found)
			assert.EqualValues(t, i+1, idx)

			b.Clear(i)
			idx, found = b.NextClear(i)
			assert.EqualValues(t, true, found)
			assert.EqualValues(t, i, idx)
		})
	}
}
