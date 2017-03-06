package bitset

import (
	"fmt"
)

const (
	allClearBits = 0
	allSetBits   = (1 << 64) - 1
)

type (
	leaf struct {
		level  *level   // level context
		numSet int      // number of bits that are set
		bits   []uint64 // bit slice of 64-bit integers
	}
)

func newAllSetLeaf(l *level) *leaf {

	bits := make([]uint64, 1+((l.total-1)/64))

	for i := range bits {
		bits[i] = allSetBits
	}

	return &leaf{
		level:  l,
		numSet: l.total,
		bits:   bits,
	}
}

func newAllClrLeaf(l *level) *leaf {

	bits := make([]uint64, 1+((l.total-1)/64))

	// for i := range bits {
	// 	bits[i] = allClearBits
	// }

	return &leaf{
		level:  l,
		numSet: 0,
		bits:   bits,
	}
}

func (n *leaf) test(idx uint64) bool {

	bindex, bmask := int(idx/64), uint64(1)<<(idx%64)
	return (n.bits[bindex] & bmask) != 0
}

func (n *leaf) set(idx uint64) (set, allset bool) {

	bindex, bmask := int(idx/64), uint64(1)<<(idx%64)

	// check if bit already set
	if (n.bits[bindex] & bmask) != 0 {
		return false, false
	}

	n.bits[bindex] |= bmask // set the bit
	n.numSet++

	return true, n.numSet == n.level.total
}

func (n *leaf) clr(idx uint64) (cleared, allclear bool) {

	bindex, bmask := int(idx/64), uint64(1)<<(idx%64)

	// check if bit already clear
	if (n.bits[bindex] & bmask) == 0 {
		return false, false
	}

	n.bits[bindex] &= ^bmask // clear the bit
	n.numSet--

	return true, n.numSet == 0
}

func (n *leaf) findset(startIdx uint64) (idx uint64, found bool) {

	l := n.level

	i := int(startIdx)
	bindex, bmask := (i / 64), uint64(1)<<(uint(i)%64)

find:
	for i < l.total {

		switch n.bits[bindex] {

		case allClearBits:
			// skip over to next uint64
			bindex++
			bmask = 1
			i += 64 - (i % 64)
			continue find

		case allSetBits:
			return uint64(i), true

		default:

			if n.bits[bindex]&bmask != 0 {
				return uint64(i), true
			}

			// move to the next bit
			if i++; i%64 == 0 {
				bindex++
				bmask = 1
			} else {
				bmask <<= 1
			}

			continue find
		}
	}

	return 0, false
}

func (n *leaf) findclr(startIdx uint64) (idx uint64, found bool) {

	l := n.level

	i := int(startIdx)
	bindex, bmask := (i / 64), uint64(1)<<(uint(i)%64)

find:
	for i < l.total {

		switch n.bits[bindex] {

		case allClearBits:
			return uint64(i), true

		case allSetBits:
			// skip over to next uint64
			bindex++
			bmask = 1
			i += 64 - (i % 64)
			continue find

		default:

			if n.bits[bindex]&bmask == 0 {
				return uint64(i), true
			}

			// move to the next bit
			if i++; i%64 == 0 {
				bindex++
				bmask = 1
			} else {
				bmask <<= 1
			}

			continue find
		}
	}

	return 0, false
}

func (n *leaf) String() string {
	return fmt.Sprintf("leaf(%p): level=%p numSet=%d bits=%v",
		n, n.level, n.numSet, n.bits)
}
