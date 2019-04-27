package promql

import (
	"container/heap"
	"context"
	"fmt"
	"math"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/gate"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/timestamp"
	"github.com/prometheus/prometheus/pkg/value"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/stats"
)

const (
	namespace	= "prometheus"
	subsystem	= "engine"
	queryTag	= "query"
	env		= "query execution"
	maxInt64	= 9223372036854774784
	minInt64	= -9223372036854775808
)

var (
	LookbackDelta			= 5 * time.Minute
	DefaultEvaluationInterval	int64
)

func SetDefaultEvaluationInterval(ev time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	atomic.StoreInt64(&DefaultEvaluationInterval, durationToInt64Millis(ev))
}
func GetDefaultEvaluationInterval() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return atomic.LoadInt64(&DefaultEvaluationInterval)
}

type engineMetrics struct {
	currentQueries		prometheus.Gauge
	maxConcurrentQueries	prometheus.Gauge
	queryQueueTime		prometheus.Summary
	queryPrepareTime	prometheus.Summary
	queryInnerEval		prometheus.Summary
	queryResultSort		prometheus.Summary
}

func convertibleToInt64(v float64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return v <= maxInt64 && v >= minInt64
}

type (
	ErrQueryTimeout		string
	ErrQueryCanceled	string
	ErrTooManySamples	string
	ErrStorage		struct{ Err error }
)

func (e ErrQueryTimeout) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("query timed out in %s", string(e))
}
func (e ErrQueryCanceled) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("query was canceled in %s", string(e))
}
func (e ErrTooManySamples) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("query processing would load too many samples into memory in %s", string(e))
}
func (e ErrStorage) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.Err.Error()
}

type Query interface {
	Exec(ctx context.Context) *Result
	Close()
	Statement() Statement
	Stats() *stats.QueryTimers
	Cancel()
}
type query struct {
	queryable	storage.Queryable
	q		string
	stmt		Statement
	stats		*stats.QueryTimers
	matrix		Matrix
	cancel		func()
	ng		*Engine
}

func (q *query) Statement() Statement {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return q.stmt
}
func (q *query) Stats() *stats.QueryTimers {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return q.stats
}
func (q *query) Cancel() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if q.cancel != nil {
		q.cancel()
	}
}
func (q *query) Close() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, s := range q.matrix {
		putPointSlice(s.Points)
	}
}
func (q *query) Exec(ctx context.Context) *Result {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetTag(queryTag, q.stmt.String())
	}
	res, warnings, err := q.ng.exec(ctx, q)
	return &Result{Err: err, Value: res, Warnings: warnings}
}
func contextDone(ctx context.Context, env string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	select {
	case <-ctx.Done():
		return contextErr(ctx.Err(), env)
	default:
		return nil
	}
}
func contextErr(err error, env string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch err {
	case context.Canceled:
		return ErrQueryCanceled(env)
	case context.DeadlineExceeded:
		return ErrQueryTimeout(env)
	default:
		return err
	}
}

type EngineOpts struct {
	Logger		log.Logger
	Reg		prometheus.Registerer
	MaxConcurrent	int
	MaxSamples	int
	Timeout		time.Duration
}
type Engine struct {
	logger			log.Logger
	metrics			*engineMetrics
	timeout			time.Duration
	gate			*gate.Gate
	maxSamplesPerQuery	int
}

func NewEngine(opts EngineOpts) *Engine {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if opts.Logger == nil {
		opts.Logger = log.NewNopLogger()
	}
	metrics := &engineMetrics{currentQueries: prometheus.NewGauge(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem, Name: "queries", Help: "The current number of queries being executed or waiting."}), maxConcurrentQueries: prometheus.NewGauge(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem, Name: "queries_concurrent_max", Help: "The max number of concurrent queries."}), queryQueueTime: prometheus.NewSummary(prometheus.SummaryOpts{Namespace: namespace, Subsystem: subsystem, Name: "query_duration_seconds", Help: "Query timings", ConstLabels: prometheus.Labels{"slice": "queue_time"}}), queryPrepareTime: prometheus.NewSummary(prometheus.SummaryOpts{Namespace: namespace, Subsystem: subsystem, Name: "query_duration_seconds", Help: "Query timings", ConstLabels: prometheus.Labels{"slice": "prepare_time"}}), queryInnerEval: prometheus.NewSummary(prometheus.SummaryOpts{Namespace: namespace, Subsystem: subsystem, Name: "query_duration_seconds", Help: "Query timings", ConstLabels: prometheus.Labels{"slice": "inner_eval"}}), queryResultSort: prometheus.NewSummary(prometheus.SummaryOpts{Namespace: namespace, Subsystem: subsystem, Name: "query_duration_seconds", Help: "Query timings", ConstLabels: prometheus.Labels{"slice": "result_sort"}})}
	metrics.maxConcurrentQueries.Set(float64(opts.MaxConcurrent))
	if opts.Reg != nil {
		opts.Reg.MustRegister(metrics.currentQueries, metrics.maxConcurrentQueries, metrics.queryQueueTime, metrics.queryPrepareTime, metrics.queryInnerEval, metrics.queryResultSort)
	}
	return &Engine{gate: gate.New(opts.MaxConcurrent), timeout: opts.Timeout, logger: opts.Logger, metrics: metrics, maxSamplesPerQuery: opts.MaxSamples}
}
func (ng *Engine) NewInstantQuery(q storage.Queryable, qs string, ts time.Time) (Query, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	expr, err := ParseExpr(qs)
	if err != nil {
		return nil, err
	}
	qry := ng.newQuery(q, expr, ts, ts, 0)
	qry.q = qs
	return qry, nil
}
func (ng *Engine) NewRangeQuery(q storage.Queryable, qs string, start, end time.Time, interval time.Duration) (Query, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	expr, err := ParseExpr(qs)
	if err != nil {
		return nil, err
	}
	if expr.Type() != ValueTypeVector && expr.Type() != ValueTypeScalar {
		return nil, fmt.Errorf("invalid expression type %q for range query, must be Scalar or instant Vector", documentedType(expr.Type()))
	}
	qry := ng.newQuery(q, expr, start, end, interval)
	qry.q = qs
	return qry, nil
}
func (ng *Engine) newQuery(q storage.Queryable, expr Expr, start, end time.Time, interval time.Duration) *query {
	_logClusterCodePath()
	defer _logClusterCodePath()
	es := &EvalStmt{Expr: expr, Start: start, End: end, Interval: interval}
	qry := &query{stmt: es, ng: ng, stats: stats.NewQueryTimers(), queryable: q}
	return qry
}

