package bitset

import (
	"fmt"
	"math"
)

type (
	inode struct {
		level      *level // level context
		nSet, nClr int    // nodes that are all-set, all-clr
		nodes      []node // child (inode/leaf) nodes
	}
)

func newInodeSet(l *level) *inode {

	nodes := make([]node, l.total)

	for i := range nodes {
		nodes[i] = newNode(l, true, true)
	}

	return &inode{
		level: l,
		nSet:  l.total,
		nClr:  0,
		nodes: nodes,
	}
}

func newInodeClr(l *level) *inode {

	nodes := make([]node, l.total)

	for i := range nodes {
		nodes[i] = newNode(l, true, false) // optional, since allclr is 'nil'
	}

	return &inode{
		level: l,
		nSet:  0,
		nClr:  l.total,
		nodes: nodes,
	}
}

func (n *inode) test(l *level, idx uint64) bool {

	i, idx := int(idx>>l.shift), idx&l.mask

	if i >= l.total {
		return false
	}

	switch next := n.nodes[i]; next.(type) {
	case *setnode:
		return true

	case *clrnode:
		return false

	default:
		return next.test(l.next, idx)
	}
}

func (in *inode) set(l *level, idx uint64) (set bool, replace node) {

	i, idx := int(idx>>l.shift), idx&l.mask

	next := in.nodes[i]

	set, repl := next.set(l.next, idx)

	if repl == next {
		return set, in
	}

	// assert( set == true ) //

	switch next.(type) {
	case *setnode:
		in.nSet--

	case *clrnode:
		in.nClr--
	}

	in.nodes[i] = repl // replace nextNode

	switch repl.(type) {
	case *setnode:
		if in.nSet++; in.nSet == l.total {
			// sparsify with an all-set node
			return true, sparsify(l, in, true)
		}

	case *clrnode:
		if in.nClr++; in.nClr == l.total {
			// sparsify with an all-clr node
			return true, sparsify(l, in, false)
		}
	}

	return true, in
}

func (in *inode) clr(l *level, idx uint64) (cleared bool, replace node) {

	i, idx := int(idx>>l.shift), idx&l.mask

	next := in.nodes[i]

	cleared, repl := next.clr(l.next, idx)

	if repl == next {
		return cleared, in
	}

	// assert( cleared == true ) //

	switch next.(type) {
	case *setnode:
		in.nSet--

	case *clrnode:
		in.nClr--
	}

	in.nodes[i] = repl // replace nextNode

	switch repl.(type) {
	case *setnode:
		if in.nSet++; in.nSet == l.total {
			// sparsify with an all-set node
			return true, sparsify(l, in, true)
		}

	case *clrnode:
		if in.nClr++; in.nClr == l.total {
			// sparsify with an all-clr node
			return true, sparsify(l, in, false)
		}
	}

	return true, in
}

func (n *inode) nextset(l *level, start uint64) (idx uint64, found bool) {

	i, idx := int(start>>l.shift), start&l.mask

	if i >= l.total {
		return 0, false
	}

	for ; i < l.total; i++ {

		next := n.nodes[i]

		if idx, found = next.nextset(l.next, idx); found {
			return (uint64(i) << l.shift) | idx, true
		}

		idx = 0
	}

	return math.MaxUint64, false
}

func (n *inode) prevset(l *level, start uint64) (idx uint64, found bool) {

	i, idx := int(start>>l.shift), start&l.mask

	if i >= l.total {
		i, idx = l.total-1, math.MaxUint64
	}

	for ; i >= 0; i-- {

		next := n.nodes[i]

		if idx, found = next.prevset(l.next, idx); found {
			return (uint64(i) << l.shift) | idx, true
		}

		idx = math.MaxUint64
	}

	return 0, false
}

func (n *inode) nextclr(l *level, start uint64) (idx uint64, found bool) {

	i, idx := int(start>>l.shift), start&l.mask

	if i >= l.total {
		return 0, false
	}

	for ; i < l.total; i++ {

		next := n.nodes[i]

		if idx, found = next.nextclr(l.next, idx); found {
			return (uint64(i) << l.shift) | idx, true
		}

		idx = 0
	}

	return math.MaxUint64, false
}

func (n *inode) prevclr(l *level, start uint64) (idx uint64, found bool) {

	i, idx := int(start>>l.shift), start&l.mask

	if i >= l.total {
		i, idx = l.total-1, math.MaxUint64
	}

	for ; i >= 0; i-- {

		next := n.nodes[i]

		if idx, found = next.prevclr(l.next, idx); found {
			return (uint64(i) << l.shift) | idx, true
		}

		idx = math.MaxUint64
	}

	return 0, false
}

func (n *inode) String() string {
	return fmt.Sprintf("inode(%p): level=%p nSet=%d nClr=%d nodes=%v",
		n, n.level, n.nSet, n.nClr, n.nodes)
}
