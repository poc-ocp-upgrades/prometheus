package remote

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage"
)

const decodeReadLimit = 32 * 1024 * 1024

type HTTPError struct {
	msg		string
	status	int
}

func (e HTTPError) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.msg
}
func (e HTTPError) Status() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.status
}
func DecodeReadRequest(r *http.Request) (*prompb.ReadRequest, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	compressed, err := ioutil.ReadAll(io.LimitReader(r.Body, decodeReadLimit))
	if err != nil {
		return nil, err
	}
	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		return nil, err
	}
	var req prompb.ReadRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		return nil, err
	}
	return &req, nil
}
func EncodeReadResponse(resp *prompb.ReadResponse, w http.ResponseWriter) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("Content-Encoding", "snappy")
	compressed := snappy.Encode(nil, data)
	_, err = w.Write(compressed)
	return err
}
func ToWriteRequest(samples []*model.Sample) *prompb.WriteRequest {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req := &prompb.WriteRequest{Timeseries: make([]prompb.TimeSeries, 0, len(samples))}
	for _, s := range samples {
		ts := prompb.TimeSeries{Labels: MetricToLabelProtos(s.Metric), Samples: []prompb.Sample{{Value: float64(s.Value), Timestamp: int64(s.Timestamp)}}}
		req.Timeseries = append(req.Timeseries, ts)
	}
	return req
}
func ToQuery(from, to int64, matchers []*labels.Matcher, p *storage.SelectParams) (*prompb.Query, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ms, err := toLabelMatchers(matchers)
	if err != nil {
		return nil, err
	}
	var rp *prompb.ReadHints
	if p != nil {
		rp = &prompb.ReadHints{StepMs: p.Step, Func: p.Func, StartMs: p.Start, EndMs: p.End}
	}
	return &prompb.Query{StartTimestampMs: from, EndTimestampMs: to, Matchers: ms, Hints: rp}, nil
}
func FromQuery(req *prompb.Query) (int64, int64, []*labels.Matcher, *storage.SelectParams, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	matchers, err := fromLabelMatchers(req.Matchers)
	if err != nil {
		return 0, 0, nil, nil, err
	}
	var selectParams *storage.SelectParams
	if req.Hints != nil {
		selectParams = &storage.SelectParams{Start: req.Hints.StartMs, End: req.Hints.EndMs, Step: req.Hints.StepMs, Func: req.Hints.Func}
	}
	return req.StartTimestampMs, req.EndTimestampMs, matchers, selectParams, nil
}
func ToQueryResult(ss storage.SeriesSet, sampleLimit int) (*prompb.QueryResult, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	numSamples := 0
	resp := &prompb.QueryResult{}
	for ss.Next() {
		series := ss.At()
		iter := series.Iterator()
		samples := []prompb.Sample{}
		for iter.Next() {
			numSamples++
			if sampleLimit > 0 && numSamples > sampleLimit {
				return nil, HTTPError{msg: fmt.Sprintf("exceeded sample limit (%d)", sampleLimit), status: http.StatusBadRequest}
			}
			ts, val := iter.At()
			samples = append(samples, prompb.Sample{Timestamp: ts, Value: val})
		}
		if err := iter.Err(); err != nil {
			return nil, err
		}
		resp.Timeseries = append(resp.Timeseries, &prompb.TimeSeries{Labels: labelsToLabelsProto(series.Labels()), Samples: samples})
	}
	if err := ss.Err(); err != nil {
		return nil, err
	}
	return resp, nil
}
func FromQueryResult(res *prompb.QueryResult) storage.SeriesSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	series := make([]storage.Series, 0, len(res.Timeseries))
	for _, ts := range res.Timeseries {
		labels := labelProtosToLabels(ts.Labels)
		if err := validateLabelsAndMetricName(labels); err != nil {
			return errSeriesSet{err: err}
		}
		series = append(series, &concreteSeries{labels: labels, samples: ts.Samples})
	}
	sort.Sort(byLabel(series))
	return &concreteSeriesSet{series: series}
}

type byLabel []storage.Series

func (a byLabel) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(a)
}
func (a byLabel) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	a[i], a[j] = a[j], a[i]
}
func (a byLabel) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return labels.Compare(a[i].Labels(), a[j].Labels()) < 0
}

type errSeriesSet struct{ err error }

func (errSeriesSet) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (errSeriesSet) At() storage.Series {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (e errSeriesSet) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.err
}

type concreteSeriesSet struct {
	cur		int
	series	[]storage.Series
}

func (c *concreteSeriesSet) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.cur++
	return c.cur-1 < len(c.series)
}
func (c *concreteSeriesSet) At() storage.Series {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.series[c.cur-1]
}
func (c *concreteSeriesSet) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}

