package kubernetes

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const (
	metaLabelPrefix		= model.MetaLabelPrefix + "kubernetes_"
	namespaceLabel		= metaLabelPrefix + "namespace"
	metricsNamespace	= "prometheus_sd_kubernetes"
)

var (
	eventCount		= prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: metricsNamespace, Name: "events_total", Help: "The number of Kubernetes events handled."}, []string{"role", "event"})
	DefaultSDConfig	= SDConfig{}
)

type Role string

const (
	RoleNode		Role	= "node"
	RolePod			Role	= "pod"
	RoleService		Role	= "service"
	RoleEndpoint	Role	= "endpoints"
	RoleIngress		Role	= "ingress"
)

func (c *Role) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := unmarshal((*string)(c)); err != nil {
		return err
	}
	switch *c {
	case RoleNode, RolePod, RoleService, RoleEndpoint, RoleIngress:
		return nil
	default:
		return fmt.Errorf("unknown Kubernetes SD role %q", *c)
	}
}

type SDConfig struct {
	APIServer			config_util.URL			`yaml:"api_server,omitempty"`
	Role				Role					`yaml:"role"`
	BasicAuth			*config_util.BasicAuth	`yaml:"basic_auth,omitempty"`
	BearerToken			config_util.Secret		`yaml:"bearer_token,omitempty"`
	BearerTokenFile		string					`yaml:"bearer_token_file,omitempty"`
	TLSConfig			config_util.TLSConfig	`yaml:"tls_config,omitempty"`
	NamespaceDiscovery	NamespaceDiscovery		`yaml:"namespaces,omitempty"`
}

func (c *SDConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = SDConfig{}
	type plain SDConfig
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if c.Role == "" {
		return fmt.Errorf("role missing (one of: pod, service, endpoints, node, ingress)")
	}
	if len(c.BearerToken) > 0 && len(c.BearerTokenFile) > 0 {
		return fmt.Errorf("at most one of bearer_token & bearer_token_file must be configured")
	}
	if c.BasicAuth != nil && (len(c.BearerToken) > 0 || len(c.BearerTokenFile) > 0) {
		return fmt.Errorf("at most one of basic_auth, bearer_token & bearer_token_file must be configured")
	}
	if c.APIServer.URL == nil && (c.BasicAuth != nil || c.BearerToken != "" || c.BearerTokenFile != "" || c.TLSConfig.CAFile != "" || c.TLSConfig.CertFile != "" || c.TLSConfig.KeyFile != "") {
		return fmt.Errorf("to use custom authentication please provide the 'api_server' URL explicitly")
	}
	return nil
}

type NamespaceDiscovery struct {
	Names []string `yaml:"names"`
}

func (c *NamespaceDiscovery) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = NamespaceDiscovery{}
	type plain NamespaceDiscovery
	return unmarshal((*plain)(c))
}
func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(eventCount)
	for _, role := range []string{"endpoints", "node", "pod", "service", "ingress"} {
		for _, evt := range []string{"add", "delete", "update"} {
			eventCount.WithLabelValues(role, evt)
		}
	}
	var (
		clientGoRequestMetricAdapterInstance		= clientGoRequestMetricAdapter{}
		clientGoCacheMetricsProviderInstance		= clientGoCacheMetricsProvider{}
		clientGoWorkqueueMetricsProviderInstance	= clientGoWorkqueueMetricsProvider{}
	)
	clientGoRequestMetricAdapterInstance.Register(prometheus.DefaultRegisterer)
	clientGoCacheMetricsProviderInstance.Register(prometheus.DefaultRegisterer)
	clientGoWorkqueueMetricsProviderInstance.Register(prometheus.DefaultRegisterer)
}

type discoverer interface {
	Run(ctx context.Context, up chan<- []*targetgroup.Group)
}
type Discovery struct {
	sync.RWMutex
	client				kubernetes.Interface
	role				Role
	logger				log.Logger
	namespaceDiscovery	*NamespaceDiscovery
	discoverers			[]discoverer
}

