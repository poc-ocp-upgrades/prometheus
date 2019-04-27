package storage

import (
	"container/heap"
	"context"
	"sort"
	"strings"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
)

type fanout struct {
	logger		log.Logger
	primary		Storage
	secondaries	[]Storage
}

func NewFanout(logger log.Logger, primary Storage, secondaries ...Storage) Storage {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &fanout{logger: logger, primary: primary, secondaries: secondaries}
}
func (f *fanout) StartTime() (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	firstTime, err := f.primary.StartTime()
	if err != nil {
		return int64(model.Latest), err
	}
	for _, storage := range f.secondaries {
		t, err := storage.StartTime()
		if err != nil {
			return int64(model.Latest), err
		}
		if t < firstTime {
			firstTime = t
		}
	}
	return firstTime, nil
}
func (f *fanout) Querier(ctx context.Context, mint, maxt int64) (Querier, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	queriers := make([]Querier, 0, 1+len(f.secondaries))
	primaryQuerier, err := f.primary.Querier(ctx, mint, maxt)
	if err != nil {
		return nil, err
	}
	queriers = append(queriers, primaryQuerier)
	for _, storage := range f.secondaries {
		querier, err := storage.Querier(ctx, mint, maxt)
		if err != nil {
			NewMergeQuerier(primaryQuerier, queriers).Close()
			return nil, err
		}
		queriers = append(queriers, querier)
	}
	return NewMergeQuerier(primaryQuerier, queriers), nil
}
func (f *fanout) Appender() (Appender, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	primary, err := f.primary.Appender()
	if err != nil {
		return nil, err
	}
	secondaries := make([]Appender, 0, len(f.secondaries))
	for _, storage := range f.secondaries {
		appender, err := storage.Appender()
		if err != nil {
			return nil, err
		}
		secondaries = append(secondaries, appender)
	}
	return &fanoutAppender{logger: f.logger, primary: primary, secondaries: secondaries}, nil
}
func (f *fanout) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := f.primary.Close(); err != nil {
		return err
	}
	var lastErr error
	for _, storage := range f.secondaries {
		if err := storage.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

type fanoutAppender struct {
	logger		log.Logger
	primary		Appender
	secondaries	[]Appender
}

func (f *fanoutAppender) Add(l labels.Labels, t int64, v float64) (uint64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ref, err := f.primary.Add(l, t, v)
	if err != nil {
		return ref, err
	}
	for _, appender := range f.secondaries {
		if _, err := appender.Add(l, t, v); err != nil {
			return 0, err
		}
	}
	return ref, nil
}
func (f *fanoutAppender) AddFast(l labels.Labels, ref uint64, t int64, v float64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := f.primary.AddFast(l, ref, t, v); err != nil {
		return err
	}
	for _, appender := range f.secondaries {
		if _, err := appender.Add(l, t, v); err != nil {
			return err
		}
	}
	return nil
}
func (f *fanoutAppender) Commit() (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = f.primary.Commit()
	for _, appender := range f.secondaries {
		if err == nil {
			err = appender.Commit()
		} else {
			if rollbackErr := appender.Rollback(); rollbackErr != nil {
				level.Error(f.logger).Log("msg", "Squashed rollback error on commit", "err", rollbackErr)
			}
		}
	}
	return
}
func (f *fanoutAppender) Rollback() (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = f.primary.Rollback()
	for _, appender := range f.secondaries {
		rollbackErr := appender.Rollback()
		if err == nil {
			err = rollbackErr
		} else if rollbackErr != nil {
			level.Error(f.logger).Log("msg", "Squashed rollback error on rollback", "err", rollbackErr)
		}
	}
	return nil
}

type mergeQuerier struct {
	primaryQuerier	Querier
	queriers	[]Querier
	failedQueriers	map[Querier]struct{}
	setQuerierMap	map[SeriesSet]Querier
}

func NewMergeQuerier(primaryQuerier Querier, queriers []Querier) Querier {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	filtered := make([]Querier, 0, len(queriers))
	for _, querier := range queriers {
		if querier != NoopQuerier() {
			filtered = append(filtered, querier)
		}
	}
	setQuerierMap := make(map[SeriesSet]Querier)
	failedQueriers := make(map[Querier]struct{})
	switch len(filtered) {
	case 0:
		return NoopQuerier()
	case 1:
		return filtered[0]
	default:
		return &mergeQuerier{primaryQuerier: primaryQuerier, queriers: filtered, failedQueriers: failedQueriers, setQuerierMap: setQuerierMap}
	}
}
func (q *mergeQuerier) Select(params *SelectParams, matchers ...*labels.Matcher) (SeriesSet, Warnings, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	seriesSets := make([]SeriesSet, 0, len(q.queriers))
	var warnings Warnings
	for _, querier := range q.queriers {
		set, wrn, err := querier.Select(params, matchers...)
		q.setQuerierMap[set] = querier
		if wrn != nil {
			warnings = append(warnings, wrn...)
		}
		if err != nil {
			q.failedQueriers[querier] = struct{}{}
			if querier != q.primaryQuerier {
				warnings = append(warnings, err)
				continue
			} else {
				return nil, nil, err
			}
		}
		seriesSets = append(seriesSets, set)
	}
	return NewMergeSeriesSet(seriesSets, q), warnings, nil
}
func (q *mergeQuerier) LabelValues(name string) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var results [][]string
	for _, querier := range q.queriers {
		values, err := querier.LabelValues(name)
		if err != nil {
			return nil, err
		}
		results = append(results, values)
	}
	return mergeStringSlices(results), nil
}
func (q *mergeQuerier) IsFailedSet(set SeriesSet) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, isFailedQuerier := q.failedQueriers[q.setQuerierMap[set]]
	return isFailedQuerier
}
func mergeStringSlices(ss [][]string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch len(ss) {
	case 0:
		return nil
	case 1:
		return ss[0]
	case 2:
		return mergeTwoStringSlices(ss[0], ss[1])
	default:
		halfway := len(ss) / 2
		return mergeTwoStringSlices(mergeStringSlices(ss[:halfway]), mergeStringSlices(ss[halfway:]))
	}
}
func mergeTwoStringSlices(a, b []string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	i, j := 0, 0
	result := make([]string, 0, len(a)+len(b))
	for i < len(a) && j < len(b) {
		switch strings.Compare(a[i], b[j]) {
		case 0:
			result = append(result, a[i])
			i++
			j++
		case -1:
			result = append(result, a[i])
			i++
		case 1:
			result = append(result, b[j])
			j++
		}
	}
	result = append(result, a[i:]...)
	result = append(result, b[j:]...)
	return result
}
func (q *mergeQuerier) LabelNames() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	labelNamesMap := make(map[string]struct{})
	for _, b := range q.queriers {
		names, err := b.LabelNames()
		if err != nil {
			return nil, errors.Wrap(err, "LabelNames() from Querier")
		}
		for _, name := range names {
			labelNamesMap[name] = struct{}{}
		}
	}
	labelNames := make([]string, 0, len(labelNamesMap))
	for name := range labelNamesMap {
		labelNames = append(labelNames, name)
	}
	sort.Strings(labelNames)
	return labelNames, nil
}
func (q *mergeQuerier) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var lastErr error
	for _, querier := range q.queriers {
		if err := querier.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

type mergeSeriesSet struct {
	currentLabels	labels.Labels
	currentSets	[]SeriesSet
	heap		seriesSetHeap
	sets		[]SeriesSet
	querier		*mergeQuerier
}

func NewMergeSeriesSet(sets []SeriesSet, querier *mergeQuerier) SeriesSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(sets) == 1 {
		return sets[0]
	}
	var h seriesSetHeap
	for _, set := range sets {
		if set == nil {
			continue
		}
		if set.Next() {
			heap.Push(&h, set)
		}
	}
	return &mergeSeriesSet{heap: h, sets: sets, querier: querier}
}
func (c *mergeSeriesSet) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for {
		for _, set := range c.currentSets {
			if set.Next() {
				heap.Push(&c.heap, set)
			}
		}
		if len(c.heap) == 0 {
			return false
		}
		c.currentSets = nil
		c.currentLabels = c.heap[0].At().Labels()
		for len(c.heap) > 0 && labels.Equal(c.currentLabels, c.heap[0].At().Labels()) {
			set := heap.Pop(&c.heap).(SeriesSet)
			if c.querier != nil && c.querier.IsFailedSet(set) {
				continue
			}
			c.currentSets = append(c.currentSets, set)
		}
		if len(c.currentSets) != 0 {
			break
		}
	}
	return true
}
func (c *mergeSeriesSet) At() Series {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(c.currentSets) == 1 {
		return c.currentSets[0].At()
	}
	series := []Series{}
	for _, seriesSet := range c.currentSets {
		series = append(series, seriesSet.At())
	}
	return &mergeSeries{labels: c.currentLabels, series: series}
}
func (c *mergeSeriesSet) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, set := range c.sets {
		if err := set.Err(); err != nil {
			return err
		}
	}
	return nil
}

