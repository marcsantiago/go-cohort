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
// such as a metric value, cost, or performance score.
type TrainingAction func() float64

// TrainingSplit controls the logic for assigning users or entities into
// "training" or "normal" flows based on a probabilistic split percentage.
// It optionally remaps cohorts and allows different actions to be executed
// depending on which path is taken.
type TrainingSplit struct {
	cohortRemapper CohortMapper
	onSplit        TrainingAction
	onNormal       TrainingAction
	rng            func() float64 // Random number generator function, expected to return 0.0 <= n < 1.0
	splitPercent   float64        // Probability of assigning to the training cohort (0.0 - 1.0)
	splitType      SplitType      // The type of cohorting strategy used for bucket assignment
}

// NewTrainingSplit creates a new TrainingSplit instance with the provided configuration.
//
// The splitPercent must be between 0.0 and 1.0 inclusive, where:
//   - 0.0: No entities are assigned to training.
//   - 1.0: All entities are assigned to training.
//   - 0.0 < p < 1.0: A random fraction (p%) of entities are assigned to training.
//
// The provided splitAction runs when the entity is assigned to the training path,
// while normalAction runs when it is not.
//
// The random number generator defaults to math/rand/v2.Float64, which is
// concurrency-safe and produces uniformly distributed floats in [0.0, 1.0).
func NewTrainingSplit(splitPercent float64,
	splitType SplitType,
	remapper CohortMapper,
	splitAction TrainingAction,
	normalAction TrainingAction,
) *TrainingSplit {
	if splitPercent < 0 || splitPercent > 1 {
		panic("splitPercent must be between 0.0 and 1.0")
	}
	return &TrainingSplit{
		cohortRemapper: remapper,
		onSplit:        splitAction,
		onNormal:       normalAction,
		rng:            rand.Float64,
		splitPercent:   splitPercent,
		splitType:      splitType,
	}
}

// RunWithTraining executes the appropriate training or normal action
// based on the configured split percentage.
//
// It returns the assigned cohort bucket and the result of the executed action.
//
// Behavior summary:
//   - If splitPercent <= 0: always run normalAction, never split.
//   - If splitPercent >= 1: always run splitAction, always split.
//   - Otherwise: randomly decide to split based on the splitPercent threshold.
//
// The cohort bucket is only assigned (via AssignCohort) when a split occurs.
// Otherwise, NoBucket is returned.
func (t *TrainingSplit) RunWithTraining(targetableID string) (Bucket, float64) {
	var result float64

	// Always take the normal path when splitPercent is 0 or less.
	if t.splitPercent <= 0 {
		if t.onNormal != nil {
			result = t.onNormal()
		}
		return NoBucket, result
	}

	// Determine if the entity should be assigned to training.
	assign := t.splitPercent >= 1 || t.shouldSplit()
	if assign {
		bucket := AssignCohort(targetableID, t.splitType)
		if t.onSplit != nil {
			result = t.onSplit()
		}
		bucket = t.cohortRemapper.Remap(bucket)
		return bucket, result
	}

	// Otherwise, run the normal (non-training) action.
	if t.onNormal != nil {
		result = t.onNormal()
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
