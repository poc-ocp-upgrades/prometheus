package rules

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"sort"
	"sync"
	"time"
	html_template "html/template"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/prometheus/prometheus/pkg/timestamp"
	"github.com/prometheus/prometheus/pkg/value"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
)

type RuleHealth string

const (
	HealthUnknown	RuleHealth	= "unknown"
	HealthGood		RuleHealth	= "ok"
	HealthBad		RuleHealth	= "err"
)
const namespace = "prometheus"

var (
	groupInterval = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "rule_group_interval_seconds"), "The interval of a rule group.", []string{"rule_group"}, nil)
)

type Metrics struct {
	evalDuration		prometheus.Summary
	evalFailures		prometheus.Counter
	evalTotal			prometheus.Counter
	iterationDuration	prometheus.Summary
	iterationsMissed	prometheus.Counter
	iterationsScheduled	prometheus.Counter
	groupLastEvalTime	*prometheus.GaugeVec
	groupLastDuration	*prometheus.GaugeVec
	groupRules			*prometheus.GaugeVec
}

func NewGroupMetrics(reg prometheus.Registerer) *Metrics {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := &Metrics{evalDuration: prometheus.NewSummary(prometheus.SummaryOpts{Namespace: namespace, Name: "rule_evaluation_duration_seconds", Help: "The duration for a rule to execute."}), evalFailures: prometheus.NewCounter(prometheus.CounterOpts{Namespace: namespace, Name: "rule_evaluation_failures_total", Help: "The total number of rule evaluation failures."}), evalTotal: prometheus.NewCounter(prometheus.CounterOpts{Namespace: namespace, Name: "rule_evaluations_total", Help: "The total number of rule evaluations."}), iterationDuration: prometheus.NewSummary(prometheus.SummaryOpts{Namespace: namespace, Name: "rule_group_duration_seconds", Help: "The duration of rule group evaluations.", Objectives: map[float64]float64{0.01: 0.001, 0.05: 0.005, 0.5: 0.05, 0.90: 0.01, 0.99: 0.001}}), iterationsMissed: prometheus.NewCounter(prometheus.CounterOpts{Namespace: namespace, Name: "rule_group_iterations_missed_total", Help: "The total number of rule group evaluations missed due to slow rule group evaluation."}), iterationsScheduled: prometheus.NewCounter(prometheus.CounterOpts{Namespace: namespace, Name: "rule_group_iterations_total", Help: "The total number of scheduled rule group evaluations, whether executed or missed."}), groupLastEvalTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Name: "rule_group_last_evaluation_timestamp_seconds", Help: "The timestamp of the last rule group evaluation in seconds."}, []string{"rule_group"}), groupLastDuration: prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Name: "rule_group_last_duration_seconds", Help: "The duration of the last rule group evaluation."}, []string{"rule_group"}), groupRules: prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Name: "rule_group_rules", Help: "The number of rules."}, []string{"rule_group"})}
	if reg != nil {
		reg.MustRegister(m.evalDuration, m.evalFailures, m.evalTotal, m.iterationDuration, m.iterationsMissed, m.iterationsScheduled, m.groupLastEvalTime, m.groupLastDuration, m.groupRules)
	}
	return m
}

type QueryFunc func(ctx context.Context, q string, t time.Time) (promql.Vector, error)

func EngineQueryFunc(engine *promql.Engine, q storage.Queryable) QueryFunc {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(ctx context.Context, qs string, t time.Time) (promql.Vector, error) {
		q, err := engine.NewInstantQuery(q, qs, t)
		if err != nil {
			return nil, err
		}
		res := q.Exec(ctx)
		if res.Err != nil {
			return nil, res.Err
		}
		switch v := res.Value.(type) {
		case promql.Vector:
			return v, nil
		case promql.Scalar:
			return promql.Vector{promql.Sample{Point: promql.Point(v), Metric: labels.Labels{}}}, nil
		default:
			return nil, fmt.Errorf("rule result is not a vector or scalar")
		}
	}
}