type seriesSetHeap []SeriesSet

func (h seriesSetHeap) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(h)
}
func (h seriesSetHeap) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h[i], h[j] = h[j], h[i]
}
func (h seriesSetHeap) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	a, b := h[i].At().Labels(), h[j].At().Labels()
	return labels.Compare(a, b) < 0
}
func (h *seriesSetHeap) Push(x interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*h = append(*h, x.(SeriesSet))
}
func (h *seriesSetHeap) Pop() interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type mergeSeries struct {
	labels	labels.Labels
	series	[]Series
}

func (m *mergeSeries) Labels() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.labels
}
func (m *mergeSeries) Iterator() SeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	iterators := make([]SeriesIterator, 0, len(m.series))
	for _, s := range m.series {
		iterators = append(iterators, s.Iterator())
	}
	return newMergeIterator(iterators)
}

type mergeIterator struct {
	iterators	[]SeriesIterator
	h		seriesIteratorHeap
}

func newMergeIterator(iterators []SeriesIterator) SeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &mergeIterator{iterators: iterators, h: nil}
}
func (c *mergeIterator) Seek(t int64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.h = seriesIteratorHeap{}
	for _, iter := range c.iterators {
		if iter.Seek(t) {
			heap.Push(&c.h, iter)
		}
	}
	return len(c.h) > 0
}
func (c *mergeIterator) At() (t int64, v float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(c.h) == 0 {
		panic("mergeIterator.At() called after .Next() returned false.")
	}
	return c.h[0].At()
}
func (c *mergeIterator) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.h == nil {
		for _, iter := range c.iterators {
			if iter.Next() {
				heap.Push(&c.h, iter)
			}
		}
		return len(c.h) > 0
	}
	if len(c.h) == 0 {
		return false
	}
	currt, _ := c.At()
	for len(c.h) > 0 {
		nextt, _ := c.h[0].At()
		if nextt != currt {
			break
		}
		iter := heap.Pop(&c.h).(SeriesIterator)
		if iter.Next() {
			heap.Push(&c.h, iter)
		}
	}
	return len(c.h) > 0
}
func (c *mergeIterator) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, iter := range c.iterators {
		if err := iter.Err(); err != nil {
			return err
		}
	}
	return nil
}

type seriesIteratorHeap []SeriesIterator

func (h seriesIteratorHeap) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(h)
}
func (h seriesIteratorHeap) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h[i], h[j] = h[j], h[i]
}
func (h seriesIteratorHeap) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	at, _ := h[i].At()
	bt, _ := h[j].At()
	return at < bt
}
func (h *seriesIteratorHeap) Push(x interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*h = append(*h, x.(SeriesIterator))
}
func (h *seriesIteratorHeap) Pop() interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
