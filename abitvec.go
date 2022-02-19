package bitvec

import "sync/atomic"

// ABitVec is a bitvector backed by an array of uint64.
// Bits are referenced by bucket (uint64) and bit within
// the bucket.
type ABitVec []uint64

// NewABitVec returns a new bitvector with the given size
func NewABitVec(size uint64) ABitVec {
	return make(ABitVec, (size+mask)>>nbits)
}

func (ABitVec) isBucketBitUnset(bucket uint64, k uint64) bool {
	return bucket&(1<<(k&mask)) == 0
}

func (ABitVec) offset(k uint64) (bucket, bit uint64) {
	return k >> nbits, 1 << (k & mask)
}

// GetBucket returns the int64 bucket
func (bv ABitVec) GetBucket(k uint64) uint64 {
	return atomic.LoadUint64(&bv[k>>nbits])
}

// GetBuckets4 returns buckets 4 at a time.
func (bv ABitVec) GetBuckets4(a, b, c, d uint64) (x, y, z, w uint64) {
	x = atomic.LoadUint64(&bv[a>>nbits])
	y = atomic.LoadUint64(&bv[b>>nbits])
	z = atomic.LoadUint64(&bv[c>>nbits])
	w = atomic.LoadUint64(&bv[d>>nbits])
	return
}

// TrySet will set the bit located at `k` if it is unset
// and will return true if the bit flipped, false otherwise.
func (bv ABitVec) TrySet(k uint64) bool {
	bucket, bit := bv.offset(k)
retry:
	old := atomic.LoadUint64(&bv[bucket])
	if old&bit != 0 {
		return false
	}
	if atomic.CompareAndSwapUint64(&bv[bucket], old, old|bit) {
		return true
	}
	goto retry
}

// TrySetWith performs TrySet but the caller is responsible
// for passing in the old bucket.
func (bv ABitVec) TrySetWith(old uint64, k uint64) bool {
	bucket, bit := bv.offset(k)
	if old&bit != 0 {
		return false
	}
retry:
	if atomic.CompareAndSwapUint64(&bv[bucket], old, old|bit) {
		return true
	}
	old = atomic.LoadUint64(&bv[bucket])
	if old&bit != 0 {
		return false
	}
	goto retry
}
