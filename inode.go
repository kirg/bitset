package bitset

import (
	"fmt"
)

type (
	inode struct {
		level  *level // level context
		allSet int    // nodes that are all-set
		allClr int    // nodes that are all-clear
		nodes  []node // (inode/leaf) nodes
	}
)

func newAllSetInode(l *level) *inode {

	nodes := make([]node, l.total)

	for i := range nodes {
		nodes[i] = newAllSetNode(l.next)
	}

	return &inode{
		level:  l,
		allSet: l.total,
		allClr: 0,
		nodes:  nodes,
	}
}

func newAllClrInode(l *level) *inode {

	nodes := make([]node, l.total)

	for i := range nodes {
		nodes[i] = newAllClrNode(l.next)
	}

	return &inode{
		level:  l,
		allSet: 0,
		allClr: l.total,
		nodes:  nodes,
	}
}

func (n *inode) test(idx uint64) bool {

	l := n.level
	i, nextIdx := int(idx>>l.shift), idx&l.mask

	switch nextNode := n.nodes[i]; nextNode {
	case allSetNode:
		return true
	case allClrNode:
		return false
	default:
		return nextNode.test(nextIdx)
	}
}

func (n *inode) set(idx uint64) (set, allset bool) {

	l := n.level
	i, nextIdx := int(idx>>l.shift), idx&l.mask

	switch nextNode := n.nodes[i]; nextNode {

	case allSetNode: // check if already set
		return false, false // no-op

	case allClrNode:

		// de-sparsify node

		l.next.delNode(nextNode)
		nextNode = l.next.newNode(false, false)

		n.nodes[i] = nextNode
		n.allClr--

		fallthrough

	default:

		// recurse down to next level
		set, allset = nextNode.set(nextIdx)

		if !allset {
			return set, false // not all-set
		}

		// sparsify with an all-set node
		l.next.delNode(nextNode)
		n.nodes[i] = l.next.newNode(true, true)

		n.allSet++
		return true, n.allSet == l.total
	}
}

func (n *inode) clr(idx uint64) (cleared, allclear bool) {

	l := n.level
	i, nextIdx := int(idx>>l.shift), idx&l.mask

	switch nextNode := n.nodes[i]; nextNode {

	case allClrNode: // check if already clear
		return false, false // no-op

	case allSetNode:

		// de-sparsify node
		l.next.delNode(nextNode)
		nextNode = l.next.newNode(true, false)

		n.nodes[i] = nextNode
		n.allClr--

		fallthrough

	default:

		// recurse down to next level
		cleared, allclear = nextNode.clr(nextIdx)

		if !allclear {
			return cleared, false // not all-set
		}

		// sparsify with an all-clr node
		l.next.delNode(nextNode)
		n.nodes[i] = l.next.newNode(false, true)

		n.allClr++
		return true, n.allClr == l.total
	}
}

func (n *inode) findset(startIdx uint64) (idx uint64, found bool) {

	l := n.level
	i, nextIdx := int(startIdx>>l.shift), startIdx&l.mask

find:
	for ; i < l.total; i++ {

		switch nextNode := n.nodes[i]; nextNode {

		case allSetNode:
			return (uint64(i) << l.shift) | nextIdx, true

		case allClrNode:
			nextIdx = 0
			continue find

		default:
			nextIdx, found = nextNode.findset(nextIdx)

			if !found {
				nextIdx = 0
				continue find
			}

			return (uint64(i) << l.shift) | nextIdx, true
		}
	}

	return 0, false
}

func (n *inode) findclr(startIdx uint64) (idx uint64, found bool) {

	l := n.level
	i, nextIdx := int(startIdx>>l.shift), startIdx&l.mask

find:
	for ; i < l.total; i++ {

		switch nextNode := n.nodes[i]; nextNode {

		case allClrNode:
			return (uint64(i) << l.shift) | nextIdx, true

		case allSetNode:
			nextIdx = 0
			continue find

		default:
			nextIdx, found = nextNode.findclr(nextIdx)

			if !found {
				nextIdx = 0
				continue find
			}

			return (uint64(i) << l.shift) | nextIdx, true
		}
	}

	return 0, false
}

func (n *inode) String() string {
	return fmt.Sprintf("inode(%p): level=%p allset=%d allclr=%d nodes=%v",
		n, n.level, n.allSet, n.allClr, n.nodes)
}
