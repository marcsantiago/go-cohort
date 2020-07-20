package cohorts

// Bucket provides a string representation of the Bucket the cohort landed in
type Bucket string

const (
	BucketA = "A"
	BucketB = "B"
	BucketC = "C"
)

func toBucket(cohortN uint64) Bucket {
	switch cohortN {
	case 0:
		return BucketA
	case 1:
		return BucketB
	case 2:
		return BucketC
	}
	return ""
}
