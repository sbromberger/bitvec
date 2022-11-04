package bitvec

import (
	"math/rand"
	"testing"
)

var GoodSets = []uint64{0, 1, 2, 3, 4, 5, 10, 11, 20, 50, 89, 98, 99}
var UnSets = []uint64{6, 7, 9, 96, 97}
var BadSets = []uint64{100, 101, 200, 1000}

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

// func TestAtomicBV(t *testing.T) {
// 	for _, k := range GoodSets {
// 		if !bv.TrySet(k) {
// 			t.Errorf("Error: could not set good index %d", k)
// 		}
// 	}
// 	for _, k := range BadSets {
// 		if bv.TrySet(k) {
// 			t.Errorf("Error: successfully set bad index %d ", k)
// 		}
// 	}

// 	for _, k := range UnSets {
// 		if x, _ := bv.Get(k); x {
// 			t.Errorf("Error: bit %d was not supposed to be set but was", k)
// 		}
// 	}
// 	if _, err := bv.Get(101); err == nil {
// 		t.Errorf("Error should be thrown for out of bounds access")
// 	}
// }

var r = rand.New(rand.NewSource(99))
var sets = r.Perm(10_000_000)

func BenchmarkSet(b *testing.B) {
	b.Run("bitvec/Set", func(b *testing.B) {
		bv := New(uint64(b.N / 64))
		b.ResetTimer()
		for n := 1; n < b.N; n++ {
			bv.TrySet(uint64(sets[n%len(sets)] % n))
		}
	})
	b.Run("abitvec/Set", func(b *testing.B) {
		abv := NewAtomic(uint64(b.N / 64))
		b.ResetTimer()
		for n := 1; n < b.N; n++ {
			abv.TrySet(uint64(sets[n%len(sets)] % n))
		}
	})
	b.Run("slice/Set", func(b *testing.B) {
		slice := make([]bool, b.N)
		b.ResetTimer()
		for n := 1; n < b.N; n++ {
			slice[sets[n%len(sets)]%n] = true
		}
	})
}

func BenchmarkGet(b *testing.B) {
	b.Run("bitvec/Get", func(b *testing.B) {
		bv := New(uint64(b.N / 64))
		for n := 0; n < b.N && sets[n]%2 == 0; n++ {
			bv.TrySet(uint64(n))
		}
		b.ResetTimer()
		for n := 1; n < b.N; n++ {
			bv.Get(uint64(sets[n%len(sets)] % n))
		}
	})
	b.Run("abitvec/Get", func(b *testing.B) {
		abv := NewAtomic(uint64(b.N / 64))
		for n := 0; n < b.N && sets[n]%2 == 0; n++ {
			abv.TrySet(uint64(n))
		}
		b.ResetTimer()
		for n := 1; n < b.N; n++ {
			abv.Get(uint64(sets[n%len(sets)] % n))
		}
	})
	b.Run("slice/Get", func(b *testing.B) {
		slice := make([]bool, b.N)
		for n := 0; n < b.N && sets[n]%2 == 0; n++ {
			slice[n] = true
		}
		b.ResetTimer()
		for n := 1; n < b.N; n++ {
			_ = slice[sets[n%len(sets)]%n]
		}
	})
}
