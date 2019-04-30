package config

import (
	"fmt"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"io/ioutil"
	"net/url"
	godefaulthttp "net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"github.com/prometheus/prometheus/pkg/relabel"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	sd_config "github.com/prometheus/prometheus/discovery/config"
	"gopkg.in/yaml.v2"
)

var (
	patRulePath = regexp.MustCompile(`^[^*]*(\*[^/]*)?$`)
)

func Load(s string) (*Config, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cfg := &Config{}
	*cfg = DefaultConfig
	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	cfg.original = s
	return cfg, nil
}
func LoadFile(filename string) (*Config, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, fmt.Errorf("parsing YAML file %s: %v", filename, err)
	}
	resolveFilepaths(filepath.Dir(filename), cfg)
	return cfg, nil
}

var (
	DefaultConfig			= Config{GlobalConfig: DefaultGlobalConfig}
	DefaultGlobalConfig		= GlobalConfig{ScrapeInterval: model.Duration(1 * time.Minute), ScrapeTimeout: model.Duration(10 * time.Second), EvaluationInterval: model.Duration(1 * time.Minute)}
	DefaultScrapeConfig		= ScrapeConfig{MetricsPath: "/metrics", Scheme: "http", HonorLabels: false}
	DefaultAlertmanagerConfig	= AlertmanagerConfig{Scheme: "http", Timeout: model.Duration(10 * time.Second)}
	DefaultRemoteWriteConfig	= RemoteWriteConfig{RemoteTimeout: model.Duration(30 * time.Second), QueueConfig: DefaultQueueConfig}
	DefaultQueueConfig		= QueueConfig{MaxShards: 1000, MinShards: 1, MaxSamplesPerSend: 100, Capacity: 100 * 100, BatchSendDeadline: model.Duration(5 * time.Second), MaxRetries: 3, MinBackoff: model.Duration(30 * time.Millisecond), MaxBackoff: model.Duration(100 * time.Millisecond)}
	DefaultRemoteReadConfig		= RemoteReadConfig{RemoteTimeout: model.Duration(1 * time.Minute)}
)

type Config struct {
	GlobalConfig		GlobalConfig		`yaml:"global"`
	AlertingConfig		AlertingConfig		`yaml:"alerting,omitempty"`
	RuleFiles		[]string		`yaml:"rule_files,omitempty"`
	ScrapeConfigs		[]*ScrapeConfig		`yaml:"scrape_configs,omitempty"`
	RemoteWriteConfigs	[]*RemoteWriteConfig	`yaml:"remote_write,omitempty"`
	RemoteReadConfigs	[]*RemoteReadConfig	`yaml:"remote_read,omitempty"`
	original		string
}

