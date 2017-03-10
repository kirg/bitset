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
		Count() uint64

		FindSet(start uint64) (idx uint64, found bool)
		FindClear(start uint64) (idx uint64, found bool)

		// NextSet(idx uint64) (uint, bool)
		// NextClear(idx uint64) (uint, bool)

		// SetBits(startIdx, endIdx uint64)
		// ClearBits(startIdx, endIdx uint64)

		Cap() uint64
		Stats() []int // #stats
	}
)

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

type (
	bitset struct {
		levels []*level
		root   node
		count  uint64
		cap    uint64
	}
)

func New(levelBits []uint) Bitset {

	levels, cap := initLevels(levelBits)

	return &bitset{
		root:   levels[len(levels)-1].newNode(false, false), // node at max level
		cap:    cap,
		levels: levels,
	}
}

func (t *bitset) Cap() uint64 {
	return t.cap
}

func (t *bitset) Test(idx uint64) (ret bool) {
	return t.root.test(idx)
}

func (t *bitset) Set(idx uint64) Bitset {

	if idx > t.cap {
		return nil
	}

	if set, _ := t.root.set(idx); set {
		t.count++
	}

	return t
}

func (t *bitset) Clear(idx uint64) Bitset {

	if idx > t.cap {
		return nil
	}

	if cleared, _ := t.root.clr(idx); cleared {
		t.count--
	}

	return t
}

func (t *bitset) Count() uint64 {

	return t.count
}

func (t *bitset) FindSet(start uint64) (idx uint64, found bool) {

	return t.root.findset(start)
}

func (t *bitset) FindClear(start uint64) (idx uint64, found bool) {

	return t.root.findclr(start)
}

func (t *bitset) Stats() (numNodes []int) {

	numNodes = make([]int, len(t.levels))

	for i, l := range t.levels {
		numNodes[i] = l.numNodes
	}

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
