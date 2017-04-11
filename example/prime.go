package main

import (
	"bitset"
	"fmt"
)

func main() {

	b := bitset.New([]uint{10, 8, 4, 2})

	fmt.Printf("computing primes upto %d:\n", b.Max())

	// start by assuming all numbers are prime, except 0 and 1; then run
	// through all numbers to check if they are a multiple of known primes
	b.SetAll().Clear(0).Clear(1)

	for n := uint64(1); n <= b.Max(); n++ {

		// run through all 'primes' under sqrt(n);
		if b.ForEachSet(func(p uint64) bool {

			// if 'p' greater than sqrt(n), stop
			if p*p > n {
				return false
			}

			// check if n is a multiple of 'p', then 'n' is not prime
			if n%p == 0 {
				b.Clear(n) // not prime
				return false
			}

			return true // continue

		}).Test(n) {

			fmt.Printf("\r%d", n)
		}
	}

	fmt.Printf("\nfound %d primes under %d (stats=%v)\n", b.Count(), b.Max(), b.Stats())
}
