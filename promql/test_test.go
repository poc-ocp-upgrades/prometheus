package promql

import (
	"math"
	"testing"
	"time"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/util/testutil"
)

func TestLazyLoader_WithSamplesTill(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		ts		time.Time
		series		[]Series
		checkOnlyError	bool
	}
	cases := []struct {
		loadString	string
		testCases	[]testCase
	}{{loadString: `
				load 10s
					metric1 1+1x10
			`, testCases: []testCase{{ts: time.Unix(40, 0), series: []Series{{Metric: labels.FromStrings("__name__", "metric1"), Points: []Point{{0, 1}, {10000, 2}, {20000, 3}, {30000, 4}, {40000, 5}}}}}, {ts: time.Unix(10, 0), series: []Series{{Metric: labels.FromStrings("__name__", "metric1"), Points: []Point{{0, 1}, {10000, 2}, {20000, 3}, {30000, 4}, {40000, 5}}}}}, {ts: time.Unix(60, 0), series: []Series{{Metric: labels.FromStrings("__name__", "metric1"), Points: []Point{{0, 1}, {10000, 2}, {20000, 3}, {30000, 4}, {40000, 5}, {50000, 6}, {60000, 7}}}}}}}, {loadString: `
				load 10s
					metric1 1+0x5
					metric2 1+1x100
			`, testCases: []testCase{{ts: time.Unix(70, 0), series: []Series{{Metric: labels.FromStrings("__name__", "metric1"), Points: []Point{{0, 1}, {10000, 1}, {20000, 1}, {30000, 1}, {40000, 1}, {50000, 1}}}, {Metric: labels.FromStrings("__name__", "metric2"), Points: []Point{{0, 1}, {10000, 2}, {20000, 3}, {30000, 4}, {40000, 5}, {50000, 6}, {60000, 7}, {70000, 8}}}}}, {ts: time.Unix(300, 0), checkOnlyError: true}}}}
	for _, c := range cases {
		suite, err := NewLazyLoader(t, c.loadString)
		testutil.Ok(t, err)
		defer suite.Close()
		for _, tc := range c.testCases {
			suite.WithSamplesTill(tc.ts, func(err error) {
				testutil.Ok(t, err)
				if tc.checkOnlyError {
					return
				}
				queryable := suite.Queryable()
				querier, err := queryable.Querier(suite.Context(), math.MinInt64, math.MaxInt64)
				testutil.Ok(t, err)
				for _, s := range tc.series {
					var matchers []*labels.Matcher
					for _, label := range s.Metric {
						m, err := labels.NewMatcher(labels.MatchEqual, label.Name, label.Value)
						testutil.Ok(t, err)
						matchers = append(matchers, m)
					}
					ss, _, err := querier.Select(nil, matchers...)
					testutil.Ok(t, err)
					testutil.Assert(t, ss.Next(), "")
					storageSeries := ss.At()
					testutil.Assert(t, !ss.Next(), "Expecting only 1 series")
					got := Series{Metric: storageSeries.Labels()}
					it := storageSeries.Iterator()
					for it.Next() {
						t, v := it.At()
						got.Points = append(got.Points, Point{T: t, V: v})
					}
					testutil.Ok(t, it.Err())
					testutil.Equals(t, s, got)
				}
			})
		}
	}
}
