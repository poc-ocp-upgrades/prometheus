package rules

import (
	"context"
	"testing"
	"time"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/timestamp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/util/testutil"
)

func TestRuleEval(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	storage := testutil.NewStorage(t)
	defer storage.Close()
	opts := promql.EngineOpts{Logger: nil, Reg: nil, MaxConcurrent: 10, MaxSamples: 10, Timeout: 10 * time.Second}
	engine := promql.NewEngine(opts)
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	now := time.Now()
	suite := []struct {
		name	string
		expr	promql.Expr
		labels	labels.Labels
		result	promql.Vector
	}{{name: "nolabels", expr: &promql.NumberLiteral{Val: 1}, labels: labels.Labels{}, result: promql.Vector{promql.Sample{Metric: labels.FromStrings("__name__", "nolabels"), Point: promql.Point{V: 1, T: timestamp.FromTime(now)}}}}, {name: "labels", expr: &promql.NumberLiteral{Val: 1}, labels: labels.FromStrings("foo", "bar"), result: promql.Vector{promql.Sample{Metric: labels.FromStrings("__name__", "labels", "foo", "bar"), Point: promql.Point{V: 1, T: timestamp.FromTime(now)}}}}}
	for _, test := range suite {
		rule := NewRecordingRule(test.name, test.expr, test.labels)
		result, err := rule.Eval(ctx, now, EngineQueryFunc(engine, storage), nil)
		testutil.Ok(t, err)
		testutil.Equals(t, result, test.result)
	}
}
func TestRecordingRuleHTMLSnippet(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	expr, err := promql.ParseExpr(`foo{html="<b>BOLD<b>"}`)
	testutil.Ok(t, err)
	rule := NewRecordingRule("testrule", expr, labels.FromStrings("html", "<b>BOLD</b>"))
	const want = `record: <a href="/test/prefix/graph?g0.expr=testrule&g0.tab=1">testrule</a>
expr: <a href="/test/prefix/graph?g0.expr=foo%7Bhtml%3D%22%3Cb%3EBOLD%3Cb%3E%22%7D&g0.tab=1">foo{html=&#34;&lt;b&gt;BOLD&lt;b&gt;&#34;}</a>
labels:
  html: '&lt;b&gt;BOLD&lt;/b&gt;'
`
	got := rule.HTMLSnippet("/test/prefix")
	testutil.Assert(t, want == got, "incorrect HTML snippet; want:\n\n%s\n\ngot:\n\n%s", want, got)
}