type testStmt func(context.Context) error

func (testStmt) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "test statement"
}
func (testStmt) stmt() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (ng *Engine) newTestQuery(f func(context.Context) error) Query {
	_logClusterCodePath()
	defer _logClusterCodePath()
	qry := &query{q: "test statement", stmt: testStmt(f), ng: ng, stats: stats.NewQueryTimers()}
	return qry
}
func (ng *Engine) exec(ctx context.Context, q *query) (Value, storage.Warnings, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ng.metrics.currentQueries.Inc()
	defer ng.metrics.currentQueries.Dec()
	ctx, cancel := context.WithTimeout(ctx, ng.timeout)
	q.cancel = cancel
	execSpanTimer, ctx := q.stats.GetSpanTimer(ctx, stats.ExecTotalTime)
	defer execSpanTimer.Finish()
	queueSpanTimer, _ := q.stats.GetSpanTimer(ctx, stats.ExecQueueTime, ng.metrics.queryQueueTime)
	if err := ng.gate.Start(ctx); err != nil {
		return nil, nil, contextErr(err, "query queue")
	}
	defer ng.gate.Done()
	queueSpanTimer.Finish()
	defer q.cancel()
	const env = "query execution"
	evalSpanTimer, ctx := q.stats.GetSpanTimer(ctx, stats.EvalTotalTime)
	defer evalSpanTimer.Finish()
	if err := contextDone(ctx, env); err != nil {
		return nil, nil, err
	}
	switch s := q.Statement().(type) {
	case *EvalStmt:
		return ng.execEvalStmt(ctx, q, s)
	case testStmt:
		return nil, nil, s(ctx)
	}
	panic(fmt.Errorf("promql.Engine.exec: unhandled statement of type %T", q.Statement()))
}
func timeMilliseconds(t time.Time) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.UnixNano() / int64(time.Millisecond/time.Nanosecond)
}
func durationMilliseconds(d time.Duration) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return int64(d / (time.Millisecond / time.Nanosecond))
}
func (ng *Engine) execEvalStmt(ctx context.Context, query *query, s *EvalStmt) (Value, storage.Warnings, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	prepareSpanTimer, ctxPrepare := query.stats.GetSpanTimer(ctx, stats.QueryPreparationTime, ng.metrics.queryPrepareTime)
	querier, warnings, err := ng.populateSeries(ctxPrepare, query.queryable, s)
	prepareSpanTimer.Finish()
	if querier != nil {
		defer querier.Close()
	}
	if err != nil {
		return nil, warnings, err
	}
	evalSpanTimer, _ := query.stats.GetSpanTimer(ctx, stats.InnerEvalTime, ng.metrics.queryInnerEval)
	if s.Start == s.End && s.Interval == 0 {
		start := timeMilliseconds(s.Start)
		evaluator := &evaluator{startTimestamp: start, endTimestamp: start, interval: 1, ctx: ctx, maxSamples: ng.maxSamplesPerQuery, defaultEvalInterval: GetDefaultEvaluationInterval(), logger: ng.logger}
		val, err := evaluator.Eval(s.Expr)
		if err != nil {
			return nil, warnings, err
		}
		evalSpanTimer.Finish()
		mat, ok := val.(Matrix)
		if !ok {
			panic(fmt.Errorf("promql.Engine.exec: invalid expression type %q", val.Type()))
		}
		query.matrix = mat
		switch s.Expr.Type() {
		case ValueTypeVector:
			vector := make(Vector, len(mat))
			for i, s := range mat {
				vector[i] = Sample{Metric: s.Metric, Point: Point{V: s.Points[0].V, T: start}}
			}
			return vector, warnings, nil
		case ValueTypeScalar:
			return Scalar{V: mat[0].Points[0].V, T: start}, warnings, nil
		case ValueTypeMatrix:
			return mat, warnings, nil
		default:
			panic(fmt.Errorf("promql.Engine.exec: unexpected expression type %q", s.Expr.Type()))
		}
	}
	evaluator := &evaluator{startTimestamp: timeMilliseconds(s.Start), endTimestamp: timeMilliseconds(s.End), interval: durationMilliseconds(s.Interval), ctx: ctx, maxSamples: ng.maxSamplesPerQuery, defaultEvalInterval: GetDefaultEvaluationInterval(), logger: ng.logger}
	val, err := evaluator.Eval(s.Expr)
	if err != nil {
		return nil, warnings, err
	}
	evalSpanTimer.Finish()
	mat, ok := val.(Matrix)
	if !ok {
		panic(fmt.Errorf("promql.Engine.exec: invalid expression type %q", val.Type()))
	}
	query.matrix = mat
	if err := contextDone(ctx, "expression evaluation"); err != nil {
		return nil, warnings, err
	}
	sortSpanTimer, _ := query.stats.GetSpanTimer(ctx, stats.ResultSortTime, ng.metrics.queryResultSort)
	sort.Sort(mat)
	sortSpanTimer.Finish()
	return mat, warnings, nil
}
func (ng *Engine) cumulativeSubqueryOffset(path []Node) time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var subqOffset time.Duration
	for _, node := range path {
		switch n := node.(type) {
		case *SubqueryExpr:
			subqOffset += n.Range + n.Offset
		}
	}
	return subqOffset
}
func (ng *Engine) populateSeries(ctx context.Context, q storage.Queryable, s *EvalStmt) (storage.Querier, storage.Warnings, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var maxOffset time.Duration
	Inspect(s.Expr, func(node Node, path []Node) error {
		subqOffset := ng.cumulativeSubqueryOffset(path)
		switch n := node.(type) {
		case *VectorSelector:
			if maxOffset < LookbackDelta+subqOffset {
				maxOffset = LookbackDelta + subqOffset
			}
			if n.Offset+LookbackDelta+subqOffset > maxOffset {
				maxOffset = n.Offset + LookbackDelta + subqOffset
			}
		case *MatrixSelector:
			if maxOffset < n.Range+subqOffset {
				maxOffset = n.Range + subqOffset
			}
			if n.Offset+n.Range+subqOffset > maxOffset {
				maxOffset = n.Offset + n.Range + subqOffset
			}
		}
		return nil
	})
	mint := s.Start.Add(-maxOffset)
	querier, err := q.Querier(ctx, timestamp.FromTime(mint), timestamp.FromTime(s.End))
	if err != nil {
		return nil, nil, err
	}
	var warnings storage.Warnings
	Inspect(s.Expr, func(node Node, path []Node) error {
		var set storage.SeriesSet
		var wrn storage.Warnings
		params := &storage.SelectParams{Start: timestamp.FromTime(s.Start), End: timestamp.FromTime(s.End), Step: durationToInt64Millis(s.Interval)}
		switch n := node.(type) {
		case *VectorSelector:
			params.Start = params.Start - durationMilliseconds(LookbackDelta)
			params.Func = extractFuncFromPath(path)
			if n.Offset > 0 {
				offsetMilliseconds := durationMilliseconds(n.Offset)
				params.Start = params.Start - offsetMilliseconds
				params.End = params.End - offsetMilliseconds
			}
			set, wrn, err = querier.Select(params, n.LabelMatchers...)
			warnings = append(warnings, wrn...)
			if err != nil {
				level.Error(ng.logger).Log("msg", "error selecting series set", "err", err)
				return err
			}
			n.unexpandedSeriesSet = set
		case *MatrixSelector:
			params.Func = extractFuncFromPath(path)
			params.Start = params.Start - durationMilliseconds(n.Range)
			if n.Offset > 0 {
				offsetMilliseconds := durationMilliseconds(n.Offset)
				params.Start = params.Start - offsetMilliseconds
				params.End = params.End - offsetMilliseconds
			}
			set, wrn, err = querier.Select(params, n.LabelMatchers...)
			warnings = append(warnings, wrn...)
			if err != nil {
				level.Error(ng.logger).Log("msg", "error selecting series set", "err", err)
				return err
			}
			n.unexpandedSeriesSet = set
		}
		return nil
	})
	return querier, warnings, err
}
func extractFuncFromPath(p []Node) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(p) == 0 {
		return ""
	}
	switch n := p[len(p)-1].(type) {
	case *AggregateExpr:
		return n.Op.String()
	case *Call:
		return n.Func.Name
	case *BinaryExpr:
		return ""
	}
	return extractFuncFromPath(p[:len(p)-1])
}
func checkForSeriesSetExpansion(ctx context.Context, expr Expr) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch e := expr.(type) {
	case *MatrixSelector:
		if e.series == nil {
			series, err := expandSeriesSet(ctx, e.unexpandedSeriesSet)
			if err != nil {
				panic(err)
			} else {
				e.series = series
			}
		}
	case *VectorSelector:
		if e.series == nil {
			series, err := expandSeriesSet(ctx, e.unexpandedSeriesSet)
			if err != nil {
				panic(err)
			} else {
				e.series = series
			}
		}
	}
	return nil
}
func expandSeriesSet(ctx context.Context, it storage.SeriesSet) (res []storage.Series, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for it.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		res = append(res, it.At())
	}
	return res, it.Err()
}

