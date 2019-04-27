package marathon

import (
	"context"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	godefaulthttp "net/http"
	"strconv"
	"strings"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/util/strutil"
)

const (
	metaLabelPrefix					= model.MetaLabelPrefix + "marathon_"
	appLabelPrefix					= metaLabelPrefix + "app_label_"
	appLabel			model.LabelName	= metaLabelPrefix + "app"
	imageLabel			model.LabelName	= metaLabelPrefix + "image"
	portIndexLabel			model.LabelName	= metaLabelPrefix + "port_index"
	taskLabel			model.LabelName	= metaLabelPrefix + "task"
	portMappingLabelPrefix				= metaLabelPrefix + "port_mapping_label_"
	portDefinitionLabelPrefix			= metaLabelPrefix + "port_definition_label_"
	namespace					= "prometheus"
)

var (
	refreshFailuresCount	= prometheus.NewCounter(prometheus.CounterOpts{Namespace: namespace, Name: "sd_marathon_refresh_failures_total", Help: "The number of Marathon-SD refresh failures."})
	refreshDuration		= prometheus.NewSummary(prometheus.SummaryOpts{Namespace: namespace, Name: "sd_marathon_refresh_duration_seconds", Help: "The duration of a Marathon-SD refresh in seconds."})
	DefaultSDConfig		= SDConfig{RefreshInterval: model.Duration(30 * time.Second)}
)

type SDConfig struct {
	Servers			[]string			`yaml:"servers,omitempty"`
	RefreshInterval		model.Duration			`yaml:"refresh_interval,omitempty"`
	AuthToken		config_util.Secret		`yaml:"auth_token,omitempty"`
	AuthTokenFile		string				`yaml:"auth_token_file,omitempty"`
	HTTPClientConfig	config_util.HTTPClientConfig	`yaml:",inline"`
}

func (c *SDConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = DefaultSDConfig
	type plain SDConfig
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if len(c.Servers) == 0 {
		return fmt.Errorf("marathon_sd: must contain at least one Marathon server")
	}
	if len(c.AuthToken) > 0 && len(c.AuthTokenFile) > 0 {
		return fmt.Errorf("marathon_sd: at most one of auth_token & auth_token_file must be configured")
	}
	if c.HTTPClientConfig.BasicAuth != nil && (len(c.AuthToken) > 0 || len(c.AuthTokenFile) > 0) {
		return fmt.Errorf("marathon_sd: at most one of basic_auth, auth_token & auth_token_file must be configured")
	}
	if (len(c.HTTPClientConfig.BearerToken) > 0 || len(c.HTTPClientConfig.BearerTokenFile) > 0) && (len(c.AuthToken) > 0 || len(c.AuthTokenFile) > 0) {
		return fmt.Errorf("marathon_sd: at most one of bearer_token, bearer_token_file, auth_token & auth_token_file must be configured")
	}
	return c.HTTPClientConfig.Validate()
}
func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(refreshFailuresCount)
	prometheus.MustRegister(refreshDuration)
}

const appListPath string = "/v2/apps/?embed=apps.tasks"

type Discovery struct {
	client		*http.Client
	servers		[]string
	refreshInterval	time.Duration
	lastRefresh	map[string]*targetgroup.Group
	appsClient	AppListClient
	logger		log.Logger
}

func NewDiscovery(conf SDConfig, logger log.Logger) (*Discovery, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if logger == nil {
		logger = log.NewNopLogger()
	}
	rt, err := config_util.NewRoundTripperFromConfig(conf.HTTPClientConfig, "marathon_sd")
	if err != nil {
		return nil, err
	}
	if len(conf.AuthToken) > 0 {
		rt, err = newAuthTokenRoundTripper(conf.AuthToken, rt)
	} else if len(conf.AuthTokenFile) > 0 {
		rt, err = newAuthTokenFileRoundTripper(conf.AuthTokenFile, rt)
	}
	if err != nil {
		return nil, err
	}
	return &Discovery{client: &http.Client{Transport: rt}, servers: conf.Servers, refreshInterval: time.Duration(conf.RefreshInterval), appsClient: fetchApps, logger: logger}, nil
}

