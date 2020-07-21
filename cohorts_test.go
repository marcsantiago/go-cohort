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
		id := uuid.NewV4()
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
		id := uuid.NewV4()
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
		id := uuid.NewV4()
		bucket := AssignCohortAB(id.String())
		bucketTwo := AssignCohortAB(id.String())
		if bucket != bucketTwo {
			t.Fatalf("AssignCohortABSameness() should have been the same got bucketAB: %s bucketTwo: %s id: %s", bucket, bucketTwo, id)
		}

		bucket = AssignCohortABC(id.String())
		bucketTwo = AssignCohortABC(id.String())
		if bucket != bucketTwo {
			t.Fatalf("AssignCohortABCSameness() should have been the same got bucketABC: %s bucketTwo: %s id: %s", bucket, bucketTwo, id)
		}
	}
}
func TestAssignCohortSamenessTwo(t *testing.T) {
	const totalCount = 100000
	id := "32e307b8-add9-4cf7-8716-049244dc9e80"
	for i := 0; i < totalCount; i++ {
		bucket := AssignCohortAB(id)
		bucketTwo := AssignCohortAB(id)

		if bucket != bucketTwo || bucket != BucketB {
			t.Fatalf("TestAssignCohortSamenessTwoAB() should have been the same got bucketAB: %s bucketTwo: %s id: %s, want bucket %s got bucket %s", bucket, bucketTwo, id, bucket, BucketB)
		}

		bucket = AssignCohortABC(id)
		bucketTwo = AssignCohortABC(id)
		if bucket != bucketTwo || bucket != BucketC {
			t.Fatalf("TestAssignCohortSamenessTwoABC() should have been the same got bucketABC: %s bucketTwo: %s id: %s, want bucket %s got bucket %s", bucket, bucketTwo, id, bucket, BucketC)
		}
	}
}

func TestAssignMultipleCohorts(t *testing.T) {
	type args struct {
		id      string
		splitBy []SplitType
	}
	tests := []struct {
		name string
		args args
		want Bucket
	}{
		{
			name: "should create a compound bucket of an A/B test and an A/B/C test",
			args: args{id: "94b58aab-fca2-4050-b0c6-119d7c1f59ca", splitBy: []SplitType{SplitCohortAB, SplitCohortABC}},
			want: Bucket("AC"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AssignMultipleCohorts(tt.args.id, tt.args.splitBy); got != tt.want {
				t.Errorf("DynamicCohort() = %v, want %v", got, tt.want)
			}
		})
	}
}

var TestBucket Bucket

//BenchmarkAssignCohort-12             	 8522774	       138 ns/op	      48 B/op	       1 allocs/op
//BenchmarkAssignCohort-12             	 8827146	       140 ns/op	      48 B/op	       1 allocs/op
//BenchmarkAssignCohort-12             	 8637229	       138 ns/op	      48 B/op	       1 allocs/op
func BenchmarkAssignCohort(b *testing.B) {
	id := uuid.NewV4()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TestBucket = AssignCohort(id.String(), SplitCohortAB)
	}
}

//BenchmarkAssignMultipleCohorts-12    	 4839728	       240 ns/op	      50 B/op	       2 allocs/op
//BenchmarkAssignMultipleCohorts-12    	 4934342	       241 ns/op	      50 B/op	       2 allocs/op
//BenchmarkAssignMultipleCohorts-12    	 4914852	       250 ns/op	      50 B/op	       2 allocs/op
func BenchmarkAssignMultipleCohorts(b *testing.B) {
	id := uuid.NewV4()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TestBucket = AssignMultipleCohorts(id.String(), []SplitType{SplitCohortAB, SplitCohortABC})
	}
}
