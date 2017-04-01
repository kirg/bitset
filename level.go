package bitset

import (
	"fmt"
	"math"
)

type (
	level struct {
		shift uint // shift to compute node index

		leaf  bool   // is leaf node
		total int    // number of child inodes/leaf-bits
		mask  uint64 // mask to compute node index
		next  *level // lower level

		numNodes int // #stats

		height int  // level
		bits   uint // number of bits in 'address'
	}
)

func initLevels(levelBits []uint) (rootLevel *level, maxIdx uint64) {

	// must at least include spec for leaf node
	if len(levelBits) == 0 {
		return nil, 0
	}

	// leaf node should have non-zero allocation
	if levelBits[0] == 0 {
		return nil, 0
	}

	levels := make([]*level, len(levelBits))
	var shift uint
	var next *level

	// compute and initialize level definitions
	for i, n := range levelBits {

		levels[i] = &level{
			shift: shift,
			mask:  (uint64(1) << shift) - 1,
			total: 1 << n,
			next:  next,
			leaf:  i == 0,

			height: i,
			bits:   n,

			numNodes: 0, // #stats
		}

		next = levels[i]
		shift += n
	}

	if shift > 64 {
		return nil, math.MaxInt64
	}

	rootLevel = levels[len(levelBits)-1]
	maxIdx = (uint64(1) << shift) - 1

	return
}

func (t *level) String() string {
	return fmt.Sprintf("level(%p): h=%d bits=%d leaf=%v mask=%d shift=%d next=%p",
		t, t.height, t.bits, t.leaf, t.mask, t.shift, t.next)
}