type evaluator struct {
	ctx			context.Context
	startTimestamp		int64
	endTimestamp		int64
	interval		int64
	maxSamples		int
	currentSamples		int
	defaultEvalInterval	int64
	logger			log.Logger
}

func (ev *evaluator) errorf(format string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ev.error(fmt.Errorf(format, args...))
}
func (ev *evaluator) error(err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	panic(err)
}
func (ev *evaluator) recover(errp *error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	e := recover()
	if e == nil {
		return
	}
	if err, ok := e.(runtime.Error); ok {
		buf := make([]byte, 64<<10)
		buf = buf[:runtime.Stack(buf, false)]
		level.Error(ev.logger).Log("msg", "runtime panic in parser", "err", e, "stacktrace", string(buf))
		*errp = fmt.Errorf("unexpected error: %s", err)
	} else {
		*errp = e.(error)
	}
}
func (ev *evaluator) Eval(expr Expr) (v Value, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer ev.recover(&err)
	return ev.eval(expr), nil
}

type EvalNodeHelper struct {
	ts				int64
	out				Vector
	dmn				map[uint64]labels.Labels
	sigf				map[uint64]uint64
	signatureToMetricWithBuckets	map[uint64]*metricWithBuckets
	regex				*regexp.Regexp
	rightSigs			map[uint64]Sample
	matchedSigs			map[uint64]map[uint64]struct{}
	resultMetric			map[uint64]labels.Labels
}

