package cohorts

import (
	"hash/maphash"
	"reflect"
	"sync"
	"unsafe"
)

var (
	_mu           sync.RWMutex
	_hashFunc     maphash.Hash
	setupHashOnce sync.Once
)

// _setupH creates an instance of maphash.Hash with a fixed seed, the reason why we need to fix the seed is to ensure
// on system restart and on individual instances that id is hashed into the same bucket for deterministic behavior
func _setupH() {
	setupHashOnce.Do(func() {
		// the internal seed value is hidden and has no public setters to be able to do this... soo
		// we are going to have to set it via reflection and unsafe package ¯\_(ツ)_/¯
		// internally they have Seed{s: s1<<32 + s2}
		// where s1 = uint64(runtime_fastrand()) AND s2 = uint64(runtime_fastrand()) ... while  s1|s2 != 0 soo..
		const seedValueOne = uint64(638976000) << 32
		const seedValueTwo = uint64(734572800)
		seedValue := seedValueOne + seedValueTwo
		var seed maphash.Seed
		_setUnexportedField(reflect.ValueOf(&seed).Elem().FieldByName("s"), seedValue)
		_hashFunc.SetSeed(seed)
	})
}

func _setUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}

// AssignCohorts returns a Bucket which is a string representation of A,B,or C contingent on the split type
func AssignCohort(identifier string, splitType SplitType) Bucket {
	_setupH()

	_mu.Lock()
	defer _mu.Unlock()
	_hashFunc.WriteString(identifier)
	cohort := _hashFunc.Sum64() % uint64(splitType)
	_hashFunc.Reset()
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
