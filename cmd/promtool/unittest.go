package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"gopkg.in/yaml.v2"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/rules"
	"github.com/prometheus/prometheus/storage"
)

func RulesUnitTest(files ...string) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	failed := false
	for _, f := range files {
		if errs := ruleUnitTest(f); errs != nil {
			fmt.Fprintln(os.Stderr, "  FAILED:")
			for _, e := range errs {
				fmt.Fprintln(os.Stderr, e.Error())
			}
			failed = true
		} else {
			fmt.Println("  SUCCESS")
		}
		fmt.Println()
	}
	if failed {
		return 1
	}
	return 0
}
func ruleUnitTest(filename string) []error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	fmt.Println("Unit Testing: ", filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return []error{err}
	}
	var unitTestInp unitTestFile
	if err := yaml.UnmarshalStrict(b, &unitTestInp); err != nil {
		return []error{err}
	}
	if unitTestInp.EvaluationInterval == 0 {
		unitTestInp.EvaluationInterval = 1 * time.Minute
	}
	mint := time.Unix(0, 0)
	maxd := unitTestInp.maxEvalTime()
	maxt := mint.Add(maxd)
	maxt = maxt.Add(unitTestInp.EvaluationInterval / 2).Round(unitTestInp.EvaluationInterval)
	groupOrderMap := make(map[string]int)
	for i, gn := range unitTestInp.GroupEvalOrder {
		if _, ok := groupOrderMap[gn]; ok {
			return []error{fmt.Errorf("group name repeated in evaluation order: %s", gn)}
		}
		groupOrderMap[gn] = i
	}
	var errs []error
	for _, t := range unitTestInp.Tests {
		ers := t.test(mint, maxt, unitTestInp.EvaluationInterval, groupOrderMap, unitTestInp.RuleFiles...)
		if ers != nil {
			errs = append(errs, ers...)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type unitTestFile struct {
	RuleFiles		[]string	`yaml:"rule_files"`
	EvaluationInterval	time.Duration	`yaml:"evaluation_interval,omitempty"`
	GroupEvalOrder		[]string	`yaml:"group_eval_order"`
	Tests			[]testGroup	`yaml:"tests"`
}

func (utf *unitTestFile) maxEvalTime() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var maxd time.Duration
	for _, t := range utf.Tests {
		d := t.maxEvalTime()
		if d > maxd {
			maxd = d
		}
	}
	return maxd
}

type testGroup struct {
	Interval	time.Duration		`yaml:"interval"`
	InputSeries	[]series		`yaml:"input_series"`
	AlertRuleTests	[]alertTestCase		`yaml:"alert_rule_test,omitempty"`
	PromqlExprTests	[]promqlTestCase	`yaml:"promql_expr_test,omitempty"`
}

func (tg *testGroup) test(mint, maxt time.Time, evalInterval time.Duration, groupOrderMap map[string]int, ruleFiles ...string) []error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	suite, err := promql.NewLazyLoader(nil, tg.seriesLoadingString())
	if err != nil {
		return []error{err}
	}
	defer suite.Close()
	opts := &rules.ManagerOptions{QueryFunc: rules.EngineQueryFunc(suite.QueryEngine(), suite.Storage()), Appendable: suite.Storage(), Context: context.Background(), NotifyFunc: func(ctx context.Context, expr string, alerts ...*rules.Alert) {
	}, Logger: &dummyLogger{}}
	m := rules.NewManager(opts)
	groupsMap, ers := m.LoadGroups(tg.Interval, ruleFiles...)
	if ers != nil {
		return ers
	}
	groups := orderedGroups(groupsMap, groupOrderMap)
	alertEvalTimesMap := map[time.Duration]struct{}{}
	alertsInTest := make(map[time.Duration]map[string]struct{})
	alertTests := make(map[time.Duration][]alertTestCase)
	for _, alert := range tg.AlertRuleTests {
		alertEvalTimesMap[alert.EvalTime] = struct{}{}
		if _, ok := alertsInTest[alert.EvalTime]; !ok {
			alertsInTest[alert.EvalTime] = make(map[string]struct{})
		}
		alertsInTest[alert.EvalTime][alert.Alertname] = struct{}{}
		alertTests[alert.EvalTime] = append(alertTests[alert.EvalTime], alert)
	}
	alertEvalTimes := make([]time.Duration, 0, len(alertEvalTimesMap))
	for k := range alertEvalTimesMap {
		alertEvalTimes = append(alertEvalTimes, k)
	}
	sort.Slice(alertEvalTimes, func(i, j int) bool {
		return alertEvalTimes[i] < alertEvalTimes[j]
	})
	curr := 0
	var errs []error
	for ts := mint; ts.Before(maxt); ts = ts.Add(evalInterval) {
		suite.WithSamplesTill(ts, func(err error) {
			if err != nil {
				errs = append(errs, err)
				return
			}
			for _, g := range groups {
				g.Eval(suite.Context(), ts)
			}
		})
		if len(errs) > 0 {
			return errs
		}
		for {
			if !(curr < len(alertEvalTimes) && ts.Sub(mint) <= alertEvalTimes[curr] && alertEvalTimes[curr] < ts.Add(evalInterval).Sub(mint)) {
				break
			}
			t := alertEvalTimes[curr]
			presentAlerts := alertsInTest[t]
			got := make(map[string]labelsAndAnnotations)
			for _, g := range groups {
				grules := g.Rules()
				for _, r := range grules {
					ar, ok := r.(*rules.AlertingRule)
					if !ok {
						continue
					}
					if _, ok := presentAlerts[ar.Name()]; !ok {
						continue
					}
					var alerts labelsAndAnnotations
					for _, a := range ar.ActiveAlerts() {
						if a.State == rules.StateFiring {
							alerts = append(alerts, labelAndAnnotation{Labels: append(labels.Labels{}, a.Labels...), Annotations: append(labels.Labels{}, a.Annotations...)})
						}
					}
					got[ar.Name()] = append(got[ar.Name()], alerts...)
				}
			}
			for _, testcase := range alertTests[t] {
				gotAlerts := got[testcase.Alertname]
				var expAlerts labelsAndAnnotations
				for _, a := range testcase.ExpAlerts {
					a.ExpLabels[labels.AlertName] = testcase.Alertname
					expAlerts = append(expAlerts, labelAndAnnotation{Labels: labels.FromMap(a.ExpLabels), Annotations: labels.FromMap(a.ExpAnnotations)})
				}
				if gotAlerts.Len() != expAlerts.Len() {
					errs = append(errs, fmt.Errorf("    alertname:%s, time:%s, \n        exp:%#v, \n        got:%#v", testcase.Alertname, testcase.EvalTime.String(), expAlerts.String(), gotAlerts.String()))
				} else {
					sort.Sort(gotAlerts)
					sort.Sort(expAlerts)
					if !reflect.DeepEqual(expAlerts, gotAlerts) {
						errs = append(errs, fmt.Errorf("    alertname:%s, time:%s, \n        exp:%#v, \n        got:%#v", testcase.Alertname, testcase.EvalTime.String(), expAlerts.String(), gotAlerts.String()))
					}
				}
			}
			curr++
		}
	}
Outer:
	for _, testCase := range tg.PromqlExprTests {
		got, err := query(suite.Context(), testCase.Expr, mint.Add(testCase.EvalTime), suite.QueryEngine(), suite.Queryable())
		if err != nil {
			errs = append(errs, fmt.Errorf("    expr:'%s', time:%s, err:%s", testCase.Expr, testCase.EvalTime.String(), err.Error()))
			continue
		}
		var gotSamples []parsedSample
		for _, s := range got {
			gotSamples = append(gotSamples, parsedSample{Labels: s.Metric.Copy(), Value: s.V})
		}
		var expSamples []parsedSample
		for _, s := range testCase.ExpSamples {
			lb, err := promql.ParseMetric(s.Labels)
			if err != nil {
				errs = append(errs, fmt.Errorf("    expr:'%s', time:%s, err:%s", testCase.Expr, testCase.EvalTime.String(), err.Error()))
				continue Outer
			}
			expSamples = append(expSamples, parsedSample{Labels: lb, Value: s.Value})
		}
		sort.Slice(expSamples, func(i, j int) bool {
			return labels.Compare(expSamples[i].Labels, expSamples[j].Labels) <= 0
		})
		sort.Slice(gotSamples, func(i, j int) bool {
			return labels.Compare(gotSamples[i].Labels, gotSamples[j].Labels) <= 0
		})
		if !reflect.DeepEqual(expSamples, gotSamples) {
			errs = append(errs, fmt.Errorf("    expr:'%s', time:%s, \n        exp:%#v, \n        got:%#v", testCase.Expr, testCase.EvalTime.String(), parsedSamplesString(expSamples), parsedSamplesString(gotSamples)))
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
func (tg *testGroup) seriesLoadingString() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := ""
	result += "load " + shortDuration(tg.Interval) + "\n"
	for _, is := range tg.InputSeries {
		result += "  " + is.Series + " " + is.Values + "\n"
	}
	return result
}
func shortDuration(d time.Duration) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}
func orderedGroups(groupsMap map[string]*rules.Group, groupOrderMap map[string]int) []*rules.Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	groups := make([]*rules.Group, 0, len(groupsMap))
	for _, g := range groupsMap {
		groups = append(groups, g)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groupOrderMap[groups[i].Name()] < groupOrderMap[groups[j].Name()]
	})
	return groups
}
func (tg *testGroup) maxEvalTime() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var maxd time.Duration
	for _, alert := range tg.AlertRuleTests {
		if alert.EvalTime > maxd {
			maxd = alert.EvalTime
		}
	}
	for _, pet := range tg.PromqlExprTests {
		if pet.EvalTime > maxd {
			maxd = pet.EvalTime
		}
	}
	return maxd
}
func query(ctx context.Context, qs string, t time.Time, engine *promql.Engine, qu storage.Queryable) (promql.Vector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	q, err := engine.NewInstantQuery(qu, qs, t)
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

type labelsAndAnnotations []labelAndAnnotation

func (la labelsAndAnnotations) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(la)
}
func (la labelsAndAnnotations) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	la[i], la[j] = la[j], la[i]
}
func (la labelsAndAnnotations) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	diff := labels.Compare(la[i].Labels, la[j].Labels)
	if diff != 0 {
		return diff < 0
	}
	return labels.Compare(la[i].Annotations, la[j].Annotations) < 0
}
func (la labelsAndAnnotations) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(la) == 0 {
		return "[]"
	}
	s := "[" + la[0].String()
	for _, l := range la[1:] {
		s += ", " + l.String()
	}
	s += "]"
	return s
}

