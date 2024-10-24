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

// _setupH creates an instance of highway hash with a fixed seed, the reason why we need to fix the seed is to ensure
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

// AssignCohort AssignCohorts returns a Bucket which is a string representation of A,B,or C contingent on the split type
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

// AssignCohortABC AssignCohortAB calls AssignCohort and fixes the split type to A/B/C
func AssignCohortABC(identifier string) Bucket {
	return AssignCohort(identifier, SplitCohortABC)
}

// AssignMultipleCohorts generates a bucket that merges the cohort on multiple split types
// e.g. two different tests running on the same user where each test has a different split type assigned
// users see that a blue banner running an A/B test and users that see cats, dogs, or clowns as an A/B/C where
// the users are the same and the tests are running at the same time but the spit is different, test 1 the user is assigned
// bucket A and in test 2 the user is an assigned bucket C, so we return AC as the bucket type. Buckets are always sorted so
// A will represent that status on A/B testing and C on A/B/C testing
func AssignMultipleCohorts(identifier string, splitBy []SplitType) Bucket {
	splitBy = filterDuplicatedBuckets(splitBy)
	buckets := make([]byte, len(splitBy))
	for i, s := range splitBy {
		buckets[i] = AssignCohort(identifier, s)[0]
	}
	return Bucket(byteSlice2String(buckets))
}

func unsafeGetBytes(s string) (b []byte) {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func byteSlice2String(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
