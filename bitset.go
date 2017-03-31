package bitset

import (
	"fmt"
)

const getStats = false

type (
	Bitset interface {
		Test(idx uint64) bool
		Set(idx uint64) Bitset
		Clear(idx uint64) Bitset
		Swap(idx uint64, set bool) (swapped bool)
		Count() uint64
		Cap() uint64
		MaxIndex() uint64

		NextSet(start uint64) (idx uint64, found bool)
		NextClear(start uint64) (idx uint64, found bool)

		// PrevSet(start uint64) (idx uint64, found bool)
		// PrevClear(start uint64) (idx uint64, found bool)

		ForEachSet(do func(idx uint64)) Bitset
		ForEachClear(do func(idx uint64)) Bitset

		ForEachSetRange(start, end uint64, do func(idx uint64)) Bitset
		ForEachClearRange(start, end uint64, do func(idx uint64)) Bitset

		// SetBits(startIdx, endIdx uint64)
		// ClearBits(startIdx, endIdx uint64)

		Stats() []int // #stats
	}
)

type (
	bitset struct {
		root      node
		rootLevel *level
		count     uint64
		max       uint64
	}
)

func New(levelBits []uint) Bitset {

	rootLevel, max := initLevels(levelBits)

	if rootLevel == nil {
		return nil
	}

	return &bitset{
		root:      newNode(rootLevel, true, false),
		rootLevel: rootLevel,
		max:       max,
	}
}

func (t *bitset) Cap() uint64 {
	return t.max + 1 // could overflow!
}

func (t *bitset) MaxIndex() uint64 {
	return t.max
}

func (t *bitset) Test(idx uint64) (ret bool) {
	return t.root.test(t.rootLevel, idx)
}

func (t *bitset) Set(idx uint64) Bitset {

	if idx > t.max {
		return nil
	}

	if set, replace := t.root.set(t.rootLevel, idx); set {

		t.count++

		if replace != t.root {
			t.root = replace
		}
	}

	return t
}

func (t *bitset) Clear(idx uint64) Bitset {

	if idx > t.max {
		return nil
	}

	if cleared, replace := t.root.clr(t.rootLevel, idx); cleared {

		t.count--

		if replace != t.root {
			t.root = replace
		}
	}

	return t
}

func (t *bitset) Swap(idx uint64, set bool) bool {

	if idx > t.max {
		return false
	}

	var swapped bool

	if set {
		if swapped, _ := t.root.set(t.rootLevel, idx); swapped {
			t.count++
		}
	} else {
		if swapped, _ := t.root.clr(t.rootLevel, idx); swapped {
			t.count--
		}
	}

	return swapped
}

func (t *bitset) Count() uint64 {

	return t.count
}

func (t *bitset) NextSet(start uint64) (idx uint64, found bool) {

	return t.root.findset(t.rootLevel, start)
}

func (t *bitset) NextClear(start uint64) (idx uint64, found bool) {

	return t.root.findclr(t.rootLevel, start)
}

func (t *bitset) ForEachSet(do func(idx uint64)) Bitset {

	for i := uint64(0); i <= t.max; i++ {

		var found bool

		i, found = t.root.findset(t.rootLevel, i) // TODO: send down 'end' to terminate search sooner

		if !found {
			break // all done
		}

		do(i)
	}

	return t
}

func (t *bitset) ForEachClear(do func(idx uint64)) Bitset {

	for i := uint64(0); i <= t.max; i++ {

		var found bool

		i, found = t.root.findclr(t.rootLevel, i) // TODO: send down 'end' to terminate search sooner

		if !found {
			break // all done
		}

		do(i)
	}

	return t
}

func (t *bitset) ForEachSetRange(start, end uint64, do func(idx uint64)) Bitset {

	for i := start; i <= end && i <= t.max; i++ {

		var found bool

		i, found = t.root.findset(t.rootLevel, start) // TODO: send down 'end' to terminate search sooner

		if !found {
			break // all done
		}

		do(i)
	}

	return t
}

func (t *bitset) ForEachClearRange(start, end uint64, do func(idx uint64)) Bitset {

	for i := start; i <= end && i <= t.max; i++ {

		var found bool

		i, found = t.root.findclr(t.rootLevel, start) // TODO: send down 'end' to terminate search sooner

		if !found {
			break // all done
		}

		do(i)
	}

	return t
}

func (t *bitset) Stats() (numNodes []int) {

	/*
		numNodes = make([]int, len(t.levels))

		for i, l := range t.levels {
			numNodes[i] = l.numNodes
		}
	*/

	return
}

func dbg0(f string, a ...interface{}) {
	fmt.Printf(f, a...)
}

func dbg1(f string, a ...interface{}) {
	fmt.Printf(f, a...)
}

func dbg2(f string, a ...interface{}) {
	fmt.Printf(f, a...)
}

/*
type BitSet interface {
	All() bool
	Any() bool
	BinaryStorageSize() int
	Bytes() []uint64
	* Clear(i uint) *BitSet
	ClearAll() *BitSet
	Clone() *BitSet
	Complement() (result *BitSet)
	Copy(c *BitSet) (count uint)
	Count() uint
	Difference(compare *BitSet) (result *BitSet)
	DifferenceCardinality(compare *BitSet) uint
	DumpAsBits() string
	Equal(c *BitSet) bool
	Flip(i uint) *BitSet
	InPlaceDifference(compare *BitSet)
	InPlaceIntersection(compare *BitSet)
	InPlaceSymmetricDifference(compare *BitSet)
	InPlaceUnion(compare *BitSet)
	Intersection(compare *BitSet) (result *BitSet)
	IntersectionCardinality(compare *BitSet) uint
	IsStrictSuperSet(other *BitSet) bool
	IsSuperSet(other *BitSet) bool
	Len() uint
	MarshalBinary() ([]byte, error)
	MarshalJSON() ([]byte, error)
	NextClear(i uint) (uint, bool)
	NextSet(i uint) (uint, bool)
	None() bool
	ReadFrom(stream io.Reader) (int64, error)
	* Set(i uint) *BitSet
	SetTo(i uint, value bool) *BitSet
	String() string
	SymmetricDifference(compare *BitSet) (result *BitSet)
	SymmetricDifferenceCardinality(compare *BitSet) uint
	* Test(i uint) bool
	Union(compare *BitSet) (result *BitSet)
	UnionCardinality(compare *BitSet) uint
	UnmarshalBinary(data []byte) error
	UnmarshalJSON(data []byte) error
	WriteTo(stream io.Writer) (int64, error)
}
*/
