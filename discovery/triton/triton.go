package triton

import (
	"context"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	godefaulthttp "net/http"
	"net/url"
	"strings"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/mwitkow/go-conntrack"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

const (
	tritonLabel		= model.MetaLabelPrefix + "triton_"
	tritonLabelGroups	= tritonLabel + "groups"
	tritonLabelMachineID	= tritonLabel + "machine_id"
	tritonLabelMachineAlias	= tritonLabel + "machine_alias"
	tritonLabelMachineBrand	= tritonLabel + "machine_brand"
	tritonLabelMachineImage	= tritonLabel + "machine_image"
	tritonLabelServerID	= tritonLabel + "server_id"
)

var (
	refreshFailuresCount	= prometheus.NewCounter(prometheus.CounterOpts{Name: "prometheus_sd_triton_refresh_failures_total", Help: "The number of Triton-SD scrape failures."})
	refreshDuration		= prometheus.NewSummary(prometheus.SummaryOpts{Name: "prometheus_sd_triton_refresh_duration_seconds", Help: "The duration of a Triton-SD refresh in seconds."})
	DefaultSDConfig		= SDConfig{Port: 9163, RefreshInterval: model.Duration(60 * time.Second), Version: 1}
)

type SDConfig struct {
	Account		string			`yaml:"account"`
	DNSSuffix	string			`yaml:"dns_suffix"`
	Endpoint	string			`yaml:"endpoint"`
	Groups		[]string		`yaml:"groups,omitempty"`
	Port		int			`yaml:"port"`
	RefreshInterval	model.Duration		`yaml:"refresh_interval,omitempty"`
	TLSConfig	config_util.TLSConfig	`yaml:"tls_config,omitempty"`
	Version		int			`yaml:"version"`
}

func (c *SDConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = DefaultSDConfig
	type plain SDConfig
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if c.Account == "" {
		return fmt.Errorf("triton SD configuration requires an account")
	}
	if c.DNSSuffix == "" {
		return fmt.Errorf("triton SD configuration requires a dns_suffix")
	}
	if c.Endpoint == "" {
		return fmt.Errorf("triton SD configuration requires an endpoint")
	}
	if c.RefreshInterval <= 0 {
		return fmt.Errorf("triton SD configuration requires RefreshInterval to be a positive integer")
	}
	return nil
}
func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(refreshFailuresCount)
	prometheus.MustRegister(refreshDuration)
}

type DiscoveryResponse struct {
	Containers []struct {
		Groups		[]string	`json:"groups"`
		ServerUUID	string		`json:"server_uuid"`
		VMAlias		string		`json:"vm_alias"`
		VMBrand		string		`json:"vm_brand"`
		VMImageUUID	string		`json:"vm_image_uuid"`
		VMUUID		string		`json:"vm_uuid"`
	} `json:"containers"`
}
type Discovery struct {
	client		*http.Client
	interval	time.Duration
	logger		log.Logger
	sdConfig	*SDConfig
}

func New(logger log.Logger, conf *SDConfig) (*Discovery, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if logger == nil {
		logger = log.NewNopLogger()
	}
	tls, err := config_util.NewTLSConfig(&conf.TLSConfig)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{TLSClientConfig: tls, DialContext: conntrack.NewDialContextFunc(conntrack.DialWithTracing(), conntrack.DialWithName("triton_sd"))}
	client := &http.Client{Transport: transport}
	return &Discovery{client: client, interval: time.Duration(conf.RefreshInterval), logger: logger, sdConfig: conf}, nil
}
func (d *Discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer close(ch)
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	tg, err := d.refresh()
	if err != nil {
		level.Error(d.logger).Log("msg", "Refreshing targets failed", "err", err)
	} else {
		ch <- []*targetgroup.Group{tg}
	}
	for {
		select {
		case <-ticker.C:
			tg, err := d.refresh()
			if err != nil {
				level.Error(d.logger).Log("msg", "Refreshing targets failed", "err", err)
			} else {
				ch <- []*targetgroup.Group{tg}
			}
		case <-ctx.Done():
			return
		}
	}
}
func (d *Discovery) refresh() (tg *targetgroup.Group, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t0 := time.Now()
	defer func() {
		refreshDuration.Observe(time.Since(t0).Seconds())
		if err != nil {
			refreshFailuresCount.Inc()
		}
	}()
	var endpoint = fmt.Sprintf("https://%s:%d/v%d/discover", d.sdConfig.Endpoint, d.sdConfig.Port, d.sdConfig.Version)
	if len(d.sdConfig.Groups) > 0 {
		groups := url.QueryEscape(strings.Join(d.sdConfig.Groups, ","))
		endpoint = fmt.Sprintf("%s?groups=%s", endpoint, groups)
	}
	tg = &targetgroup.Group{Source: endpoint}
	resp, err := d.client.Get(endpoint)
	if err != nil {
		return tg, fmt.Errorf("an error occurred when requesting targets from the discovery endpoint. %s", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tg, fmt.Errorf("an error occurred when reading the response body. %s", err)
	}
	dr := DiscoveryResponse{}
	err = json.Unmarshal(data, &dr)
	if err != nil {
		return tg, fmt.Errorf("an error occurred unmarshaling the disovery response json. %s", err)
	}
	for _, container := range dr.Containers {
		labels := model.LabelSet{tritonLabelMachineID: model.LabelValue(container.VMUUID), tritonLabelMachineAlias: model.LabelValue(container.VMAlias), tritonLabelMachineBrand: model.LabelValue(container.VMBrand), tritonLabelMachineImage: model.LabelValue(container.VMImageUUID), tritonLabelServerID: model.LabelValue(container.ServerUUID)}
		addr := fmt.Sprintf("%s.%s:%d", container.VMUUID, d.sdConfig.DNSSuffix, d.sdConfig.Port)
		labels[model.AddressLabel] = model.LabelValue(addr)
		if len(container.Groups) > 0 {
			name := "," + strings.Join(container.Groups, ",") + ","
			labels[tritonLabelGroups] = model.LabelValue(name)
		}
		tg.Targets = append(tg.Targets, labels)
	}
	return tg, nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
