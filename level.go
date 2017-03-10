package bitset

import (
	"fmt"
)

var (
	allSetNode node = &leaf{}
	allClrNode node = nil
)

type (
	node interface {
		test(idx uint64) (set bool)
		set(idx uint64) (set, allset bool)
		clr(idx uint64) (cleared, allclear bool)
		findset(startIdx uint64) (idx uint64, found bool)
		findclr(startIdx uint64) (idx uint64, found bool)
	}

	level struct {
		shift uint // shift to compute node index

		leaf  bool   // is leaf node
		total int    // number of child inodes/leaf-bits
		mask  uint64 // mask to compute node index
		next  *level // lower level

		numNodes, numAllSetNodes, numAllClrNodes int // #stats

		height int  // level
		bits   uint // number of bits in 'address'
	}
)

func initLevels(levelBits []uint) (levels []*level, cap uint64) {

	if len(levelBits) == 0 {
		return nil, 0
	}

	if levelBits[0] == 0 {
		return nil, 0
	}

	levels = make([]*level, len(levelBits))
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

			numNodes:       0, // #stats
			numAllSetNodes: 0, // #stats
			numAllClrNodes: 0, // #stats
		}

		next = levels[i]
		shift += n
	}

	cap = (uint64(1) << shift) - 1

	return
}

func (l *level) newNode(allset bool, sparse bool) (n node) {

	l.numNodes++ // #stats

	if sparse {
		if allset {
			return newAllSetNode(l)
		} else {
			return newAllClrNode(l)
		}
	} else {
		if l.leaf {
			if allset {
				return newAllSetLeaf(l)
			} else {
				return newAllClrLeaf(l)
			}
		} else {
			if allset {
				return newAllSetInode(l)
			} else {
				return newAllClrInode(l)
			}
		}
	}
}

func newAllSetNode(l *level) (n node) {

	l.numAllSetNodes++ // #stats
	return allSetNode
}

func newAllClrNode(l *level) (n node) {

	l.numAllClrNodes++ // #stats
	return allClrNode
}

func (l *level) delNode(n node) {

	switch n { // #stats
	case allSetNode: // #stats
		l.numAllSetNodes-- // #stats
	case allClrNode: // #stats
		l.numAllClrNodes-- // #stats
	default: // #stats
		l.numNodes-- // #stats
		// #stats
		if !l.leaf { // #stats
			in := n.(*inode)             // #stats
			for _, n := range in.nodes { // #stats
				l.next.delNode(n) // #stats
			} // #stats
		} // #stats
	} // #stats
}

func (t *level) String() string {
	return fmt.Sprintf("level(%p): h=%d bits=%d leaf=%v mask=%d shift=%d next=%p",
		t, t.height, t.bits, t.leaf, t.mask, t.shift, t.next)
}