type Rule interface {
	Name() string
	Eval(context.Context, time.Time, QueryFunc, *url.URL) (promql.Vector, error)
	String() string
	SetLastError(error)
	LastError() error
	SetHealth(RuleHealth)
	Health() RuleHealth
	SetEvaluationDuration(time.Duration)
	GetEvaluationDuration() time.Duration
	SetEvaluationTimestamp(time.Time)
	GetEvaluationTimestamp() time.Time
	HTMLSnippet(pathPrefix string) html_template.HTML
}
type Group struct {
	name					string
	file					string
	interval				time.Duration
	rules					[]Rule
	seriesInPreviousEval	[]map[string]labels.Labels
	opts					*ManagerOptions
	mtx						sync.Mutex
	evaluationDuration		time.Duration
	evaluationTimestamp		time.Time
	shouldRestore			bool
	done					chan struct{}
	terminated				chan struct{}
	logger					log.Logger
	metrics					*Metrics
}

func NewGroup(name, file string, interval time.Duration, rules []Rule, shouldRestore bool, opts *ManagerOptions) *Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	metrics := opts.Metrics
	if metrics == nil {
		metrics = NewGroupMetrics(opts.Registerer)
	}
	metrics.groupLastEvalTime.WithLabelValues(groupKey(file, name))
	metrics.groupLastDuration.WithLabelValues(groupKey(file, name))
	metrics.groupRules.WithLabelValues(groupKey(file, name)).Set(float64(len(rules)))
	return &Group{name: name, file: file, interval: interval, rules: rules, shouldRestore: shouldRestore, opts: opts, seriesInPreviousEval: make([]map[string]labels.Labels, len(rules)), done: make(chan struct{}), terminated: make(chan struct{}), logger: log.With(opts.Logger, "group", name), metrics: metrics}
}
func (g *Group) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return g.name
}
func (g *Group) File() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return g.file
}
func (g *Group) Rules() []Rule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return g.rules
}
func (g *Group) Interval() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return g.interval
}
func (g *Group) run(ctx context.Context) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer close(g.terminated)
	evalTimestamp := g.evalTimestamp().Add(g.interval)
	select {
	case <-time.After(time.Until(evalTimestamp)):
	case <-g.done:
		return
	}
	iter := func() {
		g.metrics.iterationsScheduled.Inc()
		start := time.Now()
		g.Eval(ctx, evalTimestamp)
		timeSinceStart := time.Since(start)
		g.metrics.iterationDuration.Observe(timeSinceStart.Seconds())
		g.setEvaluationDuration(timeSinceStart)
		g.setEvaluationTimestamp(start)
	}
	tick := time.NewTicker(g.interval)
	defer tick.Stop()
	iter()
	if g.shouldRestore {
		select {
		case <-g.done:
			return
		case <-tick.C:
			missed := (time.Since(evalTimestamp) / g.interval) - 1
			if missed > 0 {
				g.metrics.iterationsMissed.Add(float64(missed))
				g.metrics.iterationsScheduled.Add(float64(missed))
			}
			evalTimestamp = evalTimestamp.Add((missed + 1) * g.interval)
			iter()
		}
		g.RestoreForState(time.Now())
		g.shouldRestore = false
	}
	for {
		select {
		case <-g.done:
			return
		default:
			select {
			case <-g.done:
				return
			case <-tick.C:
				missed := (time.Since(evalTimestamp) / g.interval) - 1
				if missed > 0 {
					g.metrics.iterationsMissed.Add(float64(missed))
					g.metrics.iterationsScheduled.Add(float64(missed))
				}
				evalTimestamp = evalTimestamp.Add((missed + 1) * g.interval)
				iter()
			}
		}
	}
}
func (g *Group) stop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	close(g.done)
	<-g.terminated
}
func (g *Group) hash() uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := labels.New(labels.Label{"name", g.name}, labels.Label{"file", g.file})
	return l.Hash()
}
func (g *Group) GetEvaluationDuration() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.evaluationDuration
}
func (g *Group) setEvaluationDuration(dur time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	g.metrics.groupLastDuration.WithLabelValues(groupKey(g.file, g.name)).Set(dur.Seconds())
	g.mtx.Lock()
	defer g.mtx.Unlock()
	g.evaluationDuration = dur
}
func (g *Group) GetEvaluationTimestamp() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.evaluationTimestamp
}
func (g *Group) setEvaluationTimestamp(ts time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	g.metrics.groupLastEvalTime.WithLabelValues(groupKey(g.file, g.name)).Set(float64(ts.UnixNano()) / 1e9)
	g.mtx.Lock()
	defer g.mtx.Unlock()
	g.evaluationTimestamp = ts
}
func (g *Group) evalTimestamp() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		offset	= int64(g.hash() % uint64(g.interval))
		now		= time.Now().UnixNano()
		adjNow	= now - offset
		base	= adjNow - (adjNow % int64(g.interval))
	)
	return time.Unix(0, base+offset)
}
func (g *Group) CopyState(from *Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	g.evaluationDuration = from.evaluationDuration
	ruleMap := make(map[string][]int, len(from.rules))
	for fi, fromRule := range from.rules {
		l := ruleMap[fromRule.Name()]
		ruleMap[fromRule.Name()] = append(l, fi)
	}
	for i, rule := range g.rules {
		indexes := ruleMap[rule.Name()]
		if len(indexes) == 0 {
			continue
		}
		fi := indexes[0]
		g.seriesInPreviousEval[i] = from.seriesInPreviousEval[fi]
		ruleMap[rule.Name()] = indexes[1:]
		ar, ok := rule.(*AlertingRule)
		if !ok {
			continue
		}
		far, ok := from.rules[fi].(*AlertingRule)
		if !ok {
			continue
		}
		for fp, a := range far.active {
			ar.active[fp] = a
		}
	}
}
func (g *Group) Eval(ctx context.Context, ts time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i, rule := range g.rules {
		select {
		case <-g.done:
			return
		default:
		}
		func(i int, rule Rule) {
			sp, ctx := opentracing.StartSpanFromContext(ctx, "rule")
			sp.SetTag("name", rule.Name())
			defer func(t time.Time) {
				sp.Finish()
				since := time.Since(t)
				g.metrics.evalDuration.Observe(since.Seconds())
				rule.SetEvaluationDuration(since)
				rule.SetEvaluationTimestamp(t)
			}(time.Now())
			g.metrics.evalTotal.Inc()
			vector, err := rule.Eval(ctx, ts, g.opts.QueryFunc, g.opts.ExternalURL)
			if err != nil {
				if _, ok := err.(promql.ErrQueryCanceled); !ok {
					level.Warn(g.logger).Log("msg", "Evaluating rule failed", "rule", rule, "err", err)
				}
				g.metrics.evalFailures.Inc()
				return
			}
			if ar, ok := rule.(*AlertingRule); ok {
				ar.sendAlerts(ctx, ts, g.opts.ResendDelay, g.interval, g.opts.NotifyFunc)
			}
			var (
				numOutOfOrder	= 0
				numDuplicates	= 0
			)
			app, err := g.opts.Appendable.Appender()
			if err != nil {
				level.Warn(g.logger).Log("msg", "creating appender failed", "err", err)
				return
			}
			seriesReturned := make(map[string]labels.Labels, len(g.seriesInPreviousEval[i]))
			for _, s := range vector {
				if _, err := app.Add(s.Metric, s.T, s.V); err != nil {
					switch err {
					case storage.ErrOutOfOrderSample:
						numOutOfOrder++
						level.Debug(g.logger).Log("msg", "Rule evaluation result discarded", "err", err, "sample", s)
					case storage.ErrDuplicateSampleForTimestamp:
						numDuplicates++
						level.Debug(g.logger).Log("msg", "Rule evaluation result discarded", "err", err, "sample", s)
					default:
						level.Warn(g.logger).Log("msg", "Rule evaluation result discarded", "err", err, "sample", s)
					}
				} else {
					seriesReturned[s.Metric.String()] = s.Metric
				}
			}
			if numOutOfOrder > 0 {
				level.Warn(g.logger).Log("msg", "Error on ingesting out-of-order result from rule evaluation", "numDropped", numOutOfOrder)
			}
			if numDuplicates > 0 {
				level.Warn(g.logger).Log("msg", "Error on ingesting results from rule evaluation with different value but same timestamp", "numDropped", numDuplicates)
			}
			for metric, lset := range g.seriesInPreviousEval[i] {
				if _, ok := seriesReturned[metric]; !ok {
					_, err = app.Add(lset, timestamp.FromTime(ts), math.Float64frombits(value.StaleNaN))
					switch err {
					case nil:
					case storage.ErrOutOfOrderSample, storage.ErrDuplicateSampleForTimestamp:
					default:
						level.Warn(g.logger).Log("msg", "adding stale sample failed", "sample", metric, "err", err)
					}
				}
			}
			if err := app.Commit(); err != nil {
				level.Warn(g.logger).Log("msg", "rule sample appending failed", "err", err)
			} else {
				g.seriesInPreviousEval[i] = seriesReturned
			}
		}(i, rule)
	}
}
func (g *Group) RestoreForState(ts time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	maxtMS := int64(model.TimeFromUnixNano(ts.UnixNano()))
	mint := ts.Add(-g.opts.OutageTolerance)
	mintMS := int64(model.TimeFromUnixNano(mint.UnixNano()))
	q, err := g.opts.TSDB.Querier(g.opts.Context, mintMS, maxtMS)
	if err != nil {
		level.Error(g.logger).Log("msg", "Failed to get Querier", "err", err)
		return
	}
	defer func() {
		if err := q.Close(); err != nil {
			level.Error(g.logger).Log("msg", "Failed to close Querier", "err", err)
		}
	}()
	for _, rule := range g.Rules() {
		alertRule, ok := rule.(*AlertingRule)
		if !ok {
			continue
		}
		alertHoldDuration := alertRule.HoldDuration()
		if alertHoldDuration < g.opts.ForGracePeriod {
			alertRule.SetRestored(true)
			continue
		}
		alertRule.ForEachActiveAlert(func(a *Alert) {
			smpl := alertRule.forStateSample(a, time.Now(), 0)
			var matchers []*labels.Matcher
			for _, l := range smpl.Metric {
				mt, err := labels.NewMatcher(labels.MatchEqual, l.Name, l.Value)
				if err != nil {
					panic(err)
				}
				matchers = append(matchers, mt)
			}
			sset, err, _ := q.Select(nil, matchers...)
			if err != nil {
				level.Error(g.logger).Log("msg", "Failed to restore 'for' state", labels.AlertName, alertRule.Name(), "stage", "Select", "err", err)
				return
			}
			seriesFound := false
			var s storage.Series
			for sset.Next() {
				if len(sset.At().Labels()) == len(smpl.Metric) {
					s = sset.At()
					seriesFound = true
					break
				}
			}
			if !seriesFound {
				return
			}
			var t int64
			var v float64
			it := s.Iterator()
			for it.Next() {
				t, v = it.At()
			}
			if it.Err() != nil {
				level.Error(g.logger).Log("msg", "Failed to restore 'for' state", labels.AlertName, alertRule.Name(), "stage", "Iterator", "err", it.Err())
				return
			}
			if value.IsStaleNaN(v) {
				return
			}
			downAt := time.Unix(t/1000, 0)
			restoredActiveAt := time.Unix(int64(v), 0)
			timeSpentPending := downAt.Sub(restoredActiveAt)
			timeRemainingPending := alertHoldDuration - timeSpentPending
			if timeRemainingPending <= 0 {
			} else if timeRemainingPending < g.opts.ForGracePeriod {
				restoredActiveAt = ts.Add(g.opts.ForGracePeriod).Add(-alertHoldDuration)
			} else {
				downDuration := ts.Sub(downAt)
				restoredActiveAt = restoredActiveAt.Add(downDuration)
			}
			a.ActiveAt = restoredActiveAt
			level.Debug(g.logger).Log("msg", "'for' state restored", labels.AlertName, alertRule.Name(), "restored_time", a.ActiveAt.Format(time.RFC850), "labels", a.Labels.String())
		})
		alertRule.SetRestored(true)
	}
}