type labelAndAnnotation struct {
	Labels		labels.Labels
	Annotations	labels.Labels
}

func (la *labelAndAnnotation) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "Labels:" + la.Labels.String() + " Annotations:" + la.Annotations.String()
}

type series struct {
	Series	string	`yaml:"series"`
	Values	string	`yaml:"values"`
}
type alertTestCase struct {
	EvalTime	time.Duration	`yaml:"eval_time"`
	Alertname	string		`yaml:"alertname"`
	ExpAlerts	[]alert		`yaml:"exp_alerts"`
}
type alert struct {
	ExpLabels	map[string]string	`yaml:"exp_labels"`
	ExpAnnotations	map[string]string	`yaml:"exp_annotations"`
}
type promqlTestCase struct {
	Expr		string		`yaml:"expr"`
	EvalTime	time.Duration	`yaml:"eval_time"`
	ExpSamples	[]sample	`yaml:"exp_samples"`
}
type sample struct {
	Labels	string	`yaml:"labels"`
	Value	float64	`yaml:"value"`
}
type parsedSample struct {
	Labels	labels.Labels
	Value	float64
}

func parsedSamplesString(pss []parsedSample) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(pss) == 0 {
		return "nil"
	}
	s := pss[0].String()
	for _, ps := range pss[0:] {
		s += ", " + ps.String()
	}
	return s
}
func (ps *parsedSample) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ps.Labels.String() + " " + strconv.FormatFloat(ps.Value, 'E', -1, 64)
}

type dummyLogger struct{}

func (l *dummyLogger) Log(keyvals ...interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
