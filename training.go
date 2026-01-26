package cohorts

import (
	"math/rand/v2"
)

const NoBucket = Bucket("")

// CohortMapper defines a mapping between one cohort bucket and another.
// This can be used to remap users from one experimental cohort to another
// (for example, to align training or testing buckets).
type CohortMapper map[Bucket]Bucket

// Remap returns the remapped bucket if one exists in the mapper.
// If the mapper is nil or the bucket is not found, the original bucket is returned unchanged.
func (cm CohortMapper) Remap(bucket Bucket) Bucket {
	if cm == nil {
		return bucket
	}
	if remapped, ok := cm[bucket]; ok {
		return remapped
	}
	return bucket
}

// TrainingAction defines the signature for training callbacks.
// It should perform an action and return a float64 representing a measurement,
// such as a metric value, cost, or performance score. Actions are intentionally
// parameterless so callers typically provide per-call parameters by using
// closures that capture the needed values.
type TrainingAction func() float64

// TrainingSplit controls the logic for assigning users or entities into
// "training" or "normal" flows based on a probabilistic split percentage.
// The struct holds configuration used by the assignment logic; concrete
// actions are provided per-call to RunWithTraining (so the same TrainingSplit
// can be reused concurrently with different action closures).
type TrainingSplit struct {
	cohortRemapper CohortMapper
	rng            func() float64 // Random number generator function, expected to return 0.0 <= n < 1.0
	splitPercent   float64        // Probability of assigning to the training cohort (0.0 - 1.0)
	splitType      SplitType      // The type of cohorting strategy used for bucket assignment
}

// Option applies optional configuration to a TrainingSplit.
type Option func(*TrainingSplit)

// NewTrainingSplit creates a new TrainingSplit instance with the provided configuration.
//
// The splitPercent must be between 0.0 and 1.0 inclusive, where:
//   - 0.0: No entities are assigned to training.
//   - 1.0: All entities are assigned to training.
//   - 0.0 < p < 1.0: A random fraction (p) of entities are assigned to training.
//
// The returned TrainingSplit does not contain any per-call actions. Callers must
// provide the split and normal actions when invoking RunWithTraining. The RNG
// defaults to math/rand/v2.Float64 which returns uniformly distributed floats
// in [0.0, 1.0).
func NewTrainingSplit(splitPercent float64,
	splitType SplitType,
	remapper CohortMapper,
	options ...Option,
) *TrainingSplit {
	if splitPercent < 0 || splitPercent > 1 {
		panic("splitPercent must be between 0.0 and 1.0")
	}
	t := &TrainingSplit{
		cohortRemapper: remapper,
		rng:            rand.Float64,
		splitPercent:   splitPercent,
		splitType:      splitType,
	}
	for _, option := range options {
		option(t)
	}
	return t
}

// WithCustomRng overrides the default random number generator used by
// TrainingSplit when deciding whether to assign an entity to the training cohort.
//
// The provided rng function must return a float64 in the range [0.0, 1.0).
// This is useful for testing (e.g., deterministic behavior) or for injecting
// a custom source of randomness.
func WithCustomRng(rng func() float64) Option {
	return func(t *TrainingSplit) {
		t.rng = rng
	}
}

// RunWithTraining executes the appropriate training or normal action
// based on the configured split percentage. The caller provides the actions
// to run for the normal (non-training) path and the split (training) path.
//
// Parameters:
//   - targetableID: used to assign a cohort via AssignCohort when a split occurs.
//   - normalAction: action to execute when the entity is not assigned to training.
//     May be nil; if so, no action is run for the normal path.
//   - splitAction: action to execute when the entity is assigned to training.
//     May be nil; if so, no action is run for the split path.
//
// Returns:
//   - the assigned Bucket (the remapped value when split occurs), or NoBucket if not split.
//   - the float64 result returned by the executed action (0 if no action ran).
//
// Behavior summary:
//   - If splitPercent <= 0: always run normalAction, never split.
//   - If splitPercent >= 1: always run splitAction, always split.
//   - Otherwise: randomly decide to split based on the splitPercent threshold.
//
// Note: callers typically pass closures for normalAction and splitAction that
// capture per-call parameters (e.g., request-scoped values) so the TrainingSplit
// instance itself can remain reusable and concurrency-safe.
func (t *TrainingSplit) RunWithTraining(targetableID string, normalAction TrainingAction, splitAction TrainingAction) (Bucket, float64) {
	var result float64

	// Always take the normal path when splitPercent is 0 or less.
	if t.splitPercent <= 0 {
		if normalAction != nil {
			result = normalAction()
		}
		return NoBucket, result
	}

	// Determine if the entity should be assigned to training.
	assign := t.splitPercent >= 1 || t.shouldSplit()
	if assign {
		bucket := AssignCohort(targetableID, t.splitType)
		if splitAction != nil {
			result = splitAction()
		}
		bucket = t.cohortRemapper.Remap(bucket)
		return bucket, result
	}

	// Otherwise, run the normal (non-training) action.
	if normalAction != nil {
		result = normalAction()
	}
	return NoBucket, result
}

// shouldSplit decides whether to assign the current entity to the training cohort.
// It generates a random float in [0.0, 1.0) and returns true if the value is below splitPercent.
//
// Example:
//
//	splitPercent = 0.25 → ~25% chance to return true.
func (t *TrainingSplit) shouldSplit() bool {
	// 0.0 <= n < 1.0
	return t.rng() < t.splitPercent
}
