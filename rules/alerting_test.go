package rules

import (
	"testing"
	"time"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/timestamp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/util/testutil"
)

func TestAlertingRuleHTMLSnippet(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	expr, err := promql.ParseExpr(`foo{html="<b>BOLD<b>"}`)
	testutil.Ok(t, err)
	rule := NewAlertingRule("testrule", expr, 0, labels.FromStrings("html", "<b>BOLD</b>"), labels.FromStrings("html", "<b>BOLD</b>"), false, nil)
	const want = `alert: <a href="/test/prefix/graph?g0.expr=ALERTS%7Balertname%3D%22testrule%22%7D&g0.tab=1">testrule</a>
expr: <a href="/test/prefix/graph?g0.expr=foo%7Bhtml%3D%22%3Cb%3EBOLD%3Cb%3E%22%7D&g0.tab=1">foo{html=&#34;&lt;b&gt;BOLD&lt;b&gt;&#34;}</a>
labels:
  html: '&lt;b&gt;BOLD&lt;/b&gt;'
annotations:
  html: '&lt;b&gt;BOLD&lt;/b&gt;'
`
	got := rule.HTMLSnippet("/test/prefix")
	testutil.Assert(t, want == got, "incorrect HTML snippet; want:\n\n|%v|\n\ngot:\n\n|%v|", want, got)
}
func TestAlertingRuleLabelsUpdate(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	suite, err := promql.NewTest(t, `
		load 1m
			http_requests{job="app-server", instance="0"}	75 85 70 70
	`)
	testutil.Ok(t, err)
	defer suite.Close()
	err = suite.Run()
	testutil.Ok(t, err)
	expr, err := promql.ParseExpr(`http_requests < 100`)
	testutil.Ok(t, err)
	rule := NewAlertingRule("HTTPRequestRateLow", expr, time.Minute, labels.FromStrings("severity", "{{ if lt $value 80.0 }}critical{{ else }}warning{{ end }}"), nil, true, nil)
	results := []promql.Vector{{{Metric: labels.FromStrings("__name__", "ALERTS", "alertname", "HTTPRequestRateLow", "alertstate", "pending", "instance", "0", "job", "app-server", "severity", "critical"), Point: promql.Point{V: 1}}}, {{Metric: labels.FromStrings("__name__", "ALERTS", "alertname", "HTTPRequestRateLow", "alertstate", "pending", "instance", "0", "job", "app-server", "severity", "warning"), Point: promql.Point{V: 1}}}, {{Metric: labels.FromStrings("__name__", "ALERTS", "alertname", "HTTPRequestRateLow", "alertstate", "pending", "instance", "0", "job", "app-server", "severity", "critical"), Point: promql.Point{V: 1}}}, {{Metric: labels.FromStrings("__name__", "ALERTS", "alertname", "HTTPRequestRateLow", "alertstate", "firing", "instance", "0", "job", "app-server", "severity", "critical"), Point: promql.Point{V: 1}}}}
	baseTime := time.Unix(0, 0)
	for i, result := range results {
		t.Logf("case %d", i)
		evalTime := baseTime.Add(time.Duration(i) * time.Minute)
		result[0].Point.T = timestamp.FromTime(evalTime)
		res, err := rule.Eval(suite.Context(), evalTime, EngineQueryFunc(suite.QueryEngine(), suite.Storage()), nil)
		testutil.Ok(t, err)
		var filteredRes promql.Vector
		for _, smpl := range res {
			smplName := smpl.Metric.Get("__name__")
			if smplName == "ALERTS" {
				filteredRes = append(filteredRes, smpl)
			} else {
				testutil.Equals(t, smplName, "ALERTS_FOR_STATE")
			}
		}
		testutil.Equals(t, result, filteredRes)
	}
}