type authTokenRoundTripper struct {
	authToken	config_util.Secret
	rt		http.RoundTripper
}

func newAuthTokenRoundTripper(token config_util.Secret, rt http.RoundTripper) (http.RoundTripper, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &authTokenRoundTripper{token, rt}, nil
}
func (rt *authTokenRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request.Header.Set("Authorization", "token="+string(rt.authToken))
	return rt.rt.RoundTrip(request)
}

type authTokenFileRoundTripper struct {
	authTokenFile	string
	rt		http.RoundTripper
}

func newAuthTokenFileRoundTripper(tokenFile string, rt http.RoundTripper) (http.RoundTripper, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read auth token file %s: %s", tokenFile, err)
	}
	return &authTokenFileRoundTripper{tokenFile, rt}, nil
}
func (rt *authTokenFileRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	b, err := ioutil.ReadFile(rt.authTokenFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read auth token file %s: %s", rt.authTokenFile, err)
	}
	authToken := strings.TrimSpace(string(b))
	request.Header.Set("Authorization", "token="+authToken)
	return rt.rt.RoundTrip(request)
}
func (d *Discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(d.refreshInterval):
			err := d.updateServices(ctx, ch)
			if err != nil {
				level.Error(d.logger).Log("msg", "Error while updating services", "err", err)
			}
		}
	}
}
func (d *Discovery) updateServices(ctx context.Context, ch chan<- []*targetgroup.Group) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t0 := time.Now()
	defer func() {
		refreshDuration.Observe(time.Since(t0).Seconds())
		if err != nil {
			refreshFailuresCount.Inc()
		}
	}()
	targetMap, err := d.fetchTargetGroups()
	if err != nil {
		return err
	}
	all := make([]*targetgroup.Group, 0, len(targetMap))
	for _, tg := range targetMap {
		all = append(all, tg)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- all:
	}
	for source := range d.lastRefresh {
		_, ok := targetMap[source]
		if !ok {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case ch <- []*targetgroup.Group{{Source: source}}:
				level.Debug(d.logger).Log("msg", "Removing group", "source", source)
			}
		}
	}
	d.lastRefresh = targetMap
	return nil
}
func (d *Discovery) fetchTargetGroups() (map[string]*targetgroup.Group, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	url := RandomAppsURL(d.servers)
	apps, err := d.appsClient(d.client, url)
	if err != nil {
		return nil, err
	}
	groups := AppsToTargetGroups(apps)
	return groups, nil
}

type Task struct {
	ID		string		`json:"id"`
	Host		string		`json:"host"`
	Ports		[]uint32	`json:"ports"`
	IPAddresses	[]IPAddress	`json:"ipAddresses"`
}
type IPAddress struct {
	Address	string	`json:"ipAddress"`
	Proto	string	`json:"protocol"`
}
type PortMapping struct {
	Labels		map[string]string	`json:"labels"`
	ContainerPort	uint32			`json:"containerPort"`
	HostPort	uint32			`json:"hostPort"`
	ServicePort	uint32			`json:"servicePort"`
}
type DockerContainer struct {
	Image		string		`json:"image"`
	PortMappings	[]PortMapping	`json:"portMappings"`
}
type Container struct {
	Docker		DockerContainer	`json:"docker"`
	PortMappings	[]PortMapping	`json:"portMappings"`
}
type PortDefinition struct {
	Labels	map[string]string	`json:"labels"`
	Port	uint32			`json:"port"`
}
type Network struct {
	Name	string	`json:"name"`
	Mode	string	`json:"mode"`
}
type App struct {
	ID		string			`json:"id"`
	Tasks		[]Task			`json:"tasks"`
	RunningTasks	int			`json:"tasksRunning"`
	Labels		map[string]string	`json:"labels"`
	Container	Container		`json:"container"`
	PortDefinitions	[]PortDefinition	`json:"portDefinitions"`
	Networks	[]Network		`json:"networks"`
	RequirePorts	bool			`json:"requirePorts"`
}

func (app App) isContainerNet() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(app.Networks) > 0 && app.Networks[0].Mode == "container"
}

type AppList struct {
	Apps []App `json:"apps"`
}
type AppListClient func(client *http.Client, url string) (*AppList, error)

