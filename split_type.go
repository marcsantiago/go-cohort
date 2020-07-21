package cohorts

type SplitType int8

// returns -1 because when using the value in mod math the split will always be 0 and therefore fall into group A without a B group
// this in turn means the split is not running, but also nothing is panicking or breaking.
const SplitTypeUnknown = SplitType(-1)

const (
	SplitCohortAB = iota + 2
	SplitCohortABC
)
