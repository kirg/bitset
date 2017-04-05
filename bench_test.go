package bitset_test

import (
	"bitset"
	"fmt"
	"math"
	"math/rand"
	"testing"

	wb "github.com/willf/bitset"
)

var benchConfigs = [][]uint{
	[]uint{32},
	[]uint{16, 16},
	[]uint{14, 12, 6},
	[]uint{8, 8, 8, 8},
	[]uint{10, 8, 6, 4, 4},
	[]uint{4, 4, 4, 4, 4, 4, 4, 4},
}

func BenchmarkBitsetSet(bench *testing.B) {

	for _, cfg := range benchConfigs {

		b := bitset.New(cfg)

		bench.Run(fmt.Sprintf("%v", cfg), func(bench *testing.B) {

			bench.ReportAllocs()

			for i := 0; i < bench.N; i++ {
				b.Set(uint64(rand.Int31()))
			}
		})
	}

	b2 := wb.New(uint(math.MaxInt64))

	bench.Run("willf:32", func(bench *testing.B) {

		bench.ReportAllocs()

		for i := 0; i < bench.N; i++ {
			b2.Set(uint(rand.Int31()))
		}
	})
}