type concreteSeries struct {
	labels	labels.Labels
	samples	[]prompb.Sample
}

func (c *concreteSeries) Labels() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return labels.New(c.labels...)
}
func (c *concreteSeries) Iterator() storage.SeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return newConcreteSeriersIterator(c)
}

type concreteSeriesIterator struct {
	cur		int
	series	*concreteSeries
}

func newConcreteSeriersIterator(series *concreteSeries) storage.SeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &concreteSeriesIterator{cur: -1, series: series}
}
func (c *concreteSeriesIterator) Seek(t int64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.cur = sort.Search(len(c.series.samples), func(n int) bool {
		return c.series.samples[n].Timestamp >= t
	})
	return c.cur < len(c.series.samples)
}
func (c *concreteSeriesIterator) At() (t int64, v float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := c.series.samples[c.cur]
	return s.Timestamp, s.Value
}
func (c *concreteSeriesIterator) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.cur++
	return c.cur < len(c.series.samples)
}
func (c *concreteSeriesIterator) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func validateLabelsAndMetricName(ls labels.Labels) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, l := range ls {
		if l.Name == labels.MetricName && !model.IsValidMetricName(model.LabelValue(l.Value)) {
			return fmt.Errorf("invalid metric name: %v", l.Value)
		}
		if !model.LabelName(l.Name).IsValid() {
			return fmt.Errorf("invalid label name: %v", l.Name)
		}
		if !model.LabelValue(l.Value).IsValid() {
			return fmt.Errorf("invalid label value: %v", l.Value)
		}
	}
	return nil
}
func toLabelMatchers(matchers []*labels.Matcher) ([]*prompb.LabelMatcher, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pbMatchers := make([]*prompb.LabelMatcher, 0, len(matchers))
	for _, m := range matchers {
		var mType prompb.LabelMatcher_Type
		switch m.Type {
		case labels.MatchEqual:
			mType = prompb.LabelMatcher_EQ
		case labels.MatchNotEqual:
			mType = prompb.LabelMatcher_NEQ
		case labels.MatchRegexp:
			mType = prompb.LabelMatcher_RE
		case labels.MatchNotRegexp:
			mType = prompb.LabelMatcher_NRE
		default:
			return nil, fmt.Errorf("invalid matcher type")
		}
		pbMatchers = append(pbMatchers, &prompb.LabelMatcher{Type: mType, Name: m.Name, Value: m.Value})
	}
	return pbMatchers, nil
}
func fromLabelMatchers(matchers []*prompb.LabelMatcher) ([]*labels.Matcher, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make([]*labels.Matcher, 0, len(matchers))
	for _, matcher := range matchers {
		var mtype labels.MatchType
		switch matcher.Type {
		case prompb.LabelMatcher_EQ:
			mtype = labels.MatchEqual
		case prompb.LabelMatcher_NEQ:
			mtype = labels.MatchNotEqual
		case prompb.LabelMatcher_RE:
			mtype = labels.MatchRegexp
		case prompb.LabelMatcher_NRE:
			mtype = labels.MatchNotRegexp
		default:
			return nil, fmt.Errorf("invalid matcher type")
		}
		matcher, err := labels.NewMatcher(mtype, matcher.Name, matcher.Value)
		if err != nil {
			return nil, err
		}
		result = append(result, matcher)
	}
	return result, nil
}
func MetricToLabelProtos(metric model.Metric) []prompb.Label {
	_logClusterCodePath()
	defer _logClusterCodePath()
	labels := make([]prompb.Label, 0, len(metric))
	for k, v := range metric {
		labels = append(labels, prompb.Label{Name: string(k), Value: string(v)})
	}
	sort.Slice(labels, func(i int, j int) bool {
		return labels[i].Name < labels[j].Name
	})
	return labels
}
func LabelProtosToMetric(labelPairs []*prompb.Label) model.Metric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	metric := make(model.Metric, len(labelPairs))
	for _, l := range labelPairs {
		metric[model.LabelName(l.Name)] = model.LabelValue(l.Value)
	}
	return metric
}
func labelProtosToLabels(labelPairs []prompb.Label) labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(labels.Labels, 0, len(labelPairs))
	for _, l := range labelPairs {
		result = append(result, labels.Label{Name: l.Name, Value: l.Value})
	}
	sort.Sort(result)
	return result
}
func labelsToLabelsProto(labels labels.Labels) []prompb.Label {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make([]prompb.Label, 0, len(labels))
	for _, l := range labels {
		result = append(result, prompb.Label{Name: l.Name, Value: l.Value})
	}
	return result
}
func labelsToMetric(ls labels.Labels) model.Metric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	metric := make(model.Metric, len(ls))
	for _, l := range ls {
		metric[model.LabelName(l.Name)] = model.LabelValue(l.Value)
	}
	return metric
}
