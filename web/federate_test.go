package web

import (
	"bytes"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/promql"
)

var scenarios = map[string]struct {
	params		string
	accept		string
	externalLabels	model.LabelSet
	code		int
	body		string
}{"empty": {params: "", code: 200, body: ``}, "match nothing": {params: "match[]=does_not_match_anything", code: 200, body: ``}, "invalid params from the beginning": {params: "match[]=-not-a-valid-metric-name", code: 400, body: `parse error at char 1: vector selector must contain label matchers or metric name
`}, "invalid params somewhere in the middle": {params: "match[]=not-a-valid-metric-name", code: 400, body: `parse error at char 4: could not parse remaining input "-a-valid-metric"...
`}, "test_metric1": {params: "match[]=test_metric1", code: 200, body: `# TYPE test_metric1 untyped
test_metric1{foo="bar",instance="i"} 10000 6000000
test_metric1{foo="boo",instance="i"} 1 6000000
`}, "test_metric2": {params: "match[]=test_metric2", code: 200, body: `# TYPE test_metric2 untyped
test_metric2{foo="boo",instance="i"} 1 6000000
`}, "test_metric_without_labels": {params: "match[]=test_metric_without_labels", code: 200, body: `# TYPE test_metric_without_labels untyped
test_metric_without_labels{instance=""} 1001 6000000
`}, "test_stale_metric": {params: "match[]=test_metric_stale", code: 200, body: ``}, "test_old_metric": {params: "match[]=test_metric_old", code: 200, body: `# TYPE test_metric_old untyped
test_metric_old{instance=""} 981 5880000
`}, "{foo='boo'}": {params: "match[]={foo='boo'}", code: 200, body: `# TYPE test_metric1 untyped
test_metric1{foo="boo",instance="i"} 1 6000000
# TYPE test_metric2 untyped
test_metric2{foo="boo",instance="i"} 1 6000000
`}, "two matchers": {params: "match[]=test_metric1&match[]=test_metric2", code: 200, body: `# TYPE test_metric1 untyped
test_metric1{foo="bar",instance="i"} 10000 6000000
test_metric1{foo="boo",instance="i"} 1 6000000
# TYPE test_metric2 untyped
test_metric2{foo="boo",instance="i"} 1 6000000
`}, "everything": {params: "match[]={__name__=~'.%2b'}", code: 200, body: `# TYPE test_metric1 untyped
test_metric1{foo="bar",instance="i"} 10000 6000000
test_metric1{foo="boo",instance="i"} 1 6000000
# TYPE test_metric2 untyped
test_metric2{foo="boo",instance="i"} 1 6000000
# TYPE test_metric_old untyped
test_metric_old{instance=""} 981 5880000
# TYPE test_metric_without_labels untyped
test_metric_without_labels{instance=""} 1001 6000000
`}, "empty label value matches everything that doesn't have that label": {params: "match[]={foo='',__name__=~'.%2b'}", code: 200, body: `# TYPE test_metric_old untyped
test_metric_old{instance=""} 981 5880000
# TYPE test_metric_without_labels untyped
test_metric_without_labels{instance=""} 1001 6000000
`}, "empty label value for a label that doesn't exist at all, matches everything": {params: "match[]={bar='',__name__=~'.%2b'}", code: 200, body: `# TYPE test_metric1 untyped
test_metric1{foo="bar",instance="i"} 10000 6000000
test_metric1{foo="boo",instance="i"} 1 6000000
# TYPE test_metric2 untyped
test_metric2{foo="boo",instance="i"} 1 6000000
# TYPE test_metric_old untyped
test_metric_old{instance=""} 981 5880000
# TYPE test_metric_without_labels untyped
test_metric_without_labels{instance=""} 1001 6000000
`}, "external labels are added if not already present": {params: "match[]={__name__=~'.%2b'}", externalLabels: model.LabelSet{"zone": "ie", "foo": "baz"}, code: 200, body: `# TYPE test_metric1 untyped
test_metric1{foo="bar",instance="i",zone="ie"} 10000 6000000
test_metric1{foo="boo",instance="i",zone="ie"} 1 6000000
# TYPE test_metric2 untyped
test_metric2{foo="boo",instance="i",zone="ie"} 1 6000000
# TYPE test_metric_old untyped
test_metric_old{foo="baz",instance="",zone="ie"} 981 5880000
# TYPE test_metric_without_labels untyped
test_metric_without_labels{foo="baz",instance="",zone="ie"} 1001 6000000
`}, "instance is an external label": {params: "match[]={__name__=~'.%2b'}", externalLabels: model.LabelSet{"instance": "baz"}, code: 200, body: `# TYPE test_metric1 untyped
test_metric1{foo="bar",instance="i"} 10000 6000000
test_metric1{foo="boo",instance="i"} 1 6000000
# TYPE test_metric2 untyped
test_metric2{foo="boo",instance="i"} 1 6000000
# TYPE test_metric_old untyped
test_metric_old{instance="baz"} 981 5880000
# TYPE test_metric_without_labels untyped
test_metric_without_labels{instance="baz"} 1001 6000000
`}}

func TestFederation(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	suite, err := promql.NewTest(t, `
		load 1m
			test_metric1{foo="bar",instance="i"}    0+100x100
			test_metric1{foo="boo",instance="i"}    1+0x100
			test_metric2{foo="boo",instance="i"}    1+0x100
			test_metric_without_labels 1+10x100
			test_metric_stale                       1+10x99 stale
			test_metric_old                         1+10x98
	`)
	if err != nil {
		t.Fatal(err)
	}
	defer suite.Close()
	if err := suite.Run(); err != nil {
		t.Fatal(err)
	}
	h := &Handler{storage: suite.Storage(), queryEngine: suite.QueryEngine(), now: func() model.Time {
		return 101 * 60 * 1000
	}, config: &config.Config{GlobalConfig: config.GlobalConfig{}}}
	for name, scenario := range scenarios {
		h.config.GlobalConfig.ExternalLabels = scenario.externalLabels
		req := httptest.NewRequest("GET", "http://example.org/federate?"+scenario.params, nil)
		res := httptest.NewRecorder()
		h.federation(res, req)
		if got, want := res.Code, scenario.code; got != want {
			t.Errorf("Scenario %q: got code %d, want %d", name, got, want)
		}
		if got, want := normalizeBody(res.Body), scenario.body; got != want {
			t.Errorf("Scenario %q: got body\n%s\n, want\n%s\n", name, got, want)
		}
	}
}
func normalizeBody(body *bytes.Buffer) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		lines		[]string
		lastHash	int
	)
	for line, err := body.ReadString('\n'); err == nil; line, err = body.ReadString('\n') {
		if line[0] == '#' && len(lines) > 0 {
			sort.Strings(lines[lastHash+1:])
			lastHash = len(lines)
		}
		lines = append(lines, line)
	}
	if len(lines) > 0 {
		sort.Strings(lines[lastHash+1:])
	}
	return strings.Join(lines, "")
}
