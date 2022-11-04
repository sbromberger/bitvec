package bitvec

import (
	"fmt"
	"math/rand"
	"testing"
)

var GoodSets = []uint64{0, 1, 2, 3, 4, 5, 10, 11, 20, 50, 89, 98, 99}
var UnSets = []uint64{6, 7, 9, 96, 97}
var BadSets = []uint64{100, 101, 200, 1000}

var setSizes = []int{13, 256, 65536, 16777216}

func TestBV(t *testing.T) {
	bv := New(100)
	abv := NewAtomic(100)
	t.Run("valid indices", func(t *testing.T) {
		for _, k := range GoodSets {
			if !bv.TrySet(k) {
				t.Errorf("bv could not set good index %d", k)
			}
			if !abv.TrySet(k) {
				t.Errorf("abv could not set good index %d", k)
			}

			if x, _ := bv.Get(k); !x {
				t.Errorf("bv bit %d was supposed to be set but wasn't", k)
			}
			if x, _ := abv.Get(k); !x {
				t.Errorf("abv bit %d was supposed to be set but wasn't", k)
			}
		}
	})
	t.Run("invalid indices", func(t *testing.T) {
		for _, k := range BadSets {
			if bv.TrySet(k) {
				t.Errorf("bv successfully set bad index %d ", k)
			}
			if abv.TrySet(k) {
				t.Errorf("abv successfully set bad index %d ", k)
			}
		}
		if _, err := bv.Get(101); err == nil {
			t.Errorf("bv error should be thrown for out of bounds access")
		}
		if _, err := abv.Get(101); err == nil {
			t.Errorf("abv error should be thrown for out of bounds access")
		}

	})
	t.Run("unset indices", func(t *testing.T) {
		for _, k := range UnSets {
			if x, _ := bv.Get(k); x {
				t.Errorf("bv bit %d was not supposed to be set but was", k)
			}
			if x, _ := abv.Get(k); x {
				t.Errorf("abv bit %d was not supposed to be set but was", k)
			}
		}
	})

}

var r = rand.New(rand.NewSource(99))
var sets = r.Perm(1024 * 1024)

func BenchmarkSet(b *testing.B) {
	for _, setSize := range setSizes {
		b.Run(fmt.Sprintf("bitvec/Set size %d", setSize), func(b *testing.B) {
			bv := New(uint64(setSize))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bv.TrySet(uint64(sets[i%len(sets)] % setSize))
			}
		})
		b.Run(fmt.Sprintf("abitvec/Set size %d", setSize), func(b *testing.B) {
			abv := NewAtomic(uint64(setSize))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				abv.TrySet(uint64(sets[i%len(sets)] % setSize))
			}
		})
		b.Run(fmt.Sprintf("slice/Set size %d", setSize), func(b *testing.B) {
			slice := make([]bool, setSize)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				slice[sets[i%len(sets)]%setSize] = true
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	for _, setSize := range setSizes {
		b.Run(fmt.Sprintf("bitvec/Get size %d", setSize), func(b *testing.B) {
			bv := New(uint64(setSize))
			for n := 0; n < setSize && sets[n]%2 == 0; n++ {
				bv.TrySet(uint64(n))
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bv.Get(uint64(sets[i%len(sets)] % setSize))
			}
		})
		b.Run(fmt.Sprintf("abitvec/Get size %d", setSize), func(b *testing.B) {
			abv := NewAtomic(uint64(setSize))
			for n := 0; n < setSize && sets[n]%2 == 0; n++ {
				abv.TrySet(uint64(n))
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				abv.Get(uint64(sets[i%len(sets)] % setSize))
			}
		})
		b.Run(fmt.Sprintf("slice/Get size %d", setSize), func(b *testing.B) {
			slice := make([]bool, setSize)
			for n := 0; n < setSize && sets[n]%2 == 0; n++ {
				slice[n] = true
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = slice[sets[i%len(sets)]%setSize]
			}
		})
	}
}
