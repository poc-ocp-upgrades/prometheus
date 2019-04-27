package scrape

import (
	"fmt"
	"strconv"
	"testing"
	"time"
	"github.com/prometheus/prometheus/pkg/relabel"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/util/testutil"
	yaml "gopkg.in/yaml.v2"
)

func TestPopulateLabels(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cases := []struct {
		in	labels.Labels
		cfg	*config.ScrapeConfig
		res	labels.Labels
		resOrig	labels.Labels
		err	error
	}{{in: labels.FromMap(map[string]string{model.AddressLabel: "1.2.3.4:1000", "custom": "value"}), cfg: &config.ScrapeConfig{Scheme: "https", MetricsPath: "/metrics", JobName: "job"}, res: labels.FromMap(map[string]string{model.AddressLabel: "1.2.3.4:1000", model.InstanceLabel: "1.2.3.4:1000", model.SchemeLabel: "https", model.MetricsPathLabel: "/metrics", model.JobLabel: "job", "custom": "value"}), resOrig: labels.FromMap(map[string]string{model.AddressLabel: "1.2.3.4:1000", model.SchemeLabel: "https", model.MetricsPathLabel: "/metrics", model.JobLabel: "job", "custom": "value"})}, {in: labels.FromMap(map[string]string{model.AddressLabel: "1.2.3.4", model.SchemeLabel: "http", model.MetricsPathLabel: "/custom", model.JobLabel: "custom-job"}), cfg: &config.ScrapeConfig{Scheme: "https", MetricsPath: "/metrics", JobName: "job"}, res: labels.FromMap(map[string]string{model.AddressLabel: "1.2.3.4:80", model.InstanceLabel: "1.2.3.4:80", model.SchemeLabel: "http", model.MetricsPathLabel: "/custom", model.JobLabel: "custom-job"}), resOrig: labels.FromMap(map[string]string{model.AddressLabel: "1.2.3.4", model.SchemeLabel: "http", model.MetricsPathLabel: "/custom", model.JobLabel: "custom-job"})}, {in: labels.FromMap(map[string]string{model.AddressLabel: "[::1]", model.InstanceLabel: "custom-instance"}), cfg: &config.ScrapeConfig{Scheme: "https", MetricsPath: "/metrics", JobName: "job"}, res: labels.FromMap(map[string]string{model.AddressLabel: "[::1]:443", model.InstanceLabel: "custom-instance", model.SchemeLabel: "https", model.MetricsPathLabel: "/metrics", model.JobLabel: "job"}), resOrig: labels.FromMap(map[string]string{model.AddressLabel: "[::1]", model.InstanceLabel: "custom-instance", model.SchemeLabel: "https", model.MetricsPathLabel: "/metrics", model.JobLabel: "job"})}, {in: labels.FromStrings("custom", "value"), cfg: &config.ScrapeConfig{Scheme: "https", MetricsPath: "/metrics", JobName: "job"}, res: nil, resOrig: nil, err: fmt.Errorf("no address")}, {in: labels.FromStrings("custom", "host:1234"), cfg: &config.ScrapeConfig{Scheme: "https", MetricsPath: "/metrics", JobName: "job", RelabelConfigs: []*relabel.Config{{Action: relabel.Replace, Regex: relabel.MustNewRegexp("(.*)"), SourceLabels: model.LabelNames{"custom"}, Replacement: "${1}", TargetLabel: string(model.AddressLabel)}}}, res: labels.FromMap(map[string]string{model.AddressLabel: "host:1234", model.InstanceLabel: "host:1234", model.SchemeLabel: "https", model.MetricsPathLabel: "/metrics", model.JobLabel: "job", "custom": "host:1234"}), resOrig: labels.FromMap(map[string]string{model.SchemeLabel: "https", model.MetricsPathLabel: "/metrics", model.JobLabel: "job", "custom": "host:1234"})}, {in: labels.FromStrings("custom", "host:1234"), cfg: &config.ScrapeConfig{Scheme: "https", MetricsPath: "/metrics", JobName: "job", RelabelConfigs: []*relabel.Config{{Action: relabel.Replace, Regex: relabel.MustNewRegexp("(.*)"), SourceLabels: model.LabelNames{"custom"}, Replacement: "${1}", TargetLabel: string(model.AddressLabel)}}}, res: labels.FromMap(map[string]string{model.AddressLabel: "host:1234", model.InstanceLabel: "host:1234", model.SchemeLabel: "https", model.MetricsPathLabel: "/metrics", model.JobLabel: "job", "custom": "host:1234"}), resOrig: labels.FromMap(map[string]string{model.SchemeLabel: "https", model.MetricsPathLabel: "/metrics", model.JobLabel: "job", "custom": "host:1234"})}, {in: labels.FromMap(map[string]string{model.AddressLabel: "1.2.3.4:1000", "custom": "\xbd"}), cfg: &config.ScrapeConfig{Scheme: "https", MetricsPath: "/metrics", JobName: "job"}, res: nil, resOrig: nil, err: fmt.Errorf("invalid label value for \"custom\": \"\\xbd\"")}}
	for _, c := range cases {
		in := c.in.Copy()
		res, orig, err := populateLabels(c.in, c.cfg)
		testutil.Equals(t, c.err, err)
		testutil.Equals(t, c.in, in)
		testutil.Equals(t, c.res, res)
		testutil.Equals(t, c.resOrig, orig)
	}
}
func TestManagerReloadNoChange(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tsetName := "test"
	cfgText := `
scrape_configs:
 - job_name: '` + tsetName + `'
   static_configs:
   - targets: ["foo:9090"]
   - targets: ["bar:9090"]
`
	cfg := &config.Config{}
	if err := yaml.UnmarshalStrict([]byte(cfgText), cfg); err != nil {
		t.Fatalf("Unable to load YAML config cfgYaml: %s", err)
	}
	scrapeManager := NewManager(nil, nil)
	scrapeManager.ApplyConfig(cfg)
	newLoop := func(_ *Target, s scraper, _ int, _ bool, _ []*relabel.Config) loop {
		t.Fatal("reload happened")
		return nil
	}
	sp := &scrapePool{appendable: &nopAppendable{}, activeTargets: map[uint64]*Target{}, loops: map[uint64]loop{1: &testLoop{}}, newLoop: newLoop, logger: nil, config: cfg.ScrapeConfigs[0]}
	scrapeManager.scrapePools = map[string]*scrapePool{tsetName: sp}
	scrapeManager.ApplyConfig(cfg)
}
func TestManagerTargetsUpdates(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := NewManager(nil, nil)
	ts := make(chan map[string][]*targetgroup.Group)
	go m.Run(ts)
	tgSent := make(map[string][]*targetgroup.Group)
	for x := 0; x < 10; x++ {
		tgSent[strconv.Itoa(x)] = []*targetgroup.Group{{Source: strconv.Itoa(x)}}
		select {
		case ts <- tgSent:
		case <-time.After(10 * time.Millisecond):
			t.Error("Scrape manager's channel remained blocked after the set threshold.")
		}
	}
	m.mtxScrape.Lock()
	tsetActual := m.targetSets
	m.mtxScrape.Unlock()
	testutil.Equals(t, tgSent, tsetActual)
	select {
	case <-m.triggerReload:
	default:
		t.Error("No scrape loops reload was triggered after targets update.")
	}
}
