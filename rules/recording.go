package rules

import (
	"context"
	"fmt"
	"html/template"
	"net/url"
	"sync"
	"time"
	yaml "gopkg.in/yaml.v2"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/util/strutil"
)

type RecordingRule struct {
	name				string
	vector				promql.Expr
	labels				labels.Labels
	mtx					sync.Mutex
	health				RuleHealth
	evaluationTimestamp	time.Time
	lastError			error
	evaluationDuration	time.Duration
}

func NewRecordingRule(name string, vector promql.Expr, lset labels.Labels) *RecordingRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &RecordingRule{name: name, vector: vector, health: HealthUnknown, labels: lset}
}
func (rule *RecordingRule) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.name
}
func (rule *RecordingRule) Query() promql.Expr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.vector
}
func (rule *RecordingRule) Labels() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.labels
}
func (rule *RecordingRule) Eval(ctx context.Context, ts time.Time, query QueryFunc, _ *url.URL) (promql.Vector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	vector, err := query(ctx, rule.vector.String(), ts)
	if err != nil {
		rule.SetHealth(HealthBad)
		rule.SetLastError(err)
		return nil, err
	}
	for i := range vector {
		sample := &vector[i]
		lb := labels.NewBuilder(sample.Metric)
		lb.Set(labels.MetricName, rule.name)
		for _, l := range rule.labels {
			if l.Value == "" {
				lb.Del(l.Name)
			} else {
				lb.Set(l.Name, l.Value)
			}
		}
		sample.Metric = lb.Labels()
	}
	rule.SetHealth(HealthGood)
	rule.SetLastError(err)
	return vector, nil
}
func (rule *RecordingRule) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r := rulefmt.Rule{Record: rule.name, Expr: rule.vector.String(), Labels: rule.labels.Map()}
	byt, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Sprintf("error marshaling recording rule: %q", err.Error())
	}
	return string(byt)
}
func (rule *RecordingRule) SetEvaluationDuration(dur time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rule.mtx.Lock()
	defer rule.mtx.Unlock()
	rule.evaluationDuration = dur
}
func (rule *RecordingRule) SetLastError(err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rule.mtx.Lock()
	defer rule.mtx.Unlock()
	rule.lastError = err
}
func (rule *RecordingRule) LastError() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rule.mtx.Lock()
	defer rule.mtx.Unlock()
	return rule.lastError
}
func (rule *RecordingRule) SetHealth(health RuleHealth) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rule.mtx.Lock()
	defer rule.mtx.Unlock()
	rule.health = health
}
func (rule *RecordingRule) Health() RuleHealth {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rule.mtx.Lock()
	defer rule.mtx.Unlock()
	return rule.health
}
func (rule *RecordingRule) GetEvaluationDuration() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rule.mtx.Lock()
	defer rule.mtx.Unlock()
	return rule.evaluationDuration
}
func (rule *RecordingRule) SetEvaluationTimestamp(ts time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rule.mtx.Lock()
	defer rule.mtx.Unlock()
	rule.evaluationTimestamp = ts
}
func (rule *RecordingRule) GetEvaluationTimestamp() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rule.mtx.Lock()
	defer rule.mtx.Unlock()
	return rule.evaluationTimestamp
}
func (rule *RecordingRule) HTMLSnippet(pathPrefix string) template.HTML {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ruleExpr := rule.vector.String()
	labels := make(map[string]string, len(rule.labels))
	for _, l := range rule.labels {
		labels[l.Name] = template.HTMLEscapeString(l.Value)
	}
	r := rulefmt.Rule{Record: fmt.Sprintf(`<a href="%s">%s</a>`, pathPrefix+strutil.TableLinkForExpression(rule.name), rule.name), Expr: fmt.Sprintf(`<a href="%s">%s</a>`, pathPrefix+strutil.TableLinkForExpression(ruleExpr), template.HTMLEscapeString(ruleExpr)), Labels: labels}
	byt, err := yaml.Marshal(r)
	if err != nil {
		return template.HTML(fmt.Sprintf("error marshaling recording rule: %q", template.HTMLEscapeString(err.Error())))
	}
	return template.HTML(byt)
}
