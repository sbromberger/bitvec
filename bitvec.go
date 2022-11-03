// Package bitvec is bit-vector with atomic and non-atomic access
package bitvec

const (
	nbits   = 6          // 64 bits in a uint64
	ws      = 1 << nbits // constant 64
	mask    = ws - 1     // all ones
	bitsize = 2 ^ nbits
)

// BitVec is a nonatomic bit vector.
type BitVec []uint64

// New creates a non-atomic bitvector with a given size.
func New(size uint64) BitVec {
	nints := size / ws
	if size-(nints*bitsize) != 0 {
		nints++
	}

	return make(BitVec, nints)
}

func (BitVec) offset(k uint64) (bucket, bit uint64) {
	return k >> nbits, 1 << (k & mask)
}

// TrySet will try to set the bit and will return true if set
// is successful.
func (bv BitVec) TrySet(k uint64) bool {
	bucket, bit := bv.offset(k)
	old := bv[bucket]
	if old&bit != 0 {
		return false
	}
	bv[bucket] = old | bit
	return true
}
