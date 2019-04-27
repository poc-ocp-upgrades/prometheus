package kubernetes

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Endpoints struct {
	logger		log.Logger
	endpointsInf	cache.SharedInformer
	serviceInf	cache.SharedInformer
	podInf		cache.SharedInformer
	podStore	cache.Store
	endpointsStore	cache.Store
	serviceStore	cache.Store
	queue		*workqueue.Type
}

func NewEndpoints(l log.Logger, svc, eps, pod cache.SharedInformer) *Endpoints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if l == nil {
		l = log.NewNopLogger()
	}
	e := &Endpoints{logger: l, endpointsInf: eps, endpointsStore: eps.GetStore(), serviceInf: svc, serviceStore: svc.GetStore(), podInf: pod, podStore: pod.GetStore(), queue: workqueue.NewNamed("endpoints")}
	e.endpointsInf.AddEventHandler(cache.ResourceEventHandlerFuncs{AddFunc: func(o interface{}) {
		eventCount.WithLabelValues("endpoints", "add").Inc()
		e.enqueue(o)
	}, UpdateFunc: func(_, o interface{}) {
		eventCount.WithLabelValues("endpoints", "update").Inc()
		e.enqueue(o)
	}, DeleteFunc: func(o interface{}) {
		eventCount.WithLabelValues("endpoints", "delete").Inc()
		e.enqueue(o)
	}})
	serviceUpdate := func(o interface{}) {
		svc, err := convertToService(o)
		if err != nil {
			level.Error(e.logger).Log("msg", "converting to Service object failed", "err", err)
			return
		}
		ep := &apiv1.Endpoints{}
		ep.Namespace = svc.Namespace
		ep.Name = svc.Name
		obj, exists, err := e.endpointsStore.Get(ep)
		if exists && err == nil {
			e.enqueue(obj.(*apiv1.Endpoints))
		}
		if err != nil {
			level.Error(e.logger).Log("msg", "retrieving endpoints failed", "err", err)
		}
	}
	e.serviceInf.AddEventHandler(cache.ResourceEventHandlerFuncs{AddFunc: func(o interface{}) {
		eventCount.WithLabelValues("service", "add").Inc()
		serviceUpdate(o)
	}, UpdateFunc: func(_, o interface{}) {
		eventCount.WithLabelValues("service", "update").Inc()
		serviceUpdate(o)
	}, DeleteFunc: func(o interface{}) {
		eventCount.WithLabelValues("service", "delete").Inc()
		serviceUpdate(o)
	}})
	return e
}
func (e *Endpoints) enqueue(obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	e.queue.Add(key)
}
func (e *Endpoints) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer e.queue.ShutDown()
	if !cache.WaitForCacheSync(ctx.Done(), e.endpointsInf.HasSynced, e.serviceInf.HasSynced, e.podInf.HasSynced) {
		level.Error(e.logger).Log("msg", "endpoints informer unable to sync cache")
		return
	}
	go func() {
		for e.process(ctx, ch) {
		}
	}()
	<-ctx.Done()
}
func (e *Endpoints) process(ctx context.Context, ch chan<- []*targetgroup.Group) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	keyObj, quit := e.queue.Get()
	if quit {
		return false
	}
	defer e.queue.Done(keyObj)
	key := keyObj.(string)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		level.Error(e.logger).Log("msg", "splitting key failed", "key", key)
		return true
	}
	o, exists, err := e.endpointsStore.GetByKey(key)
	if err != nil {
		level.Error(e.logger).Log("msg", "getting object from store failed", "key", key)
		return true
	}
	if !exists {
		send(ctx, e.logger, RoleEndpoint, ch, &targetgroup.Group{Source: endpointsSourceFromNamespaceAndName(namespace, name)})
		return true
	}
	eps, err := convertToEndpoints(o)
	if err != nil {
		level.Error(e.logger).Log("msg", "converting to Endpoints object failed", "err", err)
		return true
	}
	send(ctx, e.logger, RoleEndpoint, ch, e.buildEndpoints(eps))
	return true
}
func convertToEndpoints(o interface{}) (*apiv1.Endpoints, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	endpoints, ok := o.(*apiv1.Endpoints)
	if ok {
		return endpoints, nil
	}
	return nil, fmt.Errorf("received unexpected object: %v", o)
}
func endpointsSource(ep *apiv1.Endpoints) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return endpointsSourceFromNamespaceAndName(ep.Namespace, ep.Name)
}
func endpointsSourceFromNamespaceAndName(namespace, name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "endpoints/" + namespace + "/" + name
}

