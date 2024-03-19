package cohorts

import "sort"

type SplitType int8

// SplitTypeUnknown returns -1 because when using the value in mod math the split will always be 0 and therefore fall into group A without a B group
// this in turn means the split is not running, but also nothing is panicking or breaking.
const SplitTypeUnknown = SplitType(-1)

const (
	SplitCohortAB = iota + 2
	SplitCohortABC
)

// SplitTypes Buckets is an alias for []Bucket with the sort interface implemented
type SplitTypes []SplitType

func (s SplitTypes) Len() int           { return len(s) }
func (s SplitTypes) Less(i, j int) bool { return s[i] < s[j] }
func (s SplitTypes) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func filterDuplicatedBuckets(splitTypes SplitTypes) []SplitType {
	sort.Sort(splitTypes)
	j := 0
	for i := 1; i < len(splitTypes); i++ {
		if splitTypes[j] == splitTypes[i] {
			continue
		}
		j++
		splitTypes[j] = splitTypes[i]
	}
	return splitTypes[:j+1]
}
