package storage

import (
	"math/rand"
	"sort"
	"testing"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/stretchr/testify/require"
)

func TestSampleRing(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cases := []struct {
		input	[]int64
		delta	int64
		size	int
	}{{input: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, delta: 2, size: 1}, {input: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, delta: 2, size: 2}, {input: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, delta: 7, size: 3}, {input: []int64{1, 2, 3, 4, 5, 16, 17, 18, 19, 20}, delta: 7, size: 1}, {input: []int64{1, 2, 3, 4, 6}, delta: 4, size: 4}}
	for _, c := range cases {
		r := newSampleRing(c.delta, c.size)
		input := []sample{}
		for _, t := range c.input {
			input = append(input, sample{t: t, v: float64(rand.Intn(100))})
		}
		for i, s := range input {
			r.add(s.t, s.v)
			buffered := r.samples()
			for _, sold := range input[:i] {
				found := false
				for _, bs := range buffered {
					if bs.t == sold.t && bs.v == sold.v {
						found = true
						break
					}
				}
				if sold.t >= s.t-c.delta && !found {
					t.Fatalf("%d: expected sample %d to be in buffer but was not; buffer %v", i, sold.t, buffered)
				}
				if sold.t < s.t-c.delta && found {
					t.Fatalf("%d: unexpected sample %d in buffer; buffer %v", i, sold.t, buffered)
				}
			}
		}
	}
}
func TestBufferedSeriesIterator(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var it *BufferedSeriesIterator
	bufferEq := func(exp []sample) {
		var b []sample
		bit := it.Buffer()
		for bit.Next() {
			t, v := bit.At()
			b = append(b, sample{t: t, v: v})
		}
		require.Equal(t, exp, b, "buffer mismatch")
	}
	sampleEq := func(ets int64, ev float64) {
		ts, v := it.Values()
		require.Equal(t, ets, ts, "timestamp mismatch")
		require.Equal(t, ev, v, "value mismatch")
	}
	it = NewBufferIterator(newListSeriesIterator([]sample{{t: 1, v: 2}, {t: 2, v: 3}, {t: 3, v: 4}, {t: 4, v: 5}, {t: 5, v: 6}, {t: 99, v: 8}, {t: 100, v: 9}, {t: 101, v: 10}}), 2)
	require.True(t, it.Seek(-123), "seek failed")
	sampleEq(1, 2)
	bufferEq(nil)
	require.True(t, it.Next(), "next failed")
	sampleEq(2, 3)
	bufferEq([]sample{{t: 1, v: 2}})
	require.True(t, it.Next(), "next failed")
	require.True(t, it.Next(), "next failed")
	require.True(t, it.Next(), "next failed")
	sampleEq(5, 6)
	bufferEq([]sample{{t: 2, v: 3}, {t: 3, v: 4}, {t: 4, v: 5}})
	require.True(t, it.Seek(5), "seek failed")
	sampleEq(5, 6)
	bufferEq([]sample{{t: 2, v: 3}, {t: 3, v: 4}, {t: 4, v: 5}})
	require.True(t, it.Seek(101), "seek failed")
	sampleEq(101, 10)
	bufferEq([]sample{{t: 99, v: 8}, {t: 100, v: 9}})
	require.False(t, it.Next(), "next succeeded unexpectedly")
}
func TestBufferedSeriesIteratorNoBadAt(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	done := false
	m := &mockSeriesIterator{seek: func(int64) bool {
		return false
	}, at: func() (int64, float64) {
		require.False(t, done)
		done = true
		return 0, 0
	}, next: func() bool {
		return !done
	}, err: func() error {
		return nil
	}}
	it := NewBufferIterator(m, 60)
	it.Next()
	it.Next()
}
func BenchmarkBufferedSeriesIterator(b *testing.B) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	it := NewBufferIterator(newFakeSeriesIterator(int64(b.N), 30), 5*60)
	b.SetBytes(int64(b.N * 16))
	b.ReportAllocs()
	b.ResetTimer()
	for it.Next() {
	}
	require.NoError(b, it.Err())
}

type mockSeriesIterator struct {
	seek	func(int64) bool
	at	func() (int64, float64)
	next	func() bool
	err	func() error
}

func (m *mockSeriesIterator) Seek(t int64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.seek(t)
}
func (m *mockSeriesIterator) At() (int64, float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.at()
}
func (m *mockSeriesIterator) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.next()
}
func (m *mockSeriesIterator) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.err()
}

type mockSeries struct {
	labels		func() labels.Labels
	iterator	func() SeriesIterator
}

func newMockSeries(lset labels.Labels, samples []sample) Series {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &mockSeries{labels: func() labels.Labels {
		return lset
	}, iterator: func() SeriesIterator {
		return newListSeriesIterator(samples)
	}}
}
func (m *mockSeries) Labels() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.labels()
}
func (m *mockSeries) Iterator() SeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.iterator()
}

type listSeriesIterator struct {
	list	[]sample
	idx	int
}

func newListSeriesIterator(list []sample) *listSeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &listSeriesIterator{list: list, idx: -1}
}
func (it *listSeriesIterator) At() (int64, float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := it.list[it.idx]
	return s.t, s.v
}
func (it *listSeriesIterator) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	it.idx++
	return it.idx < len(it.list)
}
func (it *listSeriesIterator) Seek(t int64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if it.idx == -1 {
		it.idx = 0
	}
	it.idx = sort.Search(len(it.list)-it.idx, func(i int) bool {
		s := it.list[i+it.idx]
		return s.t >= t
	})
	return it.idx < len(it.list)
}
func (it *listSeriesIterator) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}

type fakeSeriesIterator struct {
	nsamples	int64
	step		int64
	idx		int64
}

func newFakeSeriesIterator(nsamples, step int64) *fakeSeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &fakeSeriesIterator{nsamples: nsamples, step: step, idx: -1}
}
func (it *fakeSeriesIterator) At() (int64, float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return it.idx * it.step, 123
}
func (it *fakeSeriesIterator) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	it.idx++
	return it.idx < it.nsamples
}
func (it *fakeSeriesIterator) Seek(t int64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	it.idx = t / it.step
	return it.idx < it.nsamples
}
func (it *fakeSeriesIterator) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