func (enh *EvalNodeHelper) dropMetricName(l labels.Labels) labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if enh.dmn == nil {
		enh.dmn = make(map[uint64]labels.Labels, len(enh.out))
	}
	h := l.Hash()
	ret, ok := enh.dmn[h]
	if ok {
		return ret
	}
	ret = dropMetricName(l)
	enh.dmn[h] = ret
	return ret
}
func (enh *EvalNodeHelper) signatureFunc(on bool, names ...string) func(labels.Labels) uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if enh.sigf == nil {
		enh.sigf = make(map[uint64]uint64, len(enh.out))
	}
	f := signatureFunc(on, names...)
	return func(l labels.Labels) uint64 {
		h := l.Hash()
		ret, ok := enh.sigf[h]
		if ok {
			return ret
		}
		ret = f(l)
		enh.sigf[h] = ret
		return ret
	}
}
func (ev *evaluator) rangeEval(f func([]Value, *EvalNodeHelper) Vector, exprs ...Expr) Matrix {
	_logClusterCodePath()
	defer _logClusterCodePath()
	numSteps := int((ev.endTimestamp-ev.startTimestamp)/ev.interval) + 1
	matrixes := make([]Matrix, len(exprs))
	origMatrixes := make([]Matrix, len(exprs))
	originalNumSamples := ev.currentSamples
	for i, e := range exprs {
		if e != nil && e.Type() != ValueTypeString {
			matrixes[i] = ev.eval(e).(Matrix)
			origMatrixes[i] = make(Matrix, len(matrixes[i]))
			copy(origMatrixes[i], matrixes[i])
		}
	}
	vectors := make([]Vector, len(exprs))
	args := make([]Value, len(exprs))
	biggestLen := 1
	for i := range exprs {
		vectors[i] = make(Vector, 0, len(matrixes[i]))
		if len(matrixes[i]) > biggestLen {
			biggestLen = len(matrixes[i])
		}
	}
	enh := &EvalNodeHelper{out: make(Vector, 0, biggestLen)}
	seriess := make(map[uint64]Series, biggestLen)
	tempNumSamples := ev.currentSamples
	for ts := ev.startTimestamp; ts <= ev.endTimestamp; ts += ev.interval {
		ev.currentSamples = tempNumSamples
		for i := range exprs {
			vectors[i] = vectors[i][:0]
			for si, series := range matrixes[i] {
				for _, point := range series.Points {
					if point.T == ts {
						if ev.currentSamples < ev.maxSamples {
							vectors[i] = append(vectors[i], Sample{Metric: series.Metric, Point: point})
							matrixes[i][si].Points = series.Points[1:]
							ev.currentSamples++
						} else {
							ev.error(ErrTooManySamples(env))
						}
					}
					break
				}
			}
			args[i] = vectors[i]
		}
		enh.ts = ts
		result := f(args, enh)
		if result.ContainsSameLabelset() {
			ev.errorf("vector cannot contain metrics with the same labelset")
		}
		enh.out = result[:0]
		ev.currentSamples += len(result)
		tempNumSamples += len(result)
		if ev.currentSamples > ev.maxSamples {
			ev.error(ErrTooManySamples(env))
		}
		if ev.endTimestamp == ev.startTimestamp {
			mat := make(Matrix, len(result))
			for i, s := range result {
				s.Point.T = ts
				mat[i] = Series{Metric: s.Metric, Points: []Point{s.Point}}
			}
			ev.currentSamples = originalNumSamples + mat.TotalSamples()
			return mat
		}
		for _, sample := range result {
			h := sample.Metric.Hash()
			ss, ok := seriess[h]
			if !ok {
				ss = Series{Metric: sample.Metric, Points: getPointSlice(numSteps)}
			}
			sample.Point.T = ts
			ss.Points = append(ss.Points, sample.Point)
			seriess[h] = ss
		}
	}
	for _, m := range origMatrixes {
		for _, s := range m {
			putPointSlice(s.Points)
		}
	}
	mat := make(Matrix, 0, len(seriess))
	for _, ss := range seriess {
		mat = append(mat, ss)
	}
	ev.currentSamples = originalNumSamples + mat.TotalSamples()
	return mat
}
func (ev *evaluator) evalSubquery(subq *SubqueryExpr) *MatrixSelector {
	_logClusterCodePath()
	defer _logClusterCodePath()
	val := ev.eval(subq).(Matrix)
	ms := &MatrixSelector{Range: subq.Range, Offset: subq.Offset, series: make([]storage.Series, 0, len(val))}
	for _, s := range val {
		ms.series = append(ms.series, NewStorageSeries(s))
	}
	return ms
}
func (ev *evaluator) eval(expr Expr) Value {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := contextDone(ev.ctx, "expression evaluation"); err != nil {
		ev.error(err)
	}
	numSteps := int((ev.endTimestamp-ev.startTimestamp)/ev.interval) + 1
	switch e := expr.(type) {
	case *AggregateExpr:
		if s, ok := e.Param.(*StringLiteral); ok {
			return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
				return ev.aggregation(e.Op, e.Grouping, e.Without, s.Val, v[0].(Vector), enh)
			}, e.Expr)
		}
		return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
			var param float64
			if e.Param != nil {
				param = v[0].(Vector)[0].V
			}
			return ev.aggregation(e.Op, e.Grouping, e.Without, param, v[1].(Vector), enh)
		}, e.Param, e.Expr)
	case *Call:
		if e.Func.Name == "timestamp" {
			vs, ok := e.Args[0].(*VectorSelector)
			if ok {
				return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
					return e.Func.Call([]Value{ev.vectorSelector(vs, enh.ts)}, e.Args, enh)
				})
			}
		}
		var matrixArgIndex int
		var matrixArg bool
		for i, a := range e.Args {
			if _, ok := a.(*MatrixSelector); ok {
				matrixArgIndex = i
				matrixArg = true
				break
			}
			if subq, ok := a.(*SubqueryExpr); ok {
				matrixArgIndex = i
				matrixArg = true
				e.Args[i] = ev.evalSubquery(subq)
				break
			}
		}
		if !matrixArg {
			return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
				return e.Func.Call(v, e.Args, enh)
			}, e.Args...)
		}
		inArgs := make([]Value, len(e.Args))
		otherArgs := make([]Matrix, len(e.Args))
		otherInArgs := make([]Vector, len(e.Args))
		for i, e := range e.Args {
			if i != matrixArgIndex {
				otherArgs[i] = ev.eval(e).(Matrix)
				otherInArgs[i] = Vector{Sample{}}
				inArgs[i] = otherInArgs[i]
			}
		}
		sel := e.Args[matrixArgIndex].(*MatrixSelector)
		if err := checkForSeriesSetExpansion(ev.ctx, sel); err != nil {
			ev.error(err)
		}
		mat := make(Matrix, 0, len(sel.series))
		offset := durationMilliseconds(sel.Offset)
		selRange := durationMilliseconds(sel.Range)
		stepRange := selRange
		if stepRange > ev.interval {
			stepRange = ev.interval
		}
		points := getPointSlice(16)
		inMatrix := make(Matrix, 1)
		inArgs[matrixArgIndex] = inMatrix
		enh := &EvalNodeHelper{out: make(Vector, 0, 1)}
		it := storage.NewBuffer(selRange)
		for i, s := range sel.series {
			points = points[:0]
			it.Reset(s.Iterator())
			ss := Series{Metric: dropMetricName(sel.series[i].Labels()), Points: getPointSlice(numSteps)}
			inMatrix[0].Metric = sel.series[i].Labels()
			for ts, step := ev.startTimestamp, -1; ts <= ev.endTimestamp; ts += ev.interval {
				step++
				for j := range e.Args {
					if j != matrixArgIndex {
						otherInArgs[j][0].V = otherArgs[j][0].Points[step].V
					}
				}
				maxt := ts - offset
				mint := maxt - selRange
				points = ev.matrixIterSlice(it, mint, maxt, points)
				if len(points) == 0 {
					continue
				}
				inMatrix[0].Points = points
				enh.ts = ts
				outVec := e.Func.Call(inArgs, e.Args, enh)
				enh.out = outVec[:0]
				if len(outVec) > 0 {
					ss.Points = append(ss.Points, Point{V: outVec[0].Point.V, T: ts})
				}
				it.ReduceDelta(stepRange)
			}
			if len(ss.Points) > 0 {
				if ev.currentSamples < ev.maxSamples {
					mat = append(mat, ss)
					ev.currentSamples += len(ss.Points)
				} else {
					ev.error(ErrTooManySamples(env))
				}
			}
		}
		if mat.ContainsSameLabelset() {
			ev.errorf("vector cannot contain metrics with the same labelset")
		}
		putPointSlice(points)
		return mat
	case *ParenExpr:
		return ev.eval(e.Expr)
	case *UnaryExpr:
		mat := ev.eval(e.Expr).(Matrix)
		if e.Op == itemSUB {
			for i := range mat {
				mat[i].Metric = dropMetricName(mat[i].Metric)
				for j := range mat[i].Points {
					mat[i].Points[j].V = -mat[i].Points[j].V
				}
			}
			if mat.ContainsSameLabelset() {
				ev.errorf("vector cannot contain metrics with the same labelset")
			}
		}
		return mat
	case *BinaryExpr:
		switch lt, rt := e.LHS.Type(), e.RHS.Type(); {
		case lt == ValueTypeScalar && rt == ValueTypeScalar:
			return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
				val := scalarBinop(e.Op, v[0].(Vector)[0].Point.V, v[1].(Vector)[0].Point.V)
				return append(enh.out, Sample{Point: Point{V: val}})
			}, e.LHS, e.RHS)
		case lt == ValueTypeVector && rt == ValueTypeVector:
			switch e.Op {
			case itemLAND:
				return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
					return ev.VectorAnd(v[0].(Vector), v[1].(Vector), e.VectorMatching, enh)
				}, e.LHS, e.RHS)
			case itemLOR:
				return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
					return ev.VectorOr(v[0].(Vector), v[1].(Vector), e.VectorMatching, enh)
				}, e.LHS, e.RHS)
			case itemLUnless:
				return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
					return ev.VectorUnless(v[0].(Vector), v[1].(Vector), e.VectorMatching, enh)
				}, e.LHS, e.RHS)
			default:
				return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
					return ev.VectorBinop(e.Op, v[0].(Vector), v[1].(Vector), e.VectorMatching, e.ReturnBool, enh)
				}, e.LHS, e.RHS)
			}
		case lt == ValueTypeVector && rt == ValueTypeScalar:
			return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
				return ev.VectorscalarBinop(e.Op, v[0].(Vector), Scalar{V: v[1].(Vector)[0].Point.V}, false, e.ReturnBool, enh)
			}, e.LHS, e.RHS)
		case lt == ValueTypeScalar && rt == ValueTypeVector:
			return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
				return ev.VectorscalarBinop(e.Op, v[1].(Vector), Scalar{V: v[0].(Vector)[0].Point.V}, true, e.ReturnBool, enh)
			}, e.LHS, e.RHS)
		}
	case *NumberLiteral:
		return ev.rangeEval(func(v []Value, enh *EvalNodeHelper) Vector {
			return append(enh.out, Sample{Point: Point{V: e.Val}})
		})
	case *VectorSelector:
		if err := checkForSeriesSetExpansion(ev.ctx, e); err != nil {
			ev.error(err)
		}
		mat := make(Matrix, 0, len(e.series))
		it := storage.NewBuffer(durationMilliseconds(LookbackDelta))
		for i, s := range e.series {
			it.Reset(s.Iterator())
			ss := Series{Metric: e.series[i].Labels(), Points: getPointSlice(numSteps)}
			for ts := ev.startTimestamp; ts <= ev.endTimestamp; ts += ev.interval {
				_, v, ok := ev.vectorSelectorSingle(it, e, ts)
				if ok {
					if ev.currentSamples < ev.maxSamples {
						ss.Points = append(ss.Points, Point{V: v, T: ts})
						ev.currentSamples++
					} else {
						ev.error(ErrTooManySamples(env))
					}
				}
			}
			if len(ss.Points) > 0 {
				mat = append(mat, ss)
			}
		}
		return mat
	case *MatrixSelector:
		if ev.startTimestamp != ev.endTimestamp {
			panic(fmt.Errorf("cannot do range evaluation of matrix selector"))
		}
		return ev.matrixSelector(e)
	case *SubqueryExpr:
		offsetMillis := durationToInt64Millis(e.Offset)
		rangeMillis := durationToInt64Millis(e.Range)
		newEv := &evaluator{endTimestamp: ev.endTimestamp - offsetMillis, interval: ev.defaultEvalInterval, ctx: ev.ctx, currentSamples: ev.currentSamples, maxSamples: ev.maxSamples, defaultEvalInterval: ev.defaultEvalInterval, logger: ev.logger}
		if e.Step != 0 {
			newEv.interval = durationToInt64Millis(e.Step)
		}
		newEv.startTimestamp = newEv.interval * ((ev.startTimestamp - offsetMillis - rangeMillis) / newEv.interval)
		if newEv.startTimestamp < (ev.startTimestamp - offsetMillis - rangeMillis) {
			newEv.startTimestamp += newEv.interval
		}
		res := newEv.eval(e.Expr)
		ev.currentSamples = newEv.currentSamples
		return res
	}
	panic(fmt.Errorf("unhandled expression of type: %T", expr))
}
func durationToInt64Millis(d time.Duration) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return int64(d / time.Millisecond)
}
func (ev *evaluator) vectorSelector(node *VectorSelector, ts int64) Vector {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := checkForSeriesSetExpansion(ev.ctx, node); err != nil {
		ev.error(err)
	}
	var (
		vec = make(Vector, 0, len(node.series))
	)
	it := storage.NewBuffer(durationMilliseconds(LookbackDelta))
	for i, s := range node.series {
		it.Reset(s.Iterator())
		t, v, ok := ev.vectorSelectorSingle(it, node, ts)
		if ok {
			vec = append(vec, Sample{Metric: node.series[i].Labels(), Point: Point{V: v, T: t}})
			ev.currentSamples++
		}
		if ev.currentSamples >= ev.maxSamples {
			ev.error(ErrTooManySamples(env))
		}
	}
	return vec
}
func (ev *evaluator) vectorSelectorSingle(it *storage.BufferedSeriesIterator, node *VectorSelector, ts int64) (int64, float64, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	refTime := ts - durationMilliseconds(node.Offset)
	var t int64
	var v float64
	ok := it.Seek(refTime)
	if !ok {
		if it.Err() != nil {
			ev.error(it.Err())
		}
	}
	if ok {
		t, v = it.Values()
	}
	if !ok || t > refTime {
		t, v, ok = it.PeekBack(1)
		if !ok || t < refTime-durationMilliseconds(LookbackDelta) {
			return 0, 0, false
		}
	}
	if value.IsStaleNaN(v) {
		return 0, 0, false
	}
	return t, v, true
}

