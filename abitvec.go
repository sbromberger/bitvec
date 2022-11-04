package bitvec

import (
	"errors"
	"sync/atomic"
)

// ABitVec is a bitvector backed by an array of uint64.
// Bits are referenced by bucket (uint64) and bit within
// the bucket.
type ABitVec struct {
	buckets  []uint64
	capacity uint64
}

// NewAtomic returns a new bitvector with the given size
func NewAtomic(size uint64) ABitVec {
	return ABitVec{make([]uint64, (size+mask)>>nbits), size}
}

func isBucketBitUnset(bucket uint64, k uint64) bool {
	return bucket&(1<<(k&mask)) == 0
}

// GetBucket returns the int64 bucket
func (bv ABitVec) GetBucket(k uint64) uint64 {
	return atomic.LoadUint64(&bv.buckets[k>>nbits])
}

// GetBuckets4 returns buckets 4 at a time.
func (bv ABitVec) GetBuckets4(a, b, c, d uint64) (x, y, z, w uint64) {
	x = atomic.LoadUint64(&bv.buckets[a>>nbits])
	y = atomic.LoadUint64(&bv.buckets[b>>nbits])
	z = atomic.LoadUint64(&bv.buckets[c>>nbits])
	w = atomic.LoadUint64(&bv.buckets[d>>nbits])
	return
}

// TrySet will set the bit located at `k` if it is unset
// and will return true if the bit flipped, false otherwise.
func (bv ABitVec) TrySet(k uint64) bool {
	if k >= bv.capacity {
		return false
	}
	bucket, bit := offset(k)
retry:
	old := atomic.LoadUint64(&bv.buckets[bucket])
	if old&bit != 0 {
		return false
	}
	if atomic.CompareAndSwapUint64(&bv.buckets[bucket], old, old|bit) {
		return true
	}
	goto retry
}

// TrySetWith performs TrySet but the caller is responsible
// for passing in the old bucket.
func (bv ABitVec) TrySetWith(old uint64, k uint64) bool {
	if k >= bv.capacity {
		return false
	}
	bucket, bit := offset(k)
	if old&bit != 0 {
		return false
	}
retry:
	if atomic.CompareAndSwapUint64(&bv.buckets[bucket], old, old|bit) {
		return true
	}
	old = atomic.LoadUint64(&bv.buckets[bucket])
	if old&bit != 0 {
		return false
	}
	goto retry
}

// Get will return true if the bit is set; false otherwise.
func (bv ABitVec) Get(k uint64) (bool, error) {
	if k >= bv.capacity {
		return false, errors.New("Attempt to access element beyond vector bounds")
	}
	bucket, bit := offset(k)
	return bv.buckets[bucket]&bit != 0, nil
}
