package bitset

import (
	"fmt"
	"math"
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

func newLeafSet(l *level) *leaf {

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

func newLeafClr(l *level) *leaf {

	bits := make([]uint64, 1+((l.total-1)/64))

	for i := range bits {
		bits[i] = allClearBits // TODO: no-op
	}

	return &leaf{
		level:  l,
		numSet: 0,
		bits:   bits,
	}
}

func (n *leaf) test(l *level, idx uint64) bool {

	bindex, bmask := int(idx/64), uint64(1)<<(idx%64)
	return (n.bits[bindex] & bmask) != 0
}

func (n *leaf) set(l *level, idx uint64) (set bool, replace node) {

	bindex, bmask := int(idx/64), uint64(1)<<(idx%64)

	// check if bit already set
	if (n.bits[bindex] & bmask) != 0 {
		return false, n
	}

	n.bits[bindex] |= bmask // set the bit

	if n.numSet++; n.numSet == l.total {
		return true, sparsify(l, n, true)
	}

	return true, n
}

func (n *leaf) clr(l *level, idx uint64) (cleared bool, replace node) {

	bindex, bmask := int(idx/64), uint64(1)<<(idx%64)

	// check if bit already clear
	if (n.bits[bindex] & bmask) == 0 {
		return false, n
	}

	n.bits[bindex] &= ^bmask // clear the bit

	if n.numSet--; n.numSet == 0 {
		return true, sparsify(l, n, false)
	}

	return true, n
}

func (n *leaf) nextset(l *level, start uint64) (idx uint64, found bool) {

	i := int(start)
	bindex, bmask := (i / 64), uint64(1)<<(uint(i)%64)

find:
	for i < l.total {

		switch n.bits[bindex] {

		case allSetBits:
			return uint64(i), true

		case allClearBits:
			// skip over to next uint64
			bindex++
			bmask = 1
			i += 64 - (i % 64)
			continue find

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

	return math.MaxUint64, false
}

func (n *leaf) prevset(l *level, start uint64) (idx uint64, found bool) {

	if start >= uint64(l.total) {
		start = uint64(l.total - 1)
	}

	i := int(start)

	bindex, bmask := (i / 64), uint64(1)<<(uint(i)%64)

find:
	for i >= 0 {

		switch n.bits[bindex] {

		case allSetBits:
			return uint64(i), true

		case allClearBits:
			// skip over to prev uint64
			bindex--
			bmask = 1 << 63
			i -= 1 + (i % 64)

			continue find

		default:
			if n.bits[bindex]&bmask != 0 {
				return uint64(i), true
			}

			// move to the prev bit
			if i--; i%64 == 63 {
				bindex--
				bmask = 1 << 63
			} else {
				bmask >>= 1
			}

			continue find
		}
	}

	return 0, false
}

func (n *leaf) nextclr(l *level, start uint64) (idx uint64, found bool) {

	i := int(start)
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

	return math.MaxUint64, false
}

func (n *leaf) prevclr(l *level, start uint64) (idx uint64, found bool) {

	if start >= uint64(l.total) {
		start = uint64(l.total - 1)
	}

	i := int(start)

	bindex, bmask := (i / 64), uint64(1)<<(uint(i)%64)

find:
	for i >= 0 {

		switch n.bits[bindex] {

		case allClearBits:
			return uint64(i), true

		case allSetBits:
			// skip over to next uint64
			bindex--
			bmask = 1 << 63
			i -= 1 + (i % 64)

			continue find

		default:

			if n.bits[bindex]&bmask == 0 {
				return uint64(i), true
			}

			// move to the next bit
			if i--; i%64 == 63 {
				bindex--
				bmask = 1 << 63
			} else {
				bmask >>= 1
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