var pointPool = sync.Pool{}

func getPointSlice(sz int) []Point {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := pointPool.Get()
	if p != nil {
		return p.([]Point)
	}
	return make([]Point, 0, sz)
}
func putPointSlice(p []Point) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pointPool.Put(p[:0])
}
func (ev *evaluator) matrixSelector(node *MatrixSelector) Matrix {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := checkForSeriesSetExpansion(ev.ctx, node); err != nil {
		ev.error(err)
	}
	var (
		offset	= durationMilliseconds(node.Offset)
		maxt	= ev.startTimestamp - offset
		mint	= maxt - durationMilliseconds(node.Range)
		matrix	= make(Matrix, 0, len(node.series))
	)
	it := storage.NewBuffer(durationMilliseconds(node.Range))
	for i, s := range node.series {
		if err := contextDone(ev.ctx, "expression evaluation"); err != nil {
			ev.error(err)
		}
		it.Reset(s.Iterator())
		ss := Series{Metric: node.series[i].Labels()}
		ss.Points = ev.matrixIterSlice(it, mint, maxt, getPointSlice(16))
		if len(ss.Points) > 0 {
			matrix = append(matrix, ss)
		} else {
			putPointSlice(ss.Points)
		}
	}
	return matrix
}
func (ev *evaluator) matrixIterSlice(it *storage.BufferedSeriesIterator, mint, maxt int64, out []Point) []Point {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(out) > 0 && out[len(out)-1].T >= mint {
		var drop int
		for drop = 0; out[drop].T < mint; drop++ {
		}
		copy(out, out[drop:])
		out = out[:len(out)-drop]
		mint = out[len(out)-1].T + 1
	} else {
		out = out[:0]
	}
	ok := it.Seek(maxt)
	if !ok {
		if it.Err() != nil {
			ev.error(it.Err())
		}
	}
	buf := it.Buffer()
	for buf.Next() {
		t, v := buf.At()
		if value.IsStaleNaN(v) {
			continue
		}
		if t >= mint {
			if ev.currentSamples >= ev.maxSamples {
				ev.error(ErrTooManySamples(env))
			}
			out = append(out, Point{T: t, V: v})
			ev.currentSamples++
		}
	}
	if ok {
		t, v := it.Values()
		if t == maxt && !value.IsStaleNaN(v) {
			if ev.currentSamples >= ev.maxSamples {
				ev.error(ErrTooManySamples(env))
			}
			out = append(out, Point{T: t, V: v})
			ev.currentSamples++
		}
	}
	return out
}
func (ev *evaluator) VectorAnd(lhs, rhs Vector, matching *VectorMatching, enh *EvalNodeHelper) Vector {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if matching.Card != CardManyToMany {
		panic("set operations must only use many-to-many matching")
	}
	sigf := enh.signatureFunc(matching.On, matching.MatchingLabels...)
	rightSigs := map[uint64]struct{}{}
	for _, rs := range rhs {
		rightSigs[sigf(rs.Metric)] = struct{}{}
	}
	for _, ls := range lhs {
		if _, ok := rightSigs[sigf(ls.Metric)]; ok {
			enh.out = append(enh.out, ls)
		}
	}
	return enh.out
}
func (ev *evaluator) VectorOr(lhs, rhs Vector, matching *VectorMatching, enh *EvalNodeHelper) Vector {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if matching.Card != CardManyToMany {
		panic("set operations must only use many-to-many matching")
	}
	sigf := enh.signatureFunc(matching.On, matching.MatchingLabels...)
	leftSigs := map[uint64]struct{}{}
	for _, ls := range lhs {
		leftSigs[sigf(ls.Metric)] = struct{}{}
		enh.out = append(enh.out, ls)
	}
	for _, rs := range rhs {
		if _, ok := leftSigs[sigf(rs.Metric)]; !ok {
			enh.out = append(enh.out, rs)
		}
	}
	return enh.out
}
func (ev *evaluator) VectorUnless(lhs, rhs Vector, matching *VectorMatching, enh *EvalNodeHelper) Vector {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if matching.Card != CardManyToMany {
		panic("set operations must only use many-to-many matching")
	}
	sigf := enh.signatureFunc(matching.On, matching.MatchingLabels...)
	rightSigs := map[uint64]struct{}{}
	for _, rs := range rhs {
		rightSigs[sigf(rs.Metric)] = struct{}{}
	}
	for _, ls := range lhs {
		if _, ok := rightSigs[sigf(ls.Metric)]; !ok {
			enh.out = append(enh.out, ls)
		}
	}
	return enh.out
}
func (ev *evaluator) VectorBinop(op ItemType, lhs, rhs Vector, matching *VectorMatching, returnBool bool, enh *EvalNodeHelper) Vector {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if matching.Card == CardManyToMany {
		panic("many-to-many only allowed for set operators")
	}
	sigf := enh.signatureFunc(matching.On, matching.MatchingLabels...)
	if matching.Card == CardOneToMany {
		lhs, rhs = rhs, lhs
	}
	if enh.rightSigs == nil {
		enh.rightSigs = make(map[uint64]Sample, len(enh.out))
	} else {
		for k := range enh.rightSigs {
			delete(enh.rightSigs, k)
		}
	}
	rightSigs := enh.rightSigs
	for _, rs := range rhs {
		sig := sigf(rs.Metric)
		if _, found := rightSigs[sig]; found {
			ev.errorf("many-to-many matching not allowed: matching labels must be unique on one side")
		}
		rightSigs[sig] = rs
	}
	if enh.matchedSigs == nil {
		enh.matchedSigs = make(map[uint64]map[uint64]struct{}, len(rightSigs))
	} else {
		for k := range enh.matchedSigs {
			delete(enh.matchedSigs, k)
		}
	}
	matchedSigs := enh.matchedSigs
	for _, ls := range lhs {
		sig := sigf(ls.Metric)
		rs, found := rightSigs[sig]
		if !found {
			continue
		}
		vl, vr := ls.V, rs.V
		if matching.Card == CardOneToMany {
			vl, vr = vr, vl
		}
		value, keep := vectorElemBinop(op, vl, vr)
		if returnBool {
			if keep {
				value = 1.0
			} else {
				value = 0.0
			}
		} else if !keep {
			continue
		}
		metric := resultMetric(ls.Metric, rs.Metric, op, matching, enh)
		insertedSigs, exists := matchedSigs[sig]
		if matching.Card == CardOneToOne {
			if exists {
				ev.errorf("multiple matches for labels: many-to-one matching must be explicit (group_left/group_right)")
			}
			matchedSigs[sig] = nil
		} else {
			insertSig := metric.Hash()
			if !exists {
				insertedSigs = map[uint64]struct{}{}
				matchedSigs[sig] = insertedSigs
			} else if _, duplicate := insertedSigs[insertSig]; duplicate {
				ev.errorf("multiple matches for labels: grouping labels must ensure unique matches")
			}
			insertedSigs[insertSig] = struct{}{}
		}
		enh.out = append(enh.out, Sample{Metric: metric, Point: Point{V: value}})
	}
	return enh.out
}
func signatureFunc(on bool, names ...string) func(labels.Labels) uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if on {
		return func(lset labels.Labels) uint64 {
			return lset.HashForLabels(names...)
		}
	}
	return func(lset labels.Labels) uint64 {
		return lset.HashWithoutLabels(names...)
	}
}
func resultMetric(lhs, rhs labels.Labels, op ItemType, matching *VectorMatching, enh *EvalNodeHelper) labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if enh.resultMetric == nil {
		enh.resultMetric = make(map[uint64]labels.Labels, len(enh.out))
	}
	lh := lhs.Hash()
	h := (lh ^ rhs.Hash()) + lh
	if ret, ok := enh.resultMetric[h]; ok {
		return ret
	}
	lb := labels.NewBuilder(lhs)
	if shouldDropMetricName(op) {
		lb.Del(labels.MetricName)
	}
	if matching.Card == CardOneToOne {
		if matching.On {
		Outer:
			for _, l := range lhs {
				for _, n := range matching.MatchingLabels {
					if l.Name == n {
						continue Outer
					}
				}
				lb.Del(l.Name)
			}
		} else {
			lb.Del(matching.MatchingLabels...)
		}
	}
	for _, ln := range matching.Include {
		if v := rhs.Get(ln); v != "" {
			lb.Set(ln, v)
		} else {
			lb.Del(ln)
		}
	}
	ret := lb.Labels()
	enh.resultMetric[h] = ret
	return ret
}
func (ev *evaluator) VectorscalarBinop(op ItemType, lhs Vector, rhs Scalar, swap, returnBool bool, enh *EvalNodeHelper) Vector {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, lhsSample := range lhs {
		lv, rv := lhsSample.V, rhs.V
		if swap {
			lv, rv = rv, lv
		}
		value, keep := vectorElemBinop(op, lv, rv)
		if returnBool {
			if keep {
				value = 1.0
			} else {
				value = 0.0
			}
			keep = true
		}
		if keep {
			lhsSample.V = value
			if shouldDropMetricName(op) || returnBool {
				lhsSample.Metric = enh.dropMetricName(lhsSample.Metric)
			}
			enh.out = append(enh.out, lhsSample)
		}
	}
	return enh.out
}
func dropMetricName(l labels.Labels) labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return labels.NewBuilder(l).Del(labels.MetricName).Labels()
}
func scalarBinop(op ItemType, lhs, rhs float64) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch op {
	case itemADD:
		return lhs + rhs
	case itemSUB:
		return lhs - rhs
	case itemMUL:
		return lhs * rhs
	case itemDIV:
		return lhs / rhs
	case itemPOW:
		return math.Pow(lhs, rhs)
	case itemMOD:
		return math.Mod(lhs, rhs)
	case itemEQL:
		return btos(lhs == rhs)
	case itemNEQ:
		return btos(lhs != rhs)
	case itemGTR:
		return btos(lhs > rhs)
	case itemLSS:
		return btos(lhs < rhs)
	case itemGTE:
		return btos(lhs >= rhs)
	case itemLTE:
		return btos(lhs <= rhs)
	}
	panic(fmt.Errorf("operator %q not allowed for Scalar operations", op))
}
func vectorElemBinop(op ItemType, lhs, rhs float64) (float64, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch op {
	case itemADD:
		return lhs + rhs, true
	case itemSUB:
		return lhs - rhs, true
	case itemMUL:
		return lhs * rhs, true
	case itemDIV:
		return lhs / rhs, true
	case itemPOW:
		return math.Pow(lhs, rhs), true
	case itemMOD:
		return math.Mod(lhs, rhs), true
	case itemEQL:
		return lhs, lhs == rhs
	case itemNEQ:
		return lhs, lhs != rhs
	case itemGTR:
		return lhs, lhs > rhs
	case itemLSS:
		return lhs, lhs < rhs
	case itemGTE:
		return lhs, lhs >= rhs
	case itemLTE:
		return lhs, lhs <= rhs
	}
	panic(fmt.Errorf("operator %q not allowed for operations between Vectors", op))
}

