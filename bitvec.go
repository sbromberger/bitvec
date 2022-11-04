// Package bitvec is bit-vector with atomic and non-atomic access
package bitvec

import "errors"

const (
	nbits = 6          // 64 bits in a uint64
	ws    = 1 << nbits // constant 64
	mask  = ws - 1     // all ones
)

func offset(k uint64) (bucket, bit uint64) {
	return k >> nbits, 1 << (k & mask)
}

// BitVec is a nonatomic bit vector.
type BitVec struct {
	buckets  []uint64
	capacity uint64
}

// NewBitVec creates a non-atomic bitvector.
func New(size uint64) BitVec {
	nints := size / ws
	if size-(nints*ws) != 0 {
		nints++
	}

	return BitVec{make([]uint64, nints), size}
}

// TrySet will try to set the bit and will return true if set
// is successful, false if bit is already set.
func (bv BitVec) TrySet(k uint64) bool {
	if k >= bv.capacity {
		return false
	}
	bucket, bit := offset(k)
	old := bv.buckets[bucket]
	if old&bit != 0 {
		return false
	}
	bv.buckets[bucket] = old | bit
	return true
}

// Get will return true if the bit is set; false otherwise.
func (bv BitVec) Get(k uint64) (bool, error) {
	if k >= bv.capacity {
		return false, errors.New("Attempt to access element beyond vector bounds")
	}
	bucket, bit := offset(k)
	return bv.buckets[bucket]&bit != 0, nil
}
