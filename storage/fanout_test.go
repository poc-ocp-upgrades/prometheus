package storage

import (
	"fmt"
	"math"
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/prometheus/prometheus/pkg/labels"
)

func TestMergeStringSlices(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, tc := range []struct {
		input		[][]string
		expected	[]string
	}{{}, {[][]string{{"foo"}}, []string{"foo"}}, {[][]string{{"foo"}, {"bar"}}, []string{"bar", "foo"}}, {[][]string{{"foo"}, {"bar"}, {"baz"}}, []string{"bar", "baz", "foo"}}} {
		require.Equal(t, tc.expected, mergeStringSlices(tc.input))
	}
}
func TestMergeTwoStringSlices(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, tc := range []struct{ a, b, expected []string }{{[]string{}, []string{}, []string{}}, {[]string{"foo"}, nil, []string{"foo"}}, {nil, []string{"bar"}, []string{"bar"}}, {[]string{"foo"}, []string{"bar"}, []string{"bar", "foo"}}, {[]string{"foo"}, []string{"bar", "baz"}, []string{"bar", "baz", "foo"}}, {[]string{"foo"}, []string{"foo"}, []string{"foo"}}} {
		require.Equal(t, tc.expected, mergeTwoStringSlices(tc.a, tc.b))
	}
}
func TestMergeSeriesSet(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, tc := range []struct {
		input		[]SeriesSet
		expected	SeriesSet
	}{{input: []SeriesSet{newMockSeriesSet()}, expected: newMockSeriesSet()}, {input: []SeriesSet{newMockSeriesSet(newMockSeries(labels.FromStrings("bar", "baz"), []sample{{1, 1}, {2, 2}}), newMockSeries(labels.FromStrings("foo", "bar"), []sample{{0, 0}, {1, 1}}))}, expected: newMockSeriesSet(newMockSeries(labels.FromStrings("bar", "baz"), []sample{{1, 1}, {2, 2}}), newMockSeries(labels.FromStrings("foo", "bar"), []sample{{0, 0}, {1, 1}}))}, {input: []SeriesSet{newMockSeriesSet(newMockSeries(labels.FromStrings("foo", "bar"), []sample{{0, 0}, {1, 1}})), newMockSeriesSet(newMockSeries(labels.FromStrings("bar", "baz"), []sample{{1, 1}, {2, 2}}))}, expected: newMockSeriesSet(newMockSeries(labels.FromStrings("bar", "baz"), []sample{{1, 1}, {2, 2}}), newMockSeries(labels.FromStrings("foo", "bar"), []sample{{0, 0}, {1, 1}}))}, {input: []SeriesSet{newMockSeriesSet(newMockSeries(labels.FromStrings("bar", "baz"), []sample{{1, 1}, {2, 2}}), newMockSeries(labels.FromStrings("foo", "bar"), []sample{{0, 0}, {1, 1}})), newMockSeriesSet(newMockSeries(labels.FromStrings("bar", "baz"), []sample{{3, 3}, {4, 4}}), newMockSeries(labels.FromStrings("foo", "bar"), []sample{{2, 2}, {3, 3}}))}, expected: newMockSeriesSet(newMockSeries(labels.FromStrings("bar", "baz"), []sample{{1, 1}, {2, 2}, {3, 3}, {4, 4}}), newMockSeries(labels.FromStrings("foo", "bar"), []sample{{0, 0}, {1, 1}, {2, 2}, {3, 3}}))}, {input: []SeriesSet{newMockSeriesSet(newMockSeries(labels.FromStrings("foo", "bar"), []sample{{0, math.NaN()}})), newMockSeriesSet(newMockSeries(labels.FromStrings("foo", "bar"), []sample{{0, math.NaN()}}))}, expected: newMockSeriesSet(newMockSeries(labels.FromStrings("foo", "bar"), []sample{{0, math.NaN()}}))}} {
		merged := NewMergeSeriesSet(tc.input, nil)
		for merged.Next() {
			require.True(t, tc.expected.Next())
			actualSeries := merged.At()
			expectedSeries := tc.expected.At()
			require.Equal(t, expectedSeries.Labels(), actualSeries.Labels())
			require.Equal(t, drainSamples(expectedSeries.Iterator()), drainSamples(actualSeries.Iterator()))
		}
		require.False(t, tc.expected.Next())
	}
}
func TestMergeIterator(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, tc := range []struct {
		input		[]SeriesIterator
		expected	[]sample
	}{{input: []SeriesIterator{newListSeriesIterator([]sample{{0, 0}, {1, 1}})}, expected: []sample{{0, 0}, {1, 1}}}, {input: []SeriesIterator{newListSeriesIterator([]sample{{0, 0}, {1, 1}}), newListSeriesIterator([]sample{{2, 2}, {3, 3}})}, expected: []sample{{0, 0}, {1, 1}, {2, 2}, {3, 3}}}, {input: []SeriesIterator{newListSeriesIterator([]sample{{0, 0}, {3, 3}}), newListSeriesIterator([]sample{{1, 1}, {4, 4}}), newListSeriesIterator([]sample{{2, 2}, {5, 5}})}, expected: []sample{{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}}}, {input: []SeriesIterator{newListSeriesIterator([]sample{{0, 0}, {1, 1}}), newListSeriesIterator([]sample{{0, 0}, {2, 2}}), newListSeriesIterator([]sample{{2, 2}, {3, 3}})}, expected: []sample{{0, 0}, {1, 1}, {2, 2}, {3, 3}}}} {
		merged := newMergeIterator(tc.input)
		actual := drainSamples(merged)
		require.Equal(t, tc.expected, actual)
	}
}
func TestMergeIteratorSeek(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, tc := range []struct {
		input		[]SeriesIterator
		seek		int64
		expected	[]sample
	}{{input: []SeriesIterator{newListSeriesIterator([]sample{{0, 0}, {1, 1}, {2, 2}})}, seek: 1, expected: []sample{{1, 1}, {2, 2}}}, {input: []SeriesIterator{newListSeriesIterator([]sample{{0, 0}, {1, 1}}), newListSeriesIterator([]sample{{2, 2}, {3, 3}})}, seek: 2, expected: []sample{{2, 2}, {3, 3}}}, {input: []SeriesIterator{newListSeriesIterator([]sample{{0, 0}, {3, 3}}), newListSeriesIterator([]sample{{1, 1}, {4, 4}}), newListSeriesIterator([]sample{{2, 2}, {5, 5}})}, seek: 2, expected: []sample{{2, 2}, {3, 3}, {4, 4}, {5, 5}}}} {
		merged := newMergeIterator(tc.input)
		actual := []sample{}
		if merged.Seek(tc.seek) {
			t, v := merged.At()
			actual = append(actual, sample{t, v})
		}
		actual = append(actual, drainSamples(merged)...)
		require.Equal(t, tc.expected, actual)
	}
}
func drainSamples(iter SeriesIterator) []sample {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := []sample{}
	for iter.Next() {
		t, v := iter.At()
		if math.IsNaN(v) {
			v = -42
		}
		result = append(result, sample{t, v})
	}
	return result
}