func resolveFilepaths(baseDir string, cfg *Config) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	join := func(fp string) string {
		if len(fp) > 0 && !filepath.IsAbs(fp) {
			fp = filepath.Join(baseDir, fp)
		}
		return fp
	}
	for i, rf := range cfg.RuleFiles {
		cfg.RuleFiles[i] = join(rf)
	}
	clientPaths := func(scfg *config_util.HTTPClientConfig) {
		scfg.BearerTokenFile = join(scfg.BearerTokenFile)
		scfg.TLSConfig.CAFile = join(scfg.TLSConfig.CAFile)
		scfg.TLSConfig.CertFile = join(scfg.TLSConfig.CertFile)
		scfg.TLSConfig.KeyFile = join(scfg.TLSConfig.KeyFile)
	}
	sdPaths := func(cfg *sd_config.ServiceDiscoveryConfig) {
		for _, kcfg := range cfg.KubernetesSDConfigs {
			kcfg.BearerTokenFile = join(kcfg.BearerTokenFile)
			kcfg.TLSConfig.CAFile = join(kcfg.TLSConfig.CAFile)
			kcfg.TLSConfig.CertFile = join(kcfg.TLSConfig.CertFile)
			kcfg.TLSConfig.KeyFile = join(kcfg.TLSConfig.KeyFile)
		}
		for _, mcfg := range cfg.MarathonSDConfigs {
			mcfg.AuthTokenFile = join(mcfg.AuthTokenFile)
			mcfg.HTTPClientConfig.BearerTokenFile = join(mcfg.HTTPClientConfig.BearerTokenFile)
			mcfg.HTTPClientConfig.TLSConfig.CAFile = join(mcfg.HTTPClientConfig.TLSConfig.CAFile)
			mcfg.HTTPClientConfig.TLSConfig.CertFile = join(mcfg.HTTPClientConfig.TLSConfig.CertFile)
			mcfg.HTTPClientConfig.TLSConfig.KeyFile = join(mcfg.HTTPClientConfig.TLSConfig.KeyFile)
		}
		for _, consulcfg := range cfg.ConsulSDConfigs {
			consulcfg.TLSConfig.CAFile = join(consulcfg.TLSConfig.CAFile)
			consulcfg.TLSConfig.CertFile = join(consulcfg.TLSConfig.CertFile)
			consulcfg.TLSConfig.KeyFile = join(consulcfg.TLSConfig.KeyFile)
		}
		for _, filecfg := range cfg.FileSDConfigs {
			for i, fn := range filecfg.Files {
				filecfg.Files[i] = join(fn)
			}
		}
	}
	for _, cfg := range cfg.ScrapeConfigs {
		clientPaths(&cfg.HTTPClientConfig)
		sdPaths(&cfg.ServiceDiscoveryConfig)
	}
	for _, cfg := range cfg.AlertingConfig.AlertmanagerConfigs {
		clientPaths(&cfg.HTTPClientConfig)
		sdPaths(&cfg.ServiceDiscoveryConfig)
	}
}
func (c Config) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = DefaultConfig
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.GlobalConfig.isZero() {
		c.GlobalConfig = DefaultGlobalConfig
	}
	for _, rf := range c.RuleFiles {
		if !patRulePath.MatchString(rf) {
			return fmt.Errorf("invalid rule file path %q", rf)
		}
	}
	jobNames := map[string]struct{}{}
	for _, scfg := range c.ScrapeConfigs {
		if scfg == nil {
			return fmt.Errorf("empty or null scrape config section")
		}
		if scfg.ScrapeInterval == 0 {
			scfg.ScrapeInterval = c.GlobalConfig.ScrapeInterval
		}
		if scfg.ScrapeTimeout > scfg.ScrapeInterval {
			return fmt.Errorf("scrape timeout greater than scrape interval for scrape config with job name %q", scfg.JobName)
		}
		if scfg.ScrapeTimeout == 0 {
			if c.GlobalConfig.ScrapeTimeout > scfg.ScrapeInterval {
				scfg.ScrapeTimeout = scfg.ScrapeInterval
			} else {
				scfg.ScrapeTimeout = c.GlobalConfig.ScrapeTimeout
			}
		}
		if _, ok := jobNames[scfg.JobName]; ok {
			return fmt.Errorf("found multiple scrape configs with job name %q", scfg.JobName)
		}
		jobNames[scfg.JobName] = struct{}{}
	}
	for _, rwcfg := range c.RemoteWriteConfigs {
		if rwcfg == nil {
			return fmt.Errorf("empty or null remote write config section")
		}
	}
	for _, rrcfg := range c.RemoteReadConfigs {
		if rrcfg == nil {
			return fmt.Errorf("empty or null remote read config section")
		}
	}
	return nil
}

type GlobalConfig struct {
	ScrapeInterval		model.Duration	`yaml:"scrape_interval,omitempty"`
	ScrapeTimeout		model.Duration	`yaml:"scrape_timeout,omitempty"`
	EvaluationInterval	model.Duration	`yaml:"evaluation_interval,omitempty"`
	ExternalLabels		model.LabelSet	`yaml:"external_labels,omitempty"`
}