type groupedAggregation struct {
	labels		labels.Labels
	value		float64
	mean		float64
	groupCount	int
	heap		vectorByValueHeap
	reverseHeap	vectorByReverseValueHeap
}

func (ev *evaluator) aggregation(op ItemType, grouping []string, without bool, param interface{}, vec Vector, enh *EvalNodeHelper) Vector {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := map[uint64]*groupedAggregation{}
	var k int64
	if op == itemTopK || op == itemBottomK {
		f := param.(float64)
		if !convertibleToInt64(f) {
			ev.errorf("Scalar value %v overflows int64", f)
		}
		k = int64(f)
		if k < 1 {
			return Vector{}
		}
	}
	var q float64
	if op == itemQuantile {
		q = param.(float64)
	}
	var valueLabel string
	if op == itemCountValues {
		valueLabel = param.(string)
		if !model.LabelName(valueLabel).IsValid() {
			ev.errorf("invalid label name %q", valueLabel)
		}
		if !without {
			grouping = append(grouping, valueLabel)
		}
	}
	for _, s := range vec {
		metric := s.Metric
		if op == itemCountValues {
			lb := labels.NewBuilder(metric)
			lb.Set(valueLabel, strconv.FormatFloat(s.V, 'f', -1, 64))
			metric = lb.Labels()
		}
		var (
			groupingKey uint64
		)
		if without {
			groupingKey = metric.HashWithoutLabels(grouping...)
		} else {
			groupingKey = metric.HashForLabels(grouping...)
		}
		group, ok := result[groupingKey]
		if !ok {
			var m labels.Labels
			if without {
				lb := labels.NewBuilder(metric)
				lb.Del(grouping...)
				lb.Del(labels.MetricName)
				m = lb.Labels()
			} else {
				m = make(labels.Labels, 0, len(grouping))
				for _, l := range metric {
					for _, n := range grouping {
						if l.Name == n {
							m = append(m, l)
							break
						}
					}
				}
				sort.Sort(m)
			}
			result[groupingKey] = &groupedAggregation{labels: m, value: s.V, mean: s.V, groupCount: 1}
			inputVecLen := int64(len(vec))
			resultSize := k
			if k > inputVecLen {
				resultSize = inputVecLen
			}
			if op == itemStdvar || op == itemStddev {
				result[groupingKey].value = 0.0
			} else if op == itemTopK || op == itemQuantile {
				result[groupingKey].heap = make(vectorByValueHeap, 0, resultSize)
				heap.Push(&result[groupingKey].heap, &Sample{Point: Point{V: s.V}, Metric: s.Metric})
			} else if op == itemBottomK {
				result[groupingKey].reverseHeap = make(vectorByReverseValueHeap, 0, resultSize)
				heap.Push(&result[groupingKey].reverseHeap, &Sample{Point: Point{V: s.V}, Metric: s.Metric})
			}
			continue
		}
		switch op {
		case itemSum:
			group.value += s.V
		case itemAvg:
			group.groupCount++
			group.mean += (s.V - group.mean) / float64(group.groupCount)
		case itemMax:
			if group.value < s.V || math.IsNaN(group.value) {
				group.value = s.V
			}
		case itemMin:
			if group.value > s.V || math.IsNaN(group.value) {
				group.value = s.V
			}
		case itemCount, itemCountValues:
			group.groupCount++
		case itemStdvar, itemStddev:
			group.groupCount++
			delta := s.V - group.mean
			group.mean += delta / float64(group.groupCount)
			group.value += delta * (s.V - group.mean)
		case itemTopK:
			if int64(len(group.heap)) < k || group.heap[0].V < s.V || math.IsNaN(group.heap[0].V) {
				if int64(len(group.heap)) == k {
					heap.Pop(&group.heap)
				}
				heap.Push(&group.heap, &Sample{Point: Point{V: s.V}, Metric: s.Metric})
			}
		case itemBottomK:
			if int64(len(group.reverseHeap)) < k || group.reverseHeap[0].V > s.V || math.IsNaN(group.reverseHeap[0].V) {
				if int64(len(group.reverseHeap)) == k {
					heap.Pop(&group.reverseHeap)
				}
				heap.Push(&group.reverseHeap, &Sample{Point: Point{V: s.V}, Metric: s.Metric})
			}
		case itemQuantile:
			group.heap = append(group.heap, s)
		default:
			panic(fmt.Errorf("expected aggregation operator but got %q", op))
		}
	}
	for _, aggr := range result {
		switch op {
		case itemAvg:
			aggr.value = aggr.mean
		case itemCount, itemCountValues:
			aggr.value = float64(aggr.groupCount)
		case itemStdvar:
			aggr.value = aggr.value / float64(aggr.groupCount)
		case itemStddev:
			aggr.value = math.Sqrt(aggr.value / float64(aggr.groupCount))
		case itemTopK:
			sort.Sort(sort.Reverse(aggr.heap))
			for _, v := range aggr.heap {
				enh.out = append(enh.out, Sample{Metric: v.Metric, Point: Point{V: v.V}})
			}
			continue
		case itemBottomK:
			sort.Sort(sort.Reverse(aggr.reverseHeap))
			for _, v := range aggr.reverseHeap {
				enh.out = append(enh.out, Sample{Metric: v.Metric, Point: Point{V: v.V}})
			}
			continue
		case itemQuantile:
			aggr.value = quantile(q, aggr.heap)
		default:
		}
		enh.out = append(enh.out, Sample{Metric: aggr.labels, Point: Point{V: aggr.value}})
	}
	return enh.out
}
func btos(b bool) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if b {
		return 1
	}
	return 0
}
func shouldDropMetricName(op ItemType) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch op {
	case itemADD, itemSUB, itemDIV, itemMUL, itemMOD:
		return true
	default:
		return false
	}
}
func documentedType(t ValueType) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch t {
	case "vector":
		return "instant vector"
	case "matrix":
		return "range vector"
	default:
		return string(t)
	}
}