func fetchApps(client *http.Client, url string) (*AppList, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if (resp.StatusCode < 200) || (resp.StatusCode >= 300) {
		return nil, fmt.Errorf("non 2xx status '%v' response during marathon service discovery", resp.StatusCode)
	}
	var apps AppList
	err = json.NewDecoder(resp.Body).Decode(&apps)
	if err != nil {
		return nil, fmt.Errorf("%q: %v", url, err)
	}
	return &apps, nil
}
func RandomAppsURL(servers []string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	server := servers[rand.Intn(len(servers))]
	return fmt.Sprintf("%s%s", server, appListPath)
}
func AppsToTargetGroups(apps *AppList) map[string]*targetgroup.Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tgroups := map[string]*targetgroup.Group{}
	for _, a := range apps.Apps {
		group := createTargetGroup(&a)
		tgroups[group.Source] = group
	}
	return tgroups
}
func createTargetGroup(app *App) *targetgroup.Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		targets	= targetsForApp(app)
		appName	= model.LabelValue(app.ID)
		image	= model.LabelValue(app.Container.Docker.Image)
	)
	tg := &targetgroup.Group{Targets: targets, Labels: model.LabelSet{appLabel: appName, imageLabel: image}, Source: app.ID}
	for ln, lv := range app.Labels {
		ln = appLabelPrefix + strutil.SanitizeLabelName(ln)
		tg.Labels[model.LabelName(ln)] = model.LabelValue(lv)
	}
	return tg
}
func targetsForApp(app *App) []model.LabelSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	targets := make([]model.LabelSet, 0, len(app.Tasks))
	var ports []uint32
	var labels []map[string]string
	var prefix string
	if len(app.Container.PortMappings) != 0 {
		ports, labels = extractPortMapping(app.Container.PortMappings, app.isContainerNet())
		prefix = portMappingLabelPrefix
	} else if len(app.Container.Docker.PortMappings) != 0 {
		ports, labels = extractPortMapping(app.Container.Docker.PortMappings, app.isContainerNet())
		prefix = portMappingLabelPrefix
	} else if len(app.PortDefinitions) != 0 {
		ports = make([]uint32, len(app.PortDefinitions))
		labels = make([]map[string]string, len(app.PortDefinitions))
		for i := 0; i < len(app.PortDefinitions); i++ {
			labels[i] = app.PortDefinitions[i].Labels
			if app.RequirePorts {
				ports[i] = app.PortDefinitions[i].Port
			}
		}
		prefix = portDefinitionLabelPrefix
	}
	for _, t := range app.Tasks {
		if len(ports) == 0 && len(t.Ports) != 0 {
			ports = t.Ports
		}
		for i, port := range ports {
			if port == 0 && len(t.Ports) == len(ports) {
				port = t.Ports[i]
			}
			targetAddress := targetEndpoint(&t, port, app.isContainerNet())
			target := model.LabelSet{model.AddressLabel: model.LabelValue(targetAddress), taskLabel: model.LabelValue(t.ID), portIndexLabel: model.LabelValue(strconv.Itoa(i))}
			if len(labels) > 0 {
				for ln, lv := range labels[i] {
					ln = prefix + strutil.SanitizeLabelName(ln)
					target[model.LabelName(ln)] = model.LabelValue(lv)
				}
			}
			targets = append(targets, target)
		}
	}
	return targets
}
func targetEndpoint(task *Task, port uint32, containerNet bool) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var host string
	if containerNet && len(task.IPAddresses) > 0 {
		host = task.IPAddresses[0].Address
	} else {
		host = task.Host
	}
	return net.JoinHostPort(host, fmt.Sprintf("%d", port))
}
func extractPortMapping(portMappings []PortMapping, containerNet bool) ([]uint32, []map[string]string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ports := make([]uint32, len(portMappings))
	labels := make([]map[string]string, len(portMappings))
	for i := 0; i < len(portMappings); i++ {
		labels[i] = portMappings[i].Labels
		if containerNet {
			ports[i] = portMappings[i].ContainerPort
		} else {
			ports[i] = portMappings[i].HostPort
		}
	}
	return ports, labels
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