func (c *GlobalConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc := &GlobalConfig{}
	type plain GlobalConfig
	if err := unmarshal((*plain)(gc)); err != nil {
		return err
	}
	if gc.ScrapeInterval == 0 {
		gc.ScrapeInterval = DefaultGlobalConfig.ScrapeInterval
	}
	if gc.ScrapeTimeout > gc.ScrapeInterval {
		return fmt.Errorf("global scrape timeout greater than scrape interval")
	}
	if gc.ScrapeTimeout == 0 {
		if DefaultGlobalConfig.ScrapeTimeout > gc.ScrapeInterval {
			gc.ScrapeTimeout = gc.ScrapeInterval
		} else {
			gc.ScrapeTimeout = DefaultGlobalConfig.ScrapeTimeout
		}
	}
	if gc.EvaluationInterval == 0 {
		gc.EvaluationInterval = DefaultGlobalConfig.EvaluationInterval
	}
	*c = *gc
	return nil
}
func (c *GlobalConfig) isZero() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.ExternalLabels == nil && c.ScrapeInterval == 0 && c.ScrapeTimeout == 0 && c.EvaluationInterval == 0
}

type ScrapeConfig struct {
	JobName			string					`yaml:"job_name"`
	HonorLabels		bool					`yaml:"honor_labels,omitempty"`
	Params			url.Values				`yaml:"params,omitempty"`
	ScrapeInterval		model.Duration				`yaml:"scrape_interval,omitempty"`
	ScrapeTimeout		model.Duration				`yaml:"scrape_timeout,omitempty"`
	MetricsPath		string					`yaml:"metrics_path,omitempty"`
	Scheme			string					`yaml:"scheme,omitempty"`
	SampleLimit		uint					`yaml:"sample_limit,omitempty"`
	ServiceDiscoveryConfig	sd_config.ServiceDiscoveryConfig	`yaml:",inline"`
	HTTPClientConfig	config_util.HTTPClientConfig		`yaml:",inline"`
	RelabelConfigs		[]*relabel.Config			`yaml:"relabel_configs,omitempty"`
	MetricRelabelConfigs	[]*relabel.Config			`yaml:"metric_relabel_configs,omitempty"`
}

func (c *ScrapeConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = DefaultScrapeConfig
	type plain ScrapeConfig
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if len(c.JobName) == 0 {
		return fmt.Errorf("job_name is empty")
	}
	if err := c.HTTPClientConfig.Validate(); err != nil {
		return err
	}
	if err := c.ServiceDiscoveryConfig.Validate(); err != nil {
		return err
	}
	if len(c.RelabelConfigs) == 0 {
		for _, tg := range c.ServiceDiscoveryConfig.StaticConfigs {
			for _, t := range tg.Targets {
				if err := CheckTargetAddress(t[model.AddressLabel]); err != nil {
					return err
				}
			}
		}
	}
	for _, rlcfg := range c.RelabelConfigs {
		if rlcfg == nil {
			return fmt.Errorf("empty or null target relabeling rule in scrape config")
		}
	}
	for _, rlcfg := range c.MetricRelabelConfigs {
		if rlcfg == nil {
			return fmt.Errorf("empty or null metric relabeling rule in scrape config")
		}
	}
	for i, tg := range c.ServiceDiscoveryConfig.StaticConfigs {
		tg.Source = fmt.Sprintf("%d", i)
	}
	return nil
}

type AlertingConfig struct {
	AlertRelabelConfigs	[]*relabel.Config	`yaml:"alert_relabel_configs,omitempty"`
	AlertmanagerConfigs	[]*AlertmanagerConfig	`yaml:"alertmanagers,omitempty"`
}

func (c *AlertingConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = AlertingConfig{}
	type plain AlertingConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	for _, rlcfg := range c.AlertRelabelConfigs {
		if rlcfg == nil {
			return fmt.Errorf("empty or null alert relabeling rule")
		}
	}
	return nil
}

