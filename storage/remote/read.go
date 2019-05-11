package remote

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
)

var remoteReadQueries = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem, Name: "remote_read_queries", Help: "The number of in-flight remote read queries."}, []string{"client"})

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(remoteReadQueries)
}
func QueryableClient(c *Client) storage.Queryable {
	_logClusterCodePath()
	defer _logClusterCodePath()
	remoteReadQueries.WithLabelValues(c.Name())
	return storage.QueryableFunc(func(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
		return &querier{ctx: ctx, mint: mint, maxt: maxt, client: c}, nil
	})
}

type querier struct {
	ctx			context.Context
	mint, maxt	int64
	client		*Client
}

func (q *querier) Select(p *storage.SelectParams, matchers ...*labels.Matcher) (storage.SeriesSet, storage.Warnings, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	query, err := ToQuery(q.mint, q.maxt, matchers, p)
	if err != nil {
		return nil, nil, err
	}
	remoteReadGauge := remoteReadQueries.WithLabelValues(q.client.Name())
	remoteReadGauge.Inc()
	defer remoteReadGauge.Dec()
	res, err := q.client.Read(q.ctx, query)
	if err != nil {
		return nil, nil, err
	}
	return FromQueryResult(res), nil, nil
}
func (q *querier) LabelValues(name string) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
func (q *querier) LabelNames() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
func (q *querier) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func ExternalLabelsHandler(next storage.Queryable, externalLabels model.LabelSet) storage.Queryable {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return storage.QueryableFunc(func(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
		q, err := next.Querier(ctx, mint, maxt)
		if err != nil {
			return nil, err
		}
		return &externalLabelsQuerier{Querier: q, externalLabels: externalLabels}, nil
	})
}

type externalLabelsQuerier struct {
	storage.Querier
	externalLabels	model.LabelSet
}

func (q externalLabelsQuerier) Select(p *storage.SelectParams, matchers ...*labels.Matcher) (storage.SeriesSet, storage.Warnings, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m, added := q.addExternalLabels(matchers)
	s, warnings, err := q.Querier.Select(p, m...)
	if err != nil {
		return nil, warnings, err
	}
	return newSeriesSetFilter(s, added), warnings, nil
}
func PreferLocalStorageFilter(next storage.Queryable, cb startTimeCallback) storage.Queryable {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return storage.QueryableFunc(func(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
		localStartTime, err := cb()
		if err != nil {
			return nil, err
		}
		cmaxt := maxt
		if mint > localStartTime {
			return storage.NoopQuerier(), nil
		}
		if maxt > localStartTime {
			cmaxt = localStartTime
		}
		return next.Querier(ctx, mint, cmaxt)
	})
}
func RequiredMatchersFilter(next storage.Queryable, required []*labels.Matcher) storage.Queryable {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return storage.QueryableFunc(func(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
		q, err := next.Querier(ctx, mint, maxt)
		if err != nil {
			return nil, err
		}
		return &requiredMatchersQuerier{Querier: q, requiredMatchers: required}, nil
	})
}

type requiredMatchersQuerier struct {
	storage.Querier
	requiredMatchers	[]*labels.Matcher
}

func (q requiredMatchersQuerier) Select(p *storage.SelectParams, matchers ...*labels.Matcher) (storage.SeriesSet, storage.Warnings, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ms := q.requiredMatchers
	for _, m := range matchers {
		for i, r := range ms {
			if m.Type == labels.MatchEqual && m.Name == r.Name && m.Value == r.Value {
				ms = append(ms[:i], ms[i+1:]...)
				break
			}
		}
		if len(ms) == 0 {
			break
		}
	}
	if len(ms) > 0 {
		return storage.NoopSeriesSet(), nil, nil
	}
	return q.Querier.Select(p, matchers...)
}
func (q externalLabelsQuerier) addExternalLabels(ms []*labels.Matcher) ([]*labels.Matcher, model.LabelSet) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	el := make(model.LabelSet, len(q.externalLabels))
	for k, v := range q.externalLabels {
		el[k] = v
	}
	for _, m := range ms {
		if _, ok := el[model.LabelName(m.Name)]; ok {
			delete(el, model.LabelName(m.Name))
		}
	}
	for k, v := range el {
		m, err := labels.NewMatcher(labels.MatchEqual, string(k), string(v))
		if err != nil {
			panic(err)
		}
		ms = append(ms, m)
	}
	return ms, el
}
func newSeriesSetFilter(ss storage.SeriesSet, toFilter model.LabelSet) storage.SeriesSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &seriesSetFilter{SeriesSet: ss, toFilter: toFilter}
}

type seriesSetFilter struct {
	storage.SeriesSet
	toFilter	model.LabelSet
	querier		storage.Querier
}

func (ssf *seriesSetFilter) GetQuerier() storage.Querier {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ssf.querier
}
func (ssf *seriesSetFilter) SetQuerier(querier storage.Querier) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ssf.querier = querier
}
func (ssf seriesSetFilter) At() storage.Series {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return seriesFilter{Series: ssf.SeriesSet.At(), toFilter: ssf.toFilter}
}

type seriesFilter struct {
	storage.Series
	toFilter	model.LabelSet
}

func (sf seriesFilter) Labels() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	labels := sf.Series.Labels()
	for i := 0; i < len(labels); {
		if _, ok := sf.toFilter[model.LabelName(labels[i].Name)]; ok {
			labels = labels[:i+copy(labels[i:], labels[i+1:])]
			continue
		}
		i++
	}
	return labels
}
