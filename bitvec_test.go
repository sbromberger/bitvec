package bitvec

import (
	// "fmt"
	"math/rand"
	"testing"
)

var GoodSets = []uint64{0, 1, 2, 3, 4, 5, 10, 11, 20, 50, 89, 98, 99}
var UnSets = []uint64{6, 7, 9, 96, 97}
var BadSets = []uint64{100, 101, 200, 1000}

func TestBV(t *testing.T) {
	bv := New(100)
	for _, k := range GoodSets {
		if !bv.TrySet(k) {
			t.Errorf("Error: could not set good index %d", k)
		}
	}
	for _, k := range BadSets {
		if bv.TrySet(k) {
			t.Errorf("Error: successfully set bad index %d ", k)
		}
	}
	for _, k := range GoodSets {
		if x, _ := bv.Get(k); !x {
			t.Errorf("Error: bit %d was supposed to be set but wasn't", k)
		}
	}
	for _, k := range UnSets {
		if x, _ := bv.Get(k); x {
			t.Errorf("Error: bit %d was not supposed to be set but was", k)
		}
	}

	if _, err := bv.Get(101); err == nil {
		t.Errorf("Error should be thrown for out of bounds access")
	}
}

func TestAtomicBV(t *testing.T) {
	bv := NewAtomic(100)
	for _, k := range GoodSets {
		if !bv.TrySet(k) {
			t.Errorf("Error: could not set good index %d", k)
		}
	}
	for _, k := range BadSets {
		if bv.TrySet(k) {
			t.Errorf("Error: successfully set bad index %d ", k)
		}
	}

	for _, k := range GoodSets {
		if x, _ := bv.Get(k); !x {
			t.Errorf("Error: bit %d was supposed to be set but wasn't", k)
		}
	}
	for _, k := range UnSets {
		if x, _ := bv.Get(k); x {
			t.Errorf("Error: bit %d was not supposed to be set but was", k)
		}
	}
	if _, err := bv.Get(101); err == nil {
		t.Errorf("Error should be thrown for out of bounds access")
	}
}

var r = rand.New(rand.NewSource(99))
var sets = r.Perm(10_000_000)

func BenchmarkBVSet(b *testing.B) {
	bv := New(uint64(b.N / 64))
	b.ResetTimer()
	for n := 1; n < b.N; n++ {
		bv.TrySet(uint64(sets[n%len(sets)] % n))
	}
}

func BenchmarkBVAtomicSet(b *testing.B) {
	bv := NewAtomic(uint64(b.N / 64))
	b.ResetTimer()
	for n := 1; n < b.N; n++ {
		bv.TrySet(uint64(sets[n%len(sets)] % n))
	}
}
func BenchmarkSliceSet(b *testing.B) {
	bv := make([]bool, b.N)
	b.ResetTimer()
	for n := 1; n < b.N; n++ {
		bv[sets[n%len(sets)]%n] = true
	}

}

func BenchmarkBVGet(b *testing.B) {
	bv := New(uint64(b.N/64) + 1)
	b.ResetTimer()
	for n := 1; n < b.N; n++ {
		bv.Get(uint64(sets[n%len(sets)] % n))
	}
}
func BenchmarkBVAtomicGet(b *testing.B) {
	bv := NewAtomic(uint64(b.N / 64))
	b.ResetTimer()
	for n := 1; n < b.N; n++ {
		bv.Get(uint64(sets[n%len(sets)%n]))
	}
}
