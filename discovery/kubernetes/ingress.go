package kubernetes

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/util/strutil"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Ingress struct {
	logger		log.Logger
	informer	cache.SharedInformer
	store		cache.Store
	queue		*workqueue.Type
}

func NewIngress(l log.Logger, inf cache.SharedInformer) *Ingress {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := &Ingress{logger: l, informer: inf, store: inf.GetStore(), queue: workqueue.NewNamed("ingress")}
	s.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{AddFunc: func(o interface{}) {
		eventCount.WithLabelValues("ingress", "add").Inc()
		s.enqueue(o)
	}, DeleteFunc: func(o interface{}) {
		eventCount.WithLabelValues("ingress", "delete").Inc()
		s.enqueue(o)
	}, UpdateFunc: func(_, o interface{}) {
		eventCount.WithLabelValues("ingress", "update").Inc()
		s.enqueue(o)
	}})
	return s
}
func (i *Ingress) enqueue(obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	i.queue.Add(key)
}
func (i *Ingress) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer i.queue.ShutDown()
	if !cache.WaitForCacheSync(ctx.Done(), i.informer.HasSynced) {
		level.Error(i.logger).Log("msg", "ingress informer unable to sync cache")
		return
	}
	go func() {
		for i.process(ctx, ch) {
		}
	}()
	<-ctx.Done()
}
func (i *Ingress) process(ctx context.Context, ch chan<- []*targetgroup.Group) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	keyObj, quit := i.queue.Get()
	if quit {
		return false
	}
	defer i.queue.Done(keyObj)
	key := keyObj.(string)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return true
	}
	o, exists, err := i.store.GetByKey(key)
	if err != nil {
		return true
	}
	if !exists {
		send(ctx, i.logger, RoleIngress, ch, &targetgroup.Group{Source: ingressSourceFromNamespaceAndName(namespace, name)})
		return true
	}
	eps, err := convertToIngress(o)
	if err != nil {
		level.Error(i.logger).Log("msg", "converting to Ingress object failed", "err", err)
		return true
	}
	send(ctx, i.logger, RoleIngress, ch, i.buildIngress(eps))
	return true
}
func convertToIngress(o interface{}) (*v1beta1.Ingress, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ingress, ok := o.(*v1beta1.Ingress)
	if ok {
		return ingress, nil
	}
	return nil, fmt.Errorf("received unexpected object: %v", o)
}
func ingressSource(s *v1beta1.Ingress) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ingressSourceFromNamespaceAndName(s.Namespace, s.Name)
}
func ingressSourceFromNamespaceAndName(namespace, name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "ingress/" + namespace + "/" + name
}

const (
	ingressNameLabel		= metaLabelPrefix + "ingress_name"
	ingressLabelPrefix		= metaLabelPrefix + "ingress_label_"
	ingressAnnotationPrefix	= metaLabelPrefix + "ingress_annotation_"
	ingressSchemeLabel		= metaLabelPrefix + "ingress_scheme"
	ingressHostLabel		= metaLabelPrefix + "ingress_host"
	ingressPathLabel		= metaLabelPrefix + "ingress_path"
)

func ingressLabels(ingress *v1beta1.Ingress) model.LabelSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ls := make(model.LabelSet, len(ingress.Labels)+len(ingress.Annotations)+2)
	ls[ingressNameLabel] = lv(ingress.Name)
	ls[namespaceLabel] = lv(ingress.Namespace)
	for k, v := range ingress.Labels {
		ln := strutil.SanitizeLabelName(ingressLabelPrefix + k)
		ls[model.LabelName(ln)] = lv(v)
	}
	for k, v := range ingress.Annotations {
		ln := strutil.SanitizeLabelName(ingressAnnotationPrefix + k)
		ls[model.LabelName(ln)] = lv(v)
	}
	return ls
}
func pathsFromIngressRule(rv *v1beta1.IngressRuleValue) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if rv.HTTP == nil {
		return []string{"/"}
	}
	paths := make([]string, len(rv.HTTP.Paths))
	for n, p := range rv.HTTP.Paths {
		path := p.Path
		if path == "" {
			path = "/"
		}
		paths[n] = path
	}
	return paths
}
func (i *Ingress) buildIngress(ingress *v1beta1.Ingress) *targetgroup.Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tg := &targetgroup.Group{Source: ingressSource(ingress)}
	tg.Labels = ingressLabels(ingress)
	tlsHosts := make(map[string]struct{})
	for _, tls := range ingress.Spec.TLS {
		for _, host := range tls.Hosts {
			tlsHosts[host] = struct{}{}
		}
	}
	for _, rule := range ingress.Spec.Rules {
		paths := pathsFromIngressRule(&rule.IngressRuleValue)
		scheme := "http"
		_, isTLS := tlsHosts[rule.Host]
		if isTLS {
			scheme = "https"
		}
		for _, path := range paths {
			tg.Targets = append(tg.Targets, model.LabelSet{model.AddressLabel: lv(rule.Host), ingressSchemeLabel: lv(scheme), ingressHostLabel: lv(rule.Host), ingressPathLabel: lv(path)})
		}
	}
	return tg
}
