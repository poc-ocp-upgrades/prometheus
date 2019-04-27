package kubernetes

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/util/strutil"
)

type Service struct {
	logger		log.Logger
	informer	cache.SharedInformer
	store		cache.Store
	queue		*workqueue.Type
}

func NewService(l log.Logger, inf cache.SharedInformer) *Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if l == nil {
		l = log.NewNopLogger()
	}
	s := &Service{logger: l, informer: inf, store: inf.GetStore(), queue: workqueue.NewNamed("ingress")}
	s.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{AddFunc: func(o interface{}) {
		eventCount.WithLabelValues("service", "add").Inc()
		s.enqueue(o)
	}, DeleteFunc: func(o interface{}) {
		eventCount.WithLabelValues("service", "delete").Inc()
		s.enqueue(o)
	}, UpdateFunc: func(_, o interface{}) {
		eventCount.WithLabelValues("service", "update").Inc()
		s.enqueue(o)
	}})
	return s
}
func (s *Service) enqueue(obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	s.queue.Add(key)
}
func (s *Service) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer s.queue.ShutDown()
	if !cache.WaitForCacheSync(ctx.Done(), s.informer.HasSynced) {
		level.Error(s.logger).Log("msg", "service informer unable to sync cache")
		return
	}
	go func() {
		for s.process(ctx, ch) {
		}
	}()
	<-ctx.Done()
}
func (s *Service) process(ctx context.Context, ch chan<- []*targetgroup.Group) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	keyObj, quit := s.queue.Get()
	if quit {
		return false
	}
	defer s.queue.Done(keyObj)
	key := keyObj.(string)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return true
	}
	o, exists, err := s.store.GetByKey(key)
	if err != nil {
		return true
	}
	if !exists {
		send(ctx, s.logger, RoleService, ch, &targetgroup.Group{Source: serviceSourceFromNamespaceAndName(namespace, name)})
		return true
	}
	eps, err := convertToService(o)
	if err != nil {
		level.Error(s.logger).Log("msg", "converting to Service object failed", "err", err)
		return true
	}
	send(ctx, s.logger, RoleService, ch, s.buildService(eps))
	return true
}
func convertToService(o interface{}) (*apiv1.Service, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	service, ok := o.(*apiv1.Service)
	if ok {
		return service, nil
	}
	return nil, fmt.Errorf("received unexpected object: %v", o)
}
func serviceSource(s *apiv1.Service) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return serviceSourceFromNamespaceAndName(s.Namespace, s.Name)
}
func serviceSourceFromNamespaceAndName(namespace, name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "svc/" + namespace + "/" + name
}

const (
	serviceNameLabel		= metaLabelPrefix + "service_name"
	serviceLabelPrefix		= metaLabelPrefix + "service_label_"
	serviceAnnotationPrefix		= metaLabelPrefix + "service_annotation_"
	servicePortNameLabel		= metaLabelPrefix + "service_port_name"
	servicePortProtocolLabel	= metaLabelPrefix + "service_port_protocol"
	serviceClusterIPLabel		= metaLabelPrefix + "service_cluster_ip"
	serviceExternalNameLabel	= metaLabelPrefix + "service_external_name"
)

func serviceLabels(svc *apiv1.Service) model.LabelSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ls := make(model.LabelSet, len(svc.Labels)+len(svc.Annotations)+2)
	ls[serviceNameLabel] = lv(svc.Name)
	ls[namespaceLabel] = lv(svc.Namespace)
	for k, v := range svc.Labels {
		ln := strutil.SanitizeLabelName(serviceLabelPrefix + k)
		ls[model.LabelName(ln)] = lv(v)
	}
	for k, v := range svc.Annotations {
		ln := strutil.SanitizeLabelName(serviceAnnotationPrefix + k)
		ls[model.LabelName(ln)] = lv(v)
	}
	return ls
}
func (s *Service) buildService(svc *apiv1.Service) *targetgroup.Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tg := &targetgroup.Group{Source: serviceSource(svc)}
	tg.Labels = serviceLabels(svc)
	for _, port := range svc.Spec.Ports {
		addr := net.JoinHostPort(svc.Name+"."+svc.Namespace+".svc", strconv.FormatInt(int64(port.Port), 10))
		labelSet := model.LabelSet{model.AddressLabel: lv(addr), servicePortNameLabel: lv(port.Name), servicePortProtocolLabel: lv(string(port.Protocol))}
		if svc.Spec.Type == apiv1.ServiceTypeExternalName {
			labelSet[serviceExternalNameLabel] = lv(svc.Spec.ExternalName)
		} else {
			labelSet[serviceClusterIPLabel] = lv(svc.Spec.ClusterIP)
		}
		tg.Targets = append(tg.Targets, labelSet)
	}
	return tg
}
