

# cohorts ![Go](https://github.com/marcsantiago/go-cohort/workflows/Go/badge.svg?branch=master)
`import "github.com/marcsantiago/go-cohort"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
A simple data-science package, which allows deterministic hashes to segment data by cohort ![Go](<a href="https://github.com/marcsantiago/go-cohort/workflows/Go/badge.svg?branch=master">https://github.com/marcsantiago/go-cohort/workflows/Go/badge.svg?branch=master</a>)




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [type Bucket](#Bucket)
  * [func AssignCohort(identifier string, splitType SplitType) Bucket](#AssignCohort)
  * [func AssignCohortAB(identifier string) Bucket](#AssignCohortAB)
  * [func AssignCohortABC(identifier string) Bucket](#AssignCohortABC)
  * [func AssignMultipleCohorts(identifier string, splitBy []SplitType) Bucket](#AssignMultipleCohorts)
* [type SplitType](#SplitType)
* [type SplitTypes](#SplitTypes)
  * [func (s SplitTypes) Len() int](#SplitTypes.Len)
  * [func (s SplitTypes) Less(i, j int) bool](#SplitTypes.Less)
  * [func (s SplitTypes) Swap(i, j int)](#SplitTypes.Swap)


#### <a name="pkg-files">Package files</a>
[bucket.go](/src/github.com/marcsantiago/go-cohort/bucket.go) [cohorts.go](/src/github.com/marcsantiago/go-cohort/cohorts.go) [doc.go](/src/github.com/marcsantiago/go-cohort/doc.go) [split_type.go](/src/github.com/marcsantiago/go-cohort/split_type.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    BucketA = "A"
    BucketB = "B"
    BucketC = "C"
)
```
``` go
const (
    SplitCohortAB = iota + 2
    SplitCohortABC
)
```
``` go
const SplitTypeUnknown = SplitType(-1)
```
returns -1 because when using the value in mod math the split will always be 0 and therefore fall into group A without a B group
this in turn means the split is not running, but also nothing is panicking or breaking.





## <a name="Bucket">type</a> [Bucket](/src/target/bucket.go?s=95:113#L4)
``` go
type Bucket string
```
Bucket provides a string representation of the Bucket the cohort landed in







### <a name="AssignCohort">func</a> [AssignCohort](/src/target/cohorts.go?s=909:973#L40)
``` go
func AssignCohort(identifier string, splitType SplitType) Bucket
```
AssignCohorts returns a Bucket which is a string representation of A,B,or C contingent on the split type


### <a name="AssignCohortAB">func</a> [AssignCohortAB](/src/target/cohorts.go?s=1224:1269#L52)
``` go
func AssignCohortAB(identifier string) Bucket
```
AssignCohortAB calls AssignCohort and fixes the split type to A/B


### <a name="AssignCohortABC">func</a> [AssignCohortABC](/src/target/cohorts.go?s=1394:1440#L57)
``` go
func AssignCohortABC(identifier string) Bucket
```
AssignCohortAB calls AssignCohort and fixes the split type to A/B/C


### <a name="AssignMultipleCohorts">func</a> [AssignMultipleCohorts](/src/target/cohorts.go?s=2119:2192#L67)
``` go
func AssignMultipleCohorts(identifier string, splitBy []SplitType) Bucket
```
AssignMultipleCohorts generates a bucket that merges the cohort on multiple split types
e.g two different tests running on the same user where each test has a different split type assigned
users see that a blue banner running an A/B test and users that see cats, dogs, or clowns as an A/B/C where
the users are the same and the tests are running at the same time but the spit is different, test 1 the user is assigned
bucket A and in test 2 the user is a assigned bucket C, so we return AC as the bucket type. Buckets are always sorted so
A will represent that status on A/B testing and C on A/B/C testing





## <a name="SplitType">type</a> [SplitType](/src/target/split_type.go?s=32:51#L5)
``` go
type SplitType int8
```









## <a name="SplitTypes">type</a> [SplitTypes](/src/target/split_type.go?s=441:468#L17)
``` go
type SplitTypes []SplitType
```
Buckets is an alias for []Bucket with the sort interface implemented










### <a name="SplitTypes.Len">func</a> (SplitTypes) [Len](/src/target/split_type.go?s=470:499#L19)
``` go
func (s SplitTypes) Len() int
```



### <a name="SplitTypes.Less">func</a> (SplitTypes) [Less](/src/target/split_type.go?s=528:567#L20)
``` go
func (s SplitTypes) Less(i, j int) bool
```



### <a name="SplitTypes.Swap">func</a> (SplitTypes) [Swap](/src/target/split_type.go?s=591:625#L21)
``` go
func (s SplitTypes) Swap(i, j int)
```







- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
