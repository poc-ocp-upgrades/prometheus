package promql

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
)

type Value interface {
	Type() ValueType
	String() string
}

func (Matrix) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ValueTypeMatrix
}
func (Vector) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ValueTypeVector
}
func (Scalar) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ValueTypeScalar
}
func (String) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ValueTypeString
}

type ValueType string

const (
	ValueTypeNone	= "none"
	ValueTypeVector	= "vector"
	ValueTypeScalar	= "scalar"
	ValueTypeMatrix	= "matrix"
	ValueTypeString	= "string"
)

type String struct {
	T	int64
	V	string
}

func (s String) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.V
}
func (s String) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return json.Marshal([...]interface{}{float64(s.T) / 1000, s.V})
}

type Scalar struct {
	T	int64
	V	float64
}

func (s Scalar) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v := strconv.FormatFloat(s.V, 'f', -1, 64)
	return fmt.Sprintf("scalar: %v @[%v]", v, s.T)
}
func (s Scalar) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v := strconv.FormatFloat(s.V, 'f', -1, 64)
	return json.Marshal([...]interface{}{float64(s.T) / 1000, v})
}

type Series struct {
	Metric	labels.Labels	`json:"metric"`
	Points	[]Point		`json:"values"`
}

func (s Series) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	vals := make([]string, len(s.Points))
	for i, v := range s.Points {
		vals[i] = v.String()
	}
	return fmt.Sprintf("%s =>\n%s", s.Metric, strings.Join(vals, "\n"))
}

type Point struct {
	T	int64
	V	float64
}

func (p Point) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v := strconv.FormatFloat(p.V, 'f', -1, 64)
	return fmt.Sprintf("%v @[%v]", v, p.T)
}
func (p Point) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v := strconv.FormatFloat(p.V, 'f', -1, 64)
	return json.Marshal([...]interface{}{float64(p.T) / 1000, v})
}

type Sample struct {
	Point
	Metric	labels.Labels
}

func (s Sample) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s => %s", s.Metric, s.Point)
}
func (s Sample) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v := struct {
		M	labels.Labels	`json:"metric"`
		V	Point		`json:"value"`
	}{M: s.Metric, V: s.Point}
	return json.Marshal(v)
}

type Vector []Sample

func (vec Vector) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	entries := make([]string, len(vec))
	for i, s := range vec {
		entries[i] = s.String()
	}
	return strings.Join(entries, "\n")
}
func (vec Vector) ContainsSameLabelset() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := make(map[uint64]struct{}, len(vec))
	for _, s := range vec {
		hash := s.Metric.Hash()
		if _, ok := l[hash]; ok {
			return true
		}
		l[hash] = struct{}{}
	}
	return false
}

type Matrix []Series

func (m Matrix) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	strs := make([]string, len(m))
	for i, ss := range m {
		strs[i] = ss.String()
	}
	return strings.Join(strs, "\n")
}
func (m Matrix) TotalSamples() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	numSamples := 0
	for _, series := range m {
		numSamples += len(series.Points)
	}
	return numSamples
}
func (m Matrix) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(m)
}
func (m Matrix) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return labels.Compare(m[i].Metric, m[j].Metric) < 0
}
func (m Matrix) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m[i], m[j] = m[j], m[i]
}
func (m Matrix) ContainsSameLabelset() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := make(map[uint64]struct{}, len(m))
	for _, ss := range m {
		hash := ss.Metric.Hash()
		if _, ok := l[hash]; ok {
			return true
		}
		l[hash] = struct{}{}
	}
	return false
}

type Result struct {
	Err		error
	Value		Value
	Warnings	storage.Warnings
}

func (r *Result) Vector() (Vector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.Err != nil {
		return nil, r.Err
	}
	v, ok := r.Value.(Vector)
	if !ok {
		return nil, fmt.Errorf("query result is not a Vector")
	}
	return v, nil
}
func (r *Result) Matrix() (Matrix, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.Err != nil {
		return nil, r.Err
	}
	v, ok := r.Value.(Matrix)
	if !ok {
		return nil, fmt.Errorf("query result is not a range Vector")
	}
	return v, nil
}
func (r *Result) Scalar() (Scalar, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.Err != nil {
		return Scalar{}, r.Err
	}
	v, ok := r.Value.(Scalar)
	if !ok {
		return Scalar{}, fmt.Errorf("query result is not a Scalar")
	}
	return v, nil
}
func (r *Result) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.Err != nil {
		return r.Err.Error()
	}
	if r.Value == nil {
		return ""
	}
	return r.Value.String()
}

type StorageSeries struct{ series Series }

func NewStorageSeries(series Series) *StorageSeries {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &StorageSeries{series: series}
}
func (ss *StorageSeries) Labels() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ss.series.Metric
}
func (ss *StorageSeries) Iterator() storage.SeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return newStorageSeriesIterator(ss.series)
}

type storageSeriesIterator struct {
	points	[]Point
	curr	int
}

func newStorageSeriesIterator(series Series) *storageSeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &storageSeriesIterator{points: series.Points, curr: -1}
}
func (ssi *storageSeriesIterator) Seek(t int64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	i := ssi.curr
	if i < 0 {
		i = 0
	}
	for ; i < len(ssi.points); i++ {
		if ssi.points[i].T >= t {
			ssi.curr = i
			return true
		}
	}
	ssi.curr = len(ssi.points) - 1
	return false
}
func (ssi *storageSeriesIterator) At() (t int64, v float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := ssi.points[ssi.curr]
	return p.T, p.V
}
func (ssi *storageSeriesIterator) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ssi.curr++
	return ssi.curr < len(ssi.points)
}
func (ssi *storageSeriesIterator) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