type Manager struct {
	opts		*ManagerOptions
	groups		map[string]*Group
	mtx			sync.RWMutex
	block		chan struct{}
	restored	bool
	logger		log.Logger
}
type Appendable interface {
	Appender() (storage.Appender, error)
}
type NotifyFunc func(ctx context.Context, expr string, alerts ...*Alert)
type ManagerOptions struct {
	ExternalURL		*url.URL
	QueryFunc		QueryFunc
	NotifyFunc		NotifyFunc
	Context			context.Context
	Appendable		Appendable
	TSDB			storage.Storage
	Logger			log.Logger
	Registerer		prometheus.Registerer
	OutageTolerance	time.Duration
	ForGracePeriod	time.Duration
	ResendDelay		time.Duration
	Metrics			*Metrics
}

func NewManager(o *ManagerOptions) *Manager {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if o.Metrics == nil {
		o.Metrics = NewGroupMetrics(o.Registerer)
	}
	m := &Manager{groups: map[string]*Group{}, opts: o, block: make(chan struct{}), logger: o.Logger}
	if o.Registerer != nil {
		o.Registerer.MustRegister(m)
	}
	o.Metrics.iterationsMissed.Inc()
	return m
}
func (m *Manager) Run() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	close(m.block)
}
func (m *Manager) Stop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtx.Lock()
	defer m.mtx.Unlock()
	level.Info(m.logger).Log("msg", "Stopping rule manager...")
	for _, eg := range m.groups {
		eg.stop()
	}
	level.Info(m.logger).Log("msg", "Rule manager stopped")
}
func (m *Manager) Update(interval time.Duration, files []string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtx.Lock()
	defer m.mtx.Unlock()
	groups, errs := m.LoadGroups(interval, files...)
	if errs != nil {
		for _, e := range errs {
			level.Error(m.logger).Log("msg", "loading groups failed", "err", e)
		}
		return errors.New("error loading rules, previous rule set restored")
	}
	m.restored = true
	var wg sync.WaitGroup
	for _, newg := range groups {
		wg.Add(1)
		gn := groupKey(newg.name, newg.file)
		oldg, ok := m.groups[gn]
		delete(m.groups, gn)
		go func(newg *Group) {
			if ok {
				oldg.stop()
				newg.CopyState(oldg)
			}
			go func() {
				<-m.block
				newg.run(m.opts.Context)
			}()
			wg.Done()
		}(newg)
	}
	for _, oldg := range m.groups {
		oldg.stop()
	}
	wg.Wait()
	m.groups = groups
	return nil
}
func (m *Manager) LoadGroups(interval time.Duration, filenames ...string) (map[string]*Group, []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	groups := make(map[string]*Group)
	shouldRestore := !m.restored
	for _, fn := range filenames {
		rgs, errs := rulefmt.ParseFile(fn)
		if errs != nil {
			return nil, errs
		}
		for _, rg := range rgs.Groups {
			itv := interval
			if rg.Interval != 0 {
				itv = time.Duration(rg.Interval)
			}
			rules := make([]Rule, 0, len(rg.Rules))
			for _, r := range rg.Rules {
				expr, err := promql.ParseExpr(r.Expr)
				if err != nil {
					return nil, []error{err}
				}
				if r.Alert != "" {
					rules = append(rules, NewAlertingRule(r.Alert, expr, time.Duration(r.For), labels.FromMap(r.Labels), labels.FromMap(r.Annotations), m.restored, log.With(m.logger, "alert", r.Alert)))
					continue
				}
				rules = append(rules, NewRecordingRule(r.Record, expr, labels.FromMap(r.Labels)))
			}
			groups[groupKey(rg.Name, fn)] = NewGroup(rg.Name, fn, itv, rules, shouldRestore, m.opts)
		}
	}
	return groups, nil
}
func groupKey(name, file string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return name + ";" + file
}
func (m *Manager) RuleGroups() []*Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	rgs := make([]*Group, 0, len(m.groups))
	for _, g := range m.groups {
		rgs = append(rgs, g)
	}
	sort.Slice(rgs, func(i, j int) bool {
		return rgs[i].file < rgs[j].file && rgs[i].name < rgs[j].name
	})
	return rgs
}
func (m *Manager) Rules() []Rule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	var rules []Rule
	for _, g := range m.groups {
		rules = append(rules, g.rules...)
	}
	return rules
}
func (m *Manager) AlertingRules() []*AlertingRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	alerts := []*AlertingRule{}
	for _, rule := range m.Rules() {
		if alertingRule, ok := rule.(*AlertingRule); ok {
			alerts = append(alerts, alertingRule)
		}
	}
	return alerts
}
func (m *Manager) Describe(ch chan<- *prometheus.Desc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- groupInterval
}
func (m *Manager) Collect(ch chan<- prometheus.Metric) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, g := range m.RuleGroups() {
		ch <- prometheus.MustNewConstMetric(groupInterval, prometheus.GaugeValue, g.interval.Seconds(), groupKey(g.file, g.name))
	}
}
