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
	[]uint{16, 8, 4, 2, 1},
	[]uint{1, 0, 1, 0, 1, 0},
	[]uint{8, 8, 8, 8, 8, 8, 8, 8},
}

func TestBitsetNewCount(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := New(cfg)

			assert.EqualValues(t, b.Count(), 0, "new bitset should have count 0")
		})
	}
}

func TestBitsetCap(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := New(cfg)

			var bits uint = 0
			for _, t := range cfg {
				bits += t
			}

			cap := (1 << bits) - 1
			assert.EqualValues(t, cap, b.Cap())
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

func TestBitsetFindSet(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := New(cfg)

			idx, found := b.FindSet(0)
			assert.EqualValues(t, found, false)

			i := b.Cap() / 2

			b.Set(i)
			idx, found = b.FindSet(0)
			assert.EqualValues(t, true, found)
			assert.EqualValues(t, i, idx)

			idx, found = b.FindSet(i)
			assert.EqualValues(t, true, found)
			assert.EqualValues(t, i, idx)

			idx, found = b.FindSet(i + 1)
			assert.EqualValues(t, false, found)

			b.Clear(i)
			idx, found = b.FindSet(0)
			assert.EqualValues(t, false, found)

			idx, found = b.FindSet(0)
			assert.EqualValues(t, false, found)

			b.Set(b.Cap() - 1)
			idx, found = b.FindSet(0)
			assert.EqualValues(t, true, found)
			assert.EqualValues(t, b.Cap()-1, idx)
		})
	}
}

func TestBitsetFindClear(t *testing.T) {

	for _, cfg := range configs {
		t.Run(fmt.Sprintf("%v", cfg), func(t *testing.T) {

			b := New(cfg)

			idx, found := b.FindClear(0)
			assert.EqualValues(t, found, true)
			assert.EqualValues(t, idx, 0)

			i := b.Cap() / 2

			b.Set(i)
			idx, found = b.FindClear(i)
			assert.EqualValues(t, true, found)
			assert.EqualValues(t, i+1, idx)

			b.Clear(i)
			idx, found = b.FindClear(i)
			assert.EqualValues(t, true, found)
			assert.EqualValues(t, i, idx)
		})
	}
}