type mockSeriesSet struct {
	idx	int
	series	[]Series
}

func newMockSeriesSet(series ...Series) SeriesSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &mockSeriesSet{idx: -1, series: series}
}
func (m *mockSeriesSet) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.idx++
	return m.idx < len(m.series)
}
func (m *mockSeriesSet) At() Series {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.series[m.idx]
}
func (m *mockSeriesSet) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}

var result []sample

func makeSeriesSet(numSeries, numSamples int) SeriesSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	series := []Series{}
	for j := 0; j < numSeries; j++ {
		labels := labels.Labels{{Name: "foo", Value: fmt.Sprintf("bar%d", j)}}
		samples := []sample{}
		for k := 0; k < numSamples; k++ {
			samples = append(samples, sample{t: int64(k), v: float64(k)})
		}
		series = append(series, newMockSeries(labels, samples))
	}
	return newMockSeriesSet(series...)
}
func makeMergeSeriesSet(numSeriesSets, numSeries, numSamples int) SeriesSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	seriesSets := []SeriesSet{}
	for i := 0; i < numSeriesSets; i++ {
		seriesSets = append(seriesSets, makeSeriesSet(numSeries, numSamples))
	}
	return NewMergeSeriesSet(seriesSets, nil)
}
func benchmarkDrain(seriesSet SeriesSet, b *testing.B) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for n := 0; n < b.N; n++ {
		for seriesSet.Next() {
			result = drainSamples(seriesSet.At().Iterator())
		}
	}
}
func BenchmarkNoMergeSeriesSet_100_100(b *testing.B) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	seriesSet := makeSeriesSet(100, 100)
	benchmarkDrain(seriesSet, b)
}
func BenchmarkMergeSeriesSet(b *testing.B) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, bm := range []struct{ numSeriesSets, numSeries, numSamples int }{{1, 100, 100}, {10, 100, 100}, {100, 100, 100}} {
		seriesSet := makeMergeSeriesSet(bm.numSeriesSets, bm.numSeries, bm.numSamples)
		b.Run(fmt.Sprintf("%d_%d_%d", bm.numSeriesSets, bm.numSeries, bm.numSamples), func(b *testing.B) {
			benchmarkDrain(seriesSet, b)
		})
	}
}
