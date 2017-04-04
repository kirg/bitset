package bitset

import (
	"math"
	"unsafe"
)

// setnode defines a sparse node with all bits set
type setnode struct{}

func newSparseSet(l *level) *setnode {
	return (*setnode)(unsafe.Pointer(uintptr(math.MaxUint64)))
}

func newSparseClr(l *level) *clrnode {
	return (*clrnode)(nil)
}

func (sn *setnode) test(l *level, idx uint64) (set bool) {
	return true // set
}

func (sn *setnode) set(l *level, idx uint64) (set bool, replace node) {
	// set on a set-node -> no-op
	return false, sn
}

func (sn *setnode) clr(l *level, idx uint64) (cleared bool, replace node) {
	// desparsify and do 'set' on the new node
	return desparsify(l, sn, true).clr(l, idx)
}

func (sn *setnode) nextset(l *level, start uint64) (idx uint64, found bool) {
	return start, true
}

func (sn *setnode) prevset(l *level, start uint64) (idx uint64, found bool) {

	if start >= uint64(l.total) {
		return uint64(l.total - 1), true
	}

	return start, true
}

func (sn *setnode) nextclr(l *level, start uint64) (idx uint64, found bool) {
	return math.MaxUint64, false
}

func (sn *setnode) prevclr(l *level, start uint64) (idx uint64, found bool) {
	return 0, false
}

// clrnode defines a sparse node with all bits clear
type clrnode struct{}

func (cn *clrnode) test(l *level, idx uint64) (set bool) {
	return false // clear
}

func (cn *clrnode) set(l *level, idx uint64) (set bool, replace node) {
	// desparsify and do 'set' on the new node
	return desparsify(l, cn, false).set(l, idx)
}

func (cn *clrnode) clr(l *level, idx uint64) (cleared bool, replace node) {
	// clear on a clr-node -> no-op
	return false, cn
}

func (cn *clrnode) nextset(l *level, start uint64) (idx uint64, found bool) {
	return math.MaxUint64, false
}

func (cn *clrnode) prevset(l *level, start uint64) (idx uint64, found bool) {
	return 0, false
}

func (cn *clrnode) nextclr(l *level, start uint64) (idx uint64, found bool) {
	return start, true
}

func (cn *clrnode) prevclr(l *level, start uint64) (idx uint64, found bool) {

	if start >= uint64(l.total) {
		return uint64(l.total - 1), true
	}

	return start, true
}
