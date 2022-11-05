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
	bv2 := New(100)
	abv := NewAtomic(100)
	abv2 := NewAtomic(100)
	t.Run("valid indices", func(t *testing.T) {
		for _, k := range GoodSets {
			if !bv.TrySet(k) {
				t.Errorf("bv could not set good index %d", k)
			}
			if !abv.TrySet(k) {
				t.Errorf("abv could not set good index %d", k)
			}
			bv2.Set(k)
			abv2.Set(k)

			if x, _ := bv.Get(k); !x {
				t.Errorf("bv bit %d was supposed to be set but wasn't", k)
			}
			if x, _ := abv.Get(k); !x {
				t.Errorf("abv bit %d was supposed to be set but wasn't", k)
			}
			if x, _ := bv2.Get(k); !x {
				t.Errorf("bv2 bit %d was supposed to be set but wasn't", k)
			}
			if x, _ := abv2.Get(k); !x {
				t.Errorf("abv2 bit %d was supposed to be set but wasn't", k)
			}

			bv.Clear(k)
			if x, _ := bv.Get(k); x {
				t.Errorf("bv bit %d was supposed to be cleared but wasn't", k)
			}
			abv.Clear(k)
			if x, _ := abv.Get(k); x {
				t.Errorf("abv bit %d was supposed to be cleared but wasn't", k)
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
			if bv2.Set(k) == nil {
				t.Errorf("bv2 successfully set bad index %d ", k)
			}
			if abv2.Set(k) == nil {
				t.Errorf("abv2 successfully set bad index %d ", k)
			}
			if _, err := bv.Get(k); err == nil {
				t.Errorf("bv error should be thrown for out of bounds access")
			}
			if _, err := abv.Get(k); err == nil {
				t.Errorf("abv error should be thrown for out of bounds access")
			}
			if bv.Clear(k) == nil {
				t.Errorf("bv error should be thrown for out of bounds access")
			}
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
var sets = r.Perm(setSizes[len(setSizes)-1])

func BenchmarkSet(b *testing.B) {
	for _, setSize := range setSizes {
		b.Run(fmt.Sprintf("bitvec/Set size %d", setSize), func(b *testing.B) {
			bv := New(uint64(setSize))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bv.Set(uint64(sets[i%len(sets)] % setSize))
			}
		})
		b.Run(fmt.Sprintf("abitvec/Set size %d", setSize), func(b *testing.B) {
			abv := NewAtomic(uint64(setSize))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				abv.Set(uint64(sets[i%len(sets)] % setSize))
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
	fmt.Println()
}

func BenchmarkTrySet(b *testing.B) {
	for _, setSize := range setSizes {
		b.Run(fmt.Sprintf("bitvec/TrySet size %d", setSize), func(b *testing.B) {
			bv := New(uint64(setSize))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bv.TrySet(uint64(sets[i%len(sets)] % setSize))
			}
		})
		b.Run(fmt.Sprintf("abitvec/TrySet size %d", setSize), func(b *testing.B) {
			abv := NewAtomic(uint64(setSize))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				abv.TrySet(uint64(sets[i%len(sets)] % setSize))
			}
		})
	}
	fmt.Println()
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
	fmt.Println()
}
func BenchmarkClear(b *testing.B) {
	for _, setSize := range setSizes {
		b.Run(fmt.Sprintf("bitvec/Clear size %d", setSize), func(b *testing.B) {
			bv := New(uint64(setSize))
			for i := 0; i < b.N; i++ {
				bv.Set(uint64(sets[i%len(sets)] % setSize))
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bv.Clear(uint64(sets[i%len(sets)] % setSize))
			}
		})
		b.Run(fmt.Sprintf("abitvec/Clear size %d", setSize), func(b *testing.B) {
			abv := NewAtomic(uint64(setSize))
			for i := 0; i < b.N; i++ {
				abv.Set(uint64(sets[i%len(sets)] % setSize))
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				abv.Clear(uint64(sets[i%len(sets)] % setSize))
			}
		})
		b.Run(fmt.Sprintf("slice/Clear size %d", setSize), func(b *testing.B) {
			slice := make([]bool, setSize)
			for i := 0; i < b.N; i++ {
				slice[sets[i%len(sets)]%setSize] = true
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				slice[sets[i%len(sets)]%setSize] = false
			}
		})
	}
	fmt.Println()
}