const (
	endpointsNameLabel		= metaLabelPrefix + "endpoints_name"
	endpointReadyLabel		= metaLabelPrefix + "endpoint_ready"
	endpointPortNameLabel		= metaLabelPrefix + "endpoint_port_name"
	endpointPortProtocolLabel	= metaLabelPrefix + "endpoint_port_protocol"
	endpointAddressTargetKindLabel	= metaLabelPrefix + "endpoint_address_target_kind"
	endpointAddressTargetNameLabel	= metaLabelPrefix + "endpoint_address_target_name"
)

func (e *Endpoints) buildEndpoints(eps *apiv1.Endpoints) *targetgroup.Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tg := &targetgroup.Group{Source: endpointsSource(eps)}
	tg.Labels = model.LabelSet{namespaceLabel: lv(eps.Namespace), endpointsNameLabel: lv(eps.Name)}
	e.addServiceLabels(eps.Namespace, eps.Name, tg)
	type podEntry struct {
		pod		*apiv1.Pod
		servicePorts	[]apiv1.EndpointPort
	}
	seenPods := map[string]*podEntry{}
	add := func(addr apiv1.EndpointAddress, port apiv1.EndpointPort, ready string) {
		a := net.JoinHostPort(addr.IP, strconv.FormatUint(uint64(port.Port), 10))
		target := model.LabelSet{model.AddressLabel: lv(a), endpointPortNameLabel: lv(port.Name), endpointPortProtocolLabel: lv(string(port.Protocol)), endpointReadyLabel: lv(ready)}
		if addr.TargetRef != nil {
			target[model.LabelName(endpointAddressTargetKindLabel)] = lv(addr.TargetRef.Kind)
			target[model.LabelName(endpointAddressTargetNameLabel)] = lv(addr.TargetRef.Name)
		}
		pod := e.resolvePodRef(addr.TargetRef)
		if pod == nil {
			tg.Targets = append(tg.Targets, target)
			return
		}
		s := pod.Namespace + "/" + pod.Name
		sp, ok := seenPods[s]
		if !ok {
			sp = &podEntry{pod: pod}
			seenPods[s] = sp
		}
		target = target.Merge(podLabels(pod))
		for _, c := range pod.Spec.Containers {
			for _, cport := range c.Ports {
				if port.Port == cport.ContainerPort {
					ports := strconv.FormatUint(uint64(port.Port), 10)
					target[podContainerNameLabel] = lv(c.Name)
					target[podContainerPortNameLabel] = lv(cport.Name)
					target[podContainerPortNumberLabel] = lv(ports)
					target[podContainerPortProtocolLabel] = lv(string(port.Protocol))
					break
				}
			}
		}
		sp.servicePorts = append(sp.servicePorts, port)
		tg.Targets = append(tg.Targets, target)
	}
	for _, ss := range eps.Subsets {
		for _, port := range ss.Ports {
			for _, addr := range ss.Addresses {
				add(addr, port, "true")
			}
			for _, addr := range ss.NotReadyAddresses {
				add(addr, port, "false")
			}
		}
	}
	for _, pe := range seenPods {
		for _, c := range pe.pod.Spec.Containers {
			for _, cport := range c.Ports {
				hasSeenPort := func() bool {
					for _, eport := range pe.servicePorts {
						if cport.ContainerPort == eport.Port {
							return true
						}
					}
					return false
				}
				if hasSeenPort() {
					continue
				}
				a := net.JoinHostPort(pe.pod.Status.PodIP, strconv.FormatUint(uint64(cport.ContainerPort), 10))
				ports := strconv.FormatUint(uint64(cport.ContainerPort), 10)
				target := model.LabelSet{model.AddressLabel: lv(a), podContainerNameLabel: lv(c.Name), podContainerPortNameLabel: lv(cport.Name), podContainerPortNumberLabel: lv(ports), podContainerPortProtocolLabel: lv(string(cport.Protocol))}
				tg.Targets = append(tg.Targets, target.Merge(podLabels(pe.pod)))
			}
		}
	}
	return tg
}
func (e *Endpoints) resolvePodRef(ref *apiv1.ObjectReference) *apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if ref == nil || ref.Kind != "Pod" {
		return nil
	}
	p := &apiv1.Pod{}
	p.Namespace = ref.Namespace
	p.Name = ref.Name
	obj, exists, err := e.podStore.Get(p)
	if err != nil || !exists {
		return nil
	}
	if err != nil {
		level.Error(e.logger).Log("msg", "resolving pod ref failed", "err", err)
	}
	return obj.(*apiv1.Pod)
}
func (e *Endpoints) addServiceLabels(ns, name string, tg *targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	svc := &apiv1.Service{}
	svc.Namespace = ns
	svc.Name = name
	obj, exists, err := e.serviceStore.Get(svc)
	if !exists || err != nil {
		return
	}
	if err != nil {
		level.Error(e.logger).Log("msg", "retrieving service failed", "err", err)
	}
	svc = obj.(*apiv1.Service)
	tg.Labels = tg.Labels.Merge(serviceLabels(svc))
}
