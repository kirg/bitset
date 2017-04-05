package bitset

type (
	// node is the interface that is implemented by the four types of nodes:
	// - leaf (leaf node, actual storage for bits in []uint64)
	// - inode (intermediate nodes in the tree)
	// - setnode (sparse node that indicates everything under is "set")
	// - clrnode (sparse node that indicates everything under is "clear")
	node interface {
		test(l *level, idx uint64) (set bool)
		set(l *level, idx uint64) (set bool, replace node)
		clr(l *level, idx uint64) (cleared bool, replace node)
		nextset(l *level, start uint64) (idx uint64, found bool)
		nextclr(l *level, start uint64) (idx uint64, found bool)
		prevset(l *level, start uint64) (idx uint64, found bool)
		prevclr(l *level, start uint64) (idx uint64, found bool)
	}
)

func newNode(l *level, sparse, set bool) (n node) {

	if sparse {

		if set {
			return newSparseSet(l)
		}

		return newSparseClr(l)
	}

	// non-sparse node, so count
	l.numNodes++ // #stats

	if l.leaf {

		if set {
			return newLeafSet(l)
		}

		return newLeafClr(l)
	}

	if set {
		return newInodeSet(l)
	}

	return newInodeClr(l)
}

func delNode(l *level, n node) {

	// FIXME: avoid type switch (for perf)?
	switch n.(type) {
	case *setnode:
		return

	case *clrnode:
		return

	case *inode:

		l.numNodes--

		in := n.(*inode)
		for _, nx := range in.nodes {
			delNode(l.next, nx)
		}

	case *leaf:
		l.numNodes--
	}
}

// returns an allset/allclr sparse-node to replace given node
func sparsify(l *level, n node, set bool) (replace node) {

	delNode(l, n)
	return newNode(l, true, set)
}

// returns an allset/allclr 'full' node to replace a sparse node
func desparsify(l *level, n node, set bool) (replace node) {

	delNode(l, n)
	return newNode(l, false, set)
}
