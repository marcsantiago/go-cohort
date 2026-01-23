package cohorts

import (
	"math/rand/v2"
	"testing"
)

func TestCohortMapperRemap(t *testing.T) {
	type fields struct {
		mapper CohortMapper
	}
	tests := []struct {
		name   string
		fields fields
		args   struct {
			bucket Bucket
		}
		want Bucket
	}{
		{
			name: "nil mapper returns original",
			fields: fields{
				mapper: nil,
			},
			args: struct{ bucket Bucket }{
				bucket: Bucket("orig"),
			},
			want: Bucket("orig"),
		},
		{
			name: "mapped bucket returns remapped value",
			fields: fields{
				mapper: CohortMapper{
					Bucket("a"): Bucket("b"),
				},
			},
			args: struct{ bucket Bucket }{
				bucket: Bucket("a"),
			},
			want: Bucket("b"),
		},
		{
			name: "unknown bucket returns original",
			fields: fields{
				mapper: CohortMapper{
					Bucket("a"): Bucket("b"),
				},
			},
			args: struct{ bucket Bucket }{
				bucket: Bucket("unknown"),
			},
			want: Bucket("unknown"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.mapper.Remap(tt.args.bucket)
			if got != tt.want {
				t.Fatalf("Remap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTrainingSplit_Panic(t *testing.T) {
	type args struct {
		splitPercent float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "less than zero",
			args: args{splitPercent: -0.1},
		},
		{
			name: "greater than one",
			args: args{splitPercent: 1.1},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("NewTrainingSplit did not panic for splitPercent=%v", tt.args.splitPercent)
				}
			}()

			// call constructor expected to panic
			_ = NewTrainingSplit(tt.args.splitPercent, SplitType(0), nil)
		})
	}
}

func TestRunWithTraining(t *testing.T) {
	type fields struct {
		cohortRemapper CohortMapper
		normalAction   TrainingAction
		splitAction    TrainingAction
		rng            func() float64
		splitPercent   float64
		splitType      SplitType
	}
	type args struct {
		targetableID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   struct {
			bucket Bucket
			value  float64
			split  bool // whether we expect an assigned bucket
		}
	}{
		{
			name: "splitPercent 0 runs normal action and returns unspecified bucket",
			fields: fields{
				cohortRemapper: nil,
				splitAction: func() float64 {
					return 123.0
				},
				normalAction: func() float64 {
					return 3.14
				},
				rng:          func() float64 { return 0.0 },
				splitPercent: 0.0,
				splitType:    SplitCohortAB,
			},
			args: args{targetableID: "id-0"},
			want: struct {
				bucket Bucket
				value  float64
				split  bool
			}{
				bucket: NoBucket,
				value:  3.14,
				split:  false,
			},
		},
		{
			name: "splitPercent 1 runs split action and returns an assigned bucket",
			fields: fields{
				cohortRemapper: nil,
				splitAction: func() float64 {
					return 2.71
				},
				normalAction: func() float64 { return 0.0 },
				rng:          nil,
				splitPercent: 1.0,
				splitType:    SplitCohortAB,
			},
			args: args{targetableID: "id-1"},
			want: struct {
				bucket Bucket
				value  float64
				split  bool
			}{
				bucket: Bucket(""),
				value:  2.71,
				split:  true,
			},
		},
		{
			name: "intermediate split with rng forcing split (rng < splitPercent)",
			fields: fields{
				cohortRemapper: nil,
				splitAction: func() float64 {
					return 7.7
				},
				normalAction: func() float64 { return 0.0 },
				rng:          func() float64 { return 0.0 },
				splitPercent: 0.5,
				splitType:    SplitCohortAB,
			},
			args: args{targetableID: "id-2"},
			want: struct {
				bucket Bucket
				value  float64
				split  bool
			}{
				bucket: Bucket(""),
				value:  7.7,
				split:  true,
			},
		},
		{
			name: "intermediate split with rng avoiding split (rng >= splitPercent)",
			fields: fields{
				cohortRemapper: nil,
				splitAction:    func() float64 { return 9.9 },
				normalAction:   func() float64 { return 6.6 },
				rng:            func() float64 { return 0.99 },
				splitPercent:   0.5,
				splitType:      SplitCohortAB,
			},
			args: args{targetableID: "id-3"},
			want: struct {
				bucket Bucket
				value  float64
				split  bool
			}{
				bucket: NoBucket,
				value:  6.6,
				split:  false,
			},
		},
		{
			name: "cohort remapper is applied to assigned bucket",
			fields: fields{
				cohortRemapper: CohortMapper{
					Bucket("B"): Bucket("remapped"),
				},
				splitAction: func() float64 {
					return 1.23
				},
				normalAction: func() float64 { return 0.0 },
				rng:          func() float64 { return 0.0 },
				splitPercent: 0.5,
				splitType:    SplitCohortAB,
			},
			args: args{targetableID: "id-4"},
			want: struct {
				bucket Bucket
				value  float64
				split  bool
			}{
				bucket: Bucket("remapped"),
				value:  1.23,
				split:  true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ts := &TrainingSplit{
				cohortRemapper: tt.fields.cohortRemapper,
				rng:            tt.fields.rng,
				splitPercent:   tt.fields.splitPercent,
				splitType:      tt.fields.splitType,
			}
			if ts.rng == nil {
				ts.rng = rand.Float64
			}

			gotBucket, gotVal := ts.RunWithTraining(tt.args.targetableID, tt.fields.normalAction, tt.fields.splitAction)

			if gotVal != tt.want.value {
				t.Fatalf("RunWithTraining() value = %v, want %v", gotVal, tt.want.value)
			}

			if tt.want.split {
				if gotBucket == NoBucket {
					remapped := ts.cohortRemapper.Remap(gotBucket)
					if remapped != tt.want.bucket {
						t.Fatalf("RunWithTraining() bucket = %v (remapped %v), want %v", gotBucket, remapped, tt.want.bucket)
					}
				} else if tt.want.bucket != "" && ts.cohortRemapper.Remap(gotBucket) != tt.want.bucket {
					t.Fatalf("RunWithTraining() bucket = %v, want remapped %v", gotBucket, tt.want.bucket)
				}
			} else if gotBucket != NoBucket {
				t.Fatalf("RunWithTraining() bucket = %v, want NoBucket", gotBucket)
			}
		})
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/adsbynimbus/nimbus/internal/experimental/cohorts
// cpu: Apple M2
// BenchmarkRunWithTraining-8                       	 1813572	       658.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkRunWithTraining-8                       	 1821356	       666.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkRunWithTraining-8                       	 1778331	       657.8 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRunWithTraining(b *testing.B) {
	ts := NewTrainingSplit(
		0.5,
		SplitCohortABC,
		nil,
	)

	ts.rng = rand.Float64

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ts.RunWithTraining("bench-id", func() float64 { return 1.0 }, func() float64 { return 0.0 })
	}
}
