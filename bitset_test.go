package bitset

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitsetNew(t *testing.T) {

	b := New([]uint{16, 8, 4})
	assert.EqualValues(t, b.Count(), 0, "new bitset should have count 0")
}

func TestBitsetCap(t *testing.T) {

	b := New([]uint{16, 8, 4})

	cap := (1 << (16 + 8 + 4)) - 1
	assert.EqualValues(t, cap, b.Cap())
}

func TestBitsetSetClear(t *testing.T) {

	b := New([]uint{16, 8, 4})

	i := b.Cap() / 2

	b.Set(i)
	assert.EqualValues(t, b.Count(), 1)

	b.Set(i) // redundant set
	assert.EqualValues(t, b.Count(), 1)

	assert.EqualValues(t, b.Test(i), true)

	b.Clear(i)
	assert.EqualValues(t, b.Count(), 0)

	assert.EqualValues(t, b.Test(i), false)

	b.Clear(i)
	assert.EqualValues(t, b.Count(), 0)
}

func TestBitsetFindSet(t *testing.T) {

	b := New([]uint{16, 8, 4})

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
}

func TestBitsetFindClear(t *testing.T) {

	// b := New([]uint{16, 8, 4})
	b := New([]uint{4, 4, 4, 4})

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
}
