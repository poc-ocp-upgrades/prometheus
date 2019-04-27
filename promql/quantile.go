package promql

import (
	"math"
	"sort"
	"github.com/prometheus/prometheus/pkg/labels"
)

var excludedLabels = []string{labels.MetricName, labels.BucketLabel}

type bucket struct {
	upperBound	float64
	count		float64
}
type buckets []bucket

func (b buckets) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(b)
}
func (b buckets) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	b[i], b[j] = b[j], b[i]
}
func (b buckets) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return b[i].upperBound < b[j].upperBound
}

type metricWithBuckets struct {
	metric	labels.Labels
	buckets	buckets
}

func bucketQuantile(q float64, buckets buckets) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if q < 0 {
		return math.Inf(-1)
	}
	if q > 1 {
		return math.Inf(+1)
	}
	if len(buckets) < 2 {
		return math.NaN()
	}
	sort.Sort(buckets)
	if !math.IsInf(buckets[len(buckets)-1].upperBound, +1) {
		return math.NaN()
	}
	ensureMonotonic(buckets)
	rank := q * buckets[len(buckets)-1].count
	b := sort.Search(len(buckets)-1, func(i int) bool {
		return buckets[i].count >= rank
	})
	if b == len(buckets)-1 {
		return buckets[len(buckets)-2].upperBound
	}
	if b == 0 && buckets[0].upperBound <= 0 {
		return buckets[0].upperBound
	}
	var (
		bucketStart	float64
		bucketEnd	= buckets[b].upperBound
		count		= buckets[b].count
	)
	if b > 0 {
		bucketStart = buckets[b-1].upperBound
		count -= buckets[b-1].count
		rank -= buckets[b-1].count
	}
	return bucketStart + (bucketEnd-bucketStart)*(rank/count)
}
func ensureMonotonic(buckets buckets) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	max := buckets[0].count
	for i := range buckets[1:] {
		switch {
		case buckets[i].count > max:
			max = buckets[i].count
		case buckets[i].count < max:
			buckets[i].count = max
		}
	}
}
func quantile(q float64, values vectorByValueHeap) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(values) == 0 {
		return math.NaN()
	}
	if q < 0 {
		return math.Inf(-1)
	}
	if q > 1 {
		return math.Inf(+1)
	}
	sort.Sort(values)
	n := float64(len(values))
	rank := q * (n - 1)
	lowerIndex := math.Max(0, math.Floor(rank))
	upperIndex := math.Min(n-1, lowerIndex+1)
	weight := rank - math.Floor(rank)
	return values[int(lowerIndex)].V*(1-weight) + values[int(upperIndex)].V*weight
}