func (d *Discovery) getNamespaces() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	namespaces := d.namespaceDiscovery.Names
	if len(namespaces) == 0 {
		namespaces = []string{apiv1.NamespaceAll}
	}
	return namespaces
}
func New(l log.Logger, conf *SDConfig) (*Discovery, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if l == nil {
		l = log.NewNopLogger()
	}
	var (
		kcfg	*rest.Config
		err		error
	)
	if conf.APIServer.URL == nil {
		kcfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		level.Info(l).Log("msg", "Using pod service account via in-cluster config")
		if conf.TLSConfig.CAFile != "" {
			level.Warn(l).Log("msg", "Configured TLS CA file is ignored when using pod service account")
		}
		if conf.TLSConfig.CertFile != "" || conf.TLSConfig.KeyFile != "" {
			level.Warn(l).Log("msg", "Configured TLS client certificate is ignored when using pod service account")
		}
		if conf.BearerToken != "" {
			level.Warn(l).Log("msg", "Configured auth token is ignored when using pod service account")
		}
		if conf.BasicAuth != nil {
			level.Warn(l).Log("msg", "Configured basic authentication credentials are ignored when using pod service account")
		}
	} else {
		kcfg = &rest.Config{Host: conf.APIServer.String(), TLSClientConfig: rest.TLSClientConfig{CAFile: conf.TLSConfig.CAFile, CertFile: conf.TLSConfig.CertFile, KeyFile: conf.TLSConfig.KeyFile, Insecure: conf.TLSConfig.InsecureSkipVerify}}
		token := string(conf.BearerToken)
		if conf.BearerTokenFile != "" {
			bf, err := ioutil.ReadFile(conf.BearerTokenFile)
			if err != nil {
				return nil, err
			}
			token = string(bf)
		}
		kcfg.BearerToken = token
		if conf.BasicAuth != nil {
			kcfg.Username = conf.BasicAuth.Username
			kcfg.Password = string(conf.BasicAuth.Password)
		}
	}
	kcfg.UserAgent = "Prometheus/discovery"
	c, err := kubernetes.NewForConfig(kcfg)
	if err != nil {
		return nil, err
	}
	return &Discovery{client: c, logger: l, role: conf.Role, namespaceDiscovery: &conf.NamespaceDiscovery, discoverers: make([]discoverer, 0)}, nil
}

const resyncPeriod = 10 * time.Minute

func (d *Discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d.Lock()
	namespaces := d.getNamespaces()
	switch d.role {
	case RoleEndpoint:
		for _, namespace := range namespaces {
			e := d.client.CoreV1().Endpoints(namespace)
			elw := &cache.ListWatch{ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return e.List(options)
			}, WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return e.Watch(options)
			}}
			s := d.client.CoreV1().Services(namespace)
			slw := &cache.ListWatch{ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return s.List(options)
			}, WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return s.Watch(options)
			}}
			p := d.client.CoreV1().Pods(namespace)
			plw := &cache.ListWatch{ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return p.List(options)
			}, WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return p.Watch(options)
			}}
			eps := NewEndpoints(log.With(d.logger, "role", "endpoint"), cache.NewSharedInformer(slw, &apiv1.Service{}, resyncPeriod), cache.NewSharedInformer(elw, &apiv1.Endpoints{}, resyncPeriod), cache.NewSharedInformer(plw, &apiv1.Pod{}, resyncPeriod))
			d.discoverers = append(d.discoverers, eps)
			go eps.endpointsInf.Run(ctx.Done())
			go eps.serviceInf.Run(ctx.Done())
			go eps.podInf.Run(ctx.Done())
		}
	case RolePod:
		for _, namespace := range namespaces {
			p := d.client.CoreV1().Pods(namespace)
			plw := &cache.ListWatch{ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return p.List(options)
			}, WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return p.Watch(options)
			}}
			pod := NewPod(log.With(d.logger, "role", "pod"), cache.NewSharedInformer(plw, &apiv1.Pod{}, resyncPeriod))
			d.discoverers = append(d.discoverers, pod)
			go pod.informer.Run(ctx.Done())
		}
	case RoleService:
		for _, namespace := range namespaces {
			s := d.client.CoreV1().Services(namespace)
			slw := &cache.ListWatch{ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return s.List(options)
			}, WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return s.Watch(options)
			}}
			svc := NewService(log.With(d.logger, "role", "service"), cache.NewSharedInformer(slw, &apiv1.Service{}, resyncPeriod))
			d.discoverers = append(d.discoverers, svc)
			go svc.informer.Run(ctx.Done())
		}
	case RoleIngress:
		for _, namespace := range namespaces {
			i := d.client.ExtensionsV1beta1().Ingresses(namespace)
			ilw := &cache.ListWatch{ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return i.List(options)
			}, WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return i.Watch(options)
			}}
			ingress := NewIngress(log.With(d.logger, "role", "ingress"), cache.NewSharedInformer(ilw, &extensionsv1beta1.Ingress{}, resyncPeriod))
			d.discoverers = append(d.discoverers, ingress)
			go ingress.informer.Run(ctx.Done())
		}
	case RoleNode:
		nlw := &cache.ListWatch{ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return d.client.CoreV1().Nodes().List(options)
		}, WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return d.client.CoreV1().Nodes().Watch(options)
		}}
		node := NewNode(log.With(d.logger, "role", "node"), cache.NewSharedInformer(nlw, &apiv1.Node{}, resyncPeriod))
		d.discoverers = append(d.discoverers, node)
		go node.informer.Run(ctx.Done())
	default:
		level.Error(d.logger).Log("msg", "unknown Kubernetes discovery kind", "role", d.role)
	}
	var wg sync.WaitGroup
	for _, dd := range d.discoverers {
		wg.Add(1)
		go func(d discoverer) {
			defer wg.Done()
			d.Run(ctx, ch)
		}(dd)
	}
	d.Unlock()
	wg.Wait()
	<-ctx.Done()
}
func lv(s string) model.LabelValue {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return model.LabelValue(s)
}
func send(ctx context.Context, l log.Logger, role Role, ch chan<- []*targetgroup.Group, tg *targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if tg == nil {
		return
	}
	select {
	case <-ctx.Done():
	case ch <- []*targetgroup.Group{tg}:
	}
}
