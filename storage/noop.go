package storage

import (
	"math"
	"github.com/prometheus/prometheus/pkg/labels"
)

type noopQuerier struct{}

func NoopQuerier() Querier {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return noopQuerier{}
}
func (noopQuerier) Select(*SelectParams, ...*labels.Matcher) (SeriesSet, Warnings, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NoopSeriesSet(), nil, nil
}
func (noopQuerier) LabelValues(name string) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
func (noopQuerier) LabelNames() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
func (noopQuerier) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}

type noopSeriesSet struct{}

func NoopSeriesSet() SeriesSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return noopSeriesSet{}
}
func (noopSeriesSet) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (noopSeriesSet) At() Series {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (noopSeriesSet) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}

type noopSeriesIterator struct{}

var NoopSeriesIt = noopSeriesIterator{}

func (noopSeriesIterator) At() (int64, float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return math.MinInt64, 0
}
func (noopSeriesIterator) Seek(t int64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (noopSeriesIterator) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (noopSeriesIterator) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