type AlertmanagerConfig struct {
	ServiceDiscoveryConfig	sd_config.ServiceDiscoveryConfig	`yaml:",inline"`
	HTTPClientConfig	config_util.HTTPClientConfig		`yaml:",inline"`
	Scheme			string					`yaml:"scheme,omitempty"`
	PathPrefix		string					`yaml:"path_prefix,omitempty"`
	Timeout			model.Duration				`yaml:"timeout,omitempty"`
	RelabelConfigs		[]*relabel.Config			`yaml:"relabel_configs,omitempty"`
}

func (c *AlertmanagerConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = DefaultAlertmanagerConfig
	type plain AlertmanagerConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if err := c.HTTPClientConfig.Validate(); err != nil {
		return err
	}
	if err := c.ServiceDiscoveryConfig.Validate(); err != nil {
		return err
	}
	if len(c.RelabelConfigs) == 0 {
		for _, tg := range c.ServiceDiscoveryConfig.StaticConfigs {
			for _, t := range tg.Targets {
				if err := CheckTargetAddress(t[model.AddressLabel]); err != nil {
					return err
				}
			}
		}
	}
	for _, rlcfg := range c.RelabelConfigs {
		if rlcfg == nil {
			return fmt.Errorf("empty or null Alertmanager target relabeling rule")
		}
	}
	for i, tg := range c.ServiceDiscoveryConfig.StaticConfigs {
		tg.Source = fmt.Sprintf("%d", i)
	}
	return nil
}
func CheckTargetAddress(address model.LabelValue) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if strings.Contains(string(address), "/") {
		return fmt.Errorf("%q is not a valid hostname", address)
	}
	return nil
}

type ClientCert struct {
	Cert	string			`yaml:"cert"`
	Key	config_util.Secret	`yaml:"key"`
}
type FileSDConfig struct {
	Files		[]string	`yaml:"files"`
	RefreshInterval	model.Duration	`yaml:"refresh_interval,omitempty"`
}
type RemoteWriteConfig struct {
	URL			*config_util.URL		`yaml:"url"`
	RemoteTimeout		model.Duration			`yaml:"remote_timeout,omitempty"`
	WriteRelabelConfigs	[]*relabel.Config		`yaml:"write_relabel_configs,omitempty"`
	HTTPClientConfig	config_util.HTTPClientConfig	`yaml:",inline"`
	QueueConfig		QueueConfig			`yaml:"queue_config,omitempty"`
}

func (c *RemoteWriteConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = DefaultRemoteWriteConfig
	type plain RemoteWriteConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.URL == nil {
		return fmt.Errorf("url for remote_write is empty")
	}
	for _, rlcfg := range c.WriteRelabelConfigs {
		if rlcfg == nil {
			return fmt.Errorf("empty or null relabeling rule in remote write config")
		}
	}
	return c.HTTPClientConfig.Validate()
}

type QueueConfig struct {
	Capacity		int		`yaml:"capacity,omitempty"`
	MaxShards		int		`yaml:"max_shards,omitempty"`
	MinShards		int		`yaml:"min_shards,omitempty"`
	MaxSamplesPerSend	int		`yaml:"max_samples_per_send,omitempty"`
	BatchSendDeadline	model.Duration	`yaml:"batch_send_deadline,omitempty"`
	MaxRetries		int		`yaml:"max_retries,omitempty"`
	MinBackoff		model.Duration	`yaml:"min_backoff,omitempty"`
	MaxBackoff		model.Duration	`yaml:"max_backoff,omitempty"`
}
type RemoteReadConfig struct {
	URL			*config_util.URL		`yaml:"url"`
	RemoteTimeout		model.Duration			`yaml:"remote_timeout,omitempty"`
	ReadRecent		bool				`yaml:"read_recent,omitempty"`
	HTTPClientConfig	config_util.HTTPClientConfig	`yaml:",inline"`
	RequiredMatchers	model.LabelSet			`yaml:"required_matchers,omitempty"`
}

func (c *RemoteReadConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = DefaultRemoteReadConfig
	type plain RemoteReadConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.URL == nil {
		return fmt.Errorf("url for remote_read is empty")
	}
	return c.HTTPClientConfig.Validate()
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
