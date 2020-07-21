package cohorts

import (
	"hash"
	"math/rand"
	"sync"
	"unsafe"

	"github.com/minio/highwayhash"
)

var (
	_mu           sync.RWMutex
	_hashFunc     hash.Hash64
	setupHashOnce sync.Once
)

// _setupH creates an instance of highwayhash with a fixed seed, the reason why we need to fix the seed is to ensure
// on system restart and on individual instances that id is hashed into the same bucket for deterministic behavior
func _setupH() {
	setupHashOnce.Do(func() {
		const seedValueOne = uint64(638976000) << 32
		const seedValueTwo = uint64(734572800)

		seedValue := seedValueOne + seedValueTwo
		key := make([]byte, 32)
		s := rand.NewSource(int64(seedValue))
		r := rand.New(s)
		_, _ = r.Read(key)

		hh, err := highwayhash.New64(key)
		if err != nil {
			panic(err)
		}
		_hashFunc = hh
	})
}

// AssignCohorts returns a Bucket which is a string representation of A,B,or C contingent on the split type
func AssignCohort(identifier string, splitType SplitType) Bucket {
	_setupH()

	_mu.Lock()
	_hashFunc.Write(unsafeGetBytes(identifier))
	cohort := _hashFunc.Sum64() % uint64(splitType)
	_hashFunc.Reset()
	_mu.Unlock()
	return toBucket(cohort)
}

// AssignCohortAB calls AssignCohort and fixes the split type to A/B
func AssignCohortAB(identifier string) Bucket {
	return AssignCohort(identifier, SplitCohortAB)
}

// AssignCohortAB calls AssignCohort and fixes the split type to A/B/C
func AssignCohortABC(identifier string) Bucket {
	return AssignCohort(identifier, SplitCohortABC)
}

// AssignMultipleCohorts generates a bucket that sums the cohort on multiple split types
// e.g two different tests running on the same user where each test has a different split type assigned
// e.g users see that a blue banner running an A/B test and users that see cats, dogs, or clowns as an A/B/C where
// the users are the same and the tests are running at the same time
func AssignMultipleCohorts(identifier string, splitBy []SplitType) Bucket {
	buckets := make([]byte, len(splitBy))
	for i, s := range splitBy {
		buckets[i] = AssignCohort(identifier, s)[0]
	}
	return Bucket(byteSlice2String(buckets))
}

func unsafeGetBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func byteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
