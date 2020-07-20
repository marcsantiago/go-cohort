package cohorts

import (
	"fmt"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestSplitByAdvertisingIDDistributionAB(t *testing.T) {
	const (
		totalCount              = 100000
		tolerance               = 0.05
		expectedDistribution    = 0.50
		expectedNumberOfCohorts = 2
	)
	cohorts := make(map[Bucket]int)
	for i := 0; i < totalCount; i++ {
		id, _ := uuid.NewV4()
		bucket := AssignCohort(id.String(), SplitCohortAB)
		cohorts[bucket]++
	}

	if len(cohorts) != expectedNumberOfCohorts {
		fmt.Println(cohorts)
		t.Fatalf("AssignCohortAB() buckets expected %d got %d", expectedNumberOfCohorts, len(cohorts))
	}

	a, b := cohorts[BucketA], cohorts[BucketB]
	if float64(a)/float64(totalCount) > expectedDistribution && float64(a)/float64(totalCount) > expectedDistribution+tolerance {
		t.Fatalf("AssignCohortAB() distribution for A was greater expected %d +/- %.2f got %d", expectedNumberOfCohorts, tolerance, len(cohorts))
	}

	if float64(a)/float64(totalCount) < expectedDistribution && float64(a)/float64(totalCount) < expectedDistribution-tolerance {
		t.Fatalf("AssignCohortAB() distribution for A was less expected %d +/- %.2f got %d", expectedNumberOfCohorts, tolerance, len(cohorts))
	}

	if float64(b)/float64(totalCount) > expectedDistribution && float64(b)/float64(totalCount) > expectedDistribution+tolerance {
		t.Fatalf("AssignCohortAB() distribution for B was greater expected %d +/- %.2f got %d", expectedNumberOfCohorts, tolerance, len(cohorts))
	}

	if float64(b)/float64(totalCount) < expectedDistribution && float64(b)/float64(totalCount) < expectedDistribution-tolerance {
		t.Fatalf("AssignCohortAB() distribution for B was less expected %d +/- %.2f got %d", expectedNumberOfCohorts, tolerance, len(cohorts))
	}

}

func TestSplitByAdvertisingIDDistributionABC(t *testing.T) {
	const (
		totalCount              = 100000
		tolerance               = 0.05
		expectedDistribution    = 0.33
		expectedNumberOfCohorts = 3
	)
	cohorts := make(map[Bucket]int)
	for i := 0; i < totalCount; i++ {
		id, _ := uuid.NewV4()
		bucket := AssignCohort(id.String(), SplitCohortABC)
		cohorts[bucket]++
	}

	if len(cohorts) != expectedNumberOfCohorts {
		t.Fatalf("AssignCohortABC() buckets expected %d got %d", expectedNumberOfCohorts, len(cohorts))
	}

	a, b, c := cohorts[BucketA], cohorts[BucketB], cohorts[BucketC]
	if float64(a)/float64(totalCount) > expectedDistribution && float64(a)/float64(totalCount) > expectedDistribution+tolerance {
		t.Fatalf("AssignCohortABC() distribution for A was greater expected %d +/- %.2f got %d", expectedNumberOfCohorts, tolerance, len(cohorts))
	}

	if float64(a)/float64(totalCount) < expectedDistribution && float64(a)/float64(totalCount) < expectedDistribution-tolerance {
		t.Fatalf("AssignCohortABC() distribution for A was less expected %d +/- %.2f got %d", expectedNumberOfCohorts, tolerance, len(cohorts))
	}

	if float64(b)/float64(totalCount) > expectedDistribution && float64(b)/float64(totalCount) > expectedDistribution+tolerance {
		t.Fatalf("AssignCohortABC() distribution for B was greater expected %d +/- %.2f got %d", expectedNumberOfCohorts, tolerance, len(cohorts))
	}

	if float64(b)/float64(totalCount) < expectedDistribution && float64(b)/float64(totalCount) < expectedDistribution-tolerance {
		t.Fatalf("SplitByAdvertisingIDDistributionAB() distribution for B was less expected %d +/- %.2f got %d", expectedNumberOfCohorts, tolerance, len(cohorts))
	}

	if float64(c)/float64(totalCount) > expectedDistribution && float64(c)/float64(totalCount) > expectedDistribution+tolerance {
		t.Fatalf("AssignCohortABC() distribution for C was greater expected %d +/- %.2f got %d", expectedNumberOfCohorts, tolerance, len(cohorts))
	}

	if float64(c)/float64(totalCount) < expectedDistribution && float64(c)/float64(totalCount) < expectedDistribution-tolerance {
		t.Fatalf("AssignCohortABC() distribution for C was less expected %d +/- %.2f got %d", expectedNumberOfCohorts, tolerance, len(cohorts))
	}
}

func TestAssignCohortSameness(t *testing.T) {
	const totalCount = 100000
	for i := 0; i < totalCount; i++ {
		id, _ := uuid.NewV4()
		bucket := AssignCohortAB(id.String())
		bucketTwo := AssignCohortAB(id.String())
		if bucket != bucketTwo {
			t.Fatalf("AssignCohortABSameness() should have been the same got  bucketAB: %s bucketTwo: %s id: %s", bucket, bucketTwo, id)
		}

		bucket = AssignCohortABC(id.String())
		bucketTwo = AssignCohortABC(id.String())
		if bucket != bucketTwo {
			t.Fatalf("AssignCohortABCSameness() should have been the same got  bucketAB: %s bucketTwo: %s id: %s", bucket, bucketTwo, id)
		}
	}
}

var TestBucket Bucket

//goos: darwin
//goarch: amd64
//pkg: github.com/timehop/nimbus/internal/experimental/cohorts
//BenchmarkSplitByAdvertisingIDABC-12    	11897707	       103 ns/op	      48 B/op	       1 allocs/op
//BenchmarkSplitByAdvertisingIDABC-12    	10915790	       101 ns/op	      48 B/op	       1 allocs/op
//BenchmarkSplitByAdvertisingIDABC-12    	11791810	       103 ns/op	      48 B/op	       1 allocs/op
func BenchmarkSplitByAdvertisingIDABC(b *testing.B) {
	id, _ := uuid.NewV4()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TestBucket = AssignCohort(id.String(), SplitCohortAB)
	}
}
