package kubernetes

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/util/strutil"
)

type Pod struct {
	informer	cache.SharedInformer
	store		cache.Store
	logger		log.Logger
	queue		*workqueue.Type
}

func NewPod(l log.Logger, pods cache.SharedInformer) *Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if l == nil {
		l = log.NewNopLogger()
	}
	p := &Pod{informer: pods, store: pods.GetStore(), logger: l, queue: workqueue.NewNamed("pod")}
	p.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{AddFunc: func(o interface{}) {
		eventCount.WithLabelValues("pod", "add").Inc()
		p.enqueue(o)
	}, DeleteFunc: func(o interface{}) {
		eventCount.WithLabelValues("pod", "delete").Inc()
		p.enqueue(o)
	}, UpdateFunc: func(_, o interface{}) {
		eventCount.WithLabelValues("pod", "update").Inc()
		p.enqueue(o)
	}})
	return p
}
func (p *Pod) enqueue(obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	p.queue.Add(key)
}
func (p *Pod) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer p.queue.ShutDown()
	if !cache.WaitForCacheSync(ctx.Done(), p.informer.HasSynced) {
		level.Error(p.logger).Log("msg", "pod informer unable to sync cache")
		return
	}
	go func() {
		for p.process(ctx, ch) {
		}
	}()
	<-ctx.Done()
}
func (p *Pod) process(ctx context.Context, ch chan<- []*targetgroup.Group) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	keyObj, quit := p.queue.Get()
	if quit {
		return false
	}
	defer p.queue.Done(keyObj)
	key := keyObj.(string)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return true
	}
	o, exists, err := p.store.GetByKey(key)
	if err != nil {
		return true
	}
	if !exists {
		send(ctx, p.logger, RolePod, ch, &targetgroup.Group{Source: podSourceFromNamespaceAndName(namespace, name)})
		return true
	}
	eps, err := convertToPod(o)
	if err != nil {
		level.Error(p.logger).Log("msg", "converting to Pod object failed", "err", err)
		return true
	}
	send(ctx, p.logger, RolePod, ch, p.buildPod(eps))
	return true
}
func convertToPod(o interface{}) (*apiv1.Pod, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pod, ok := o.(*apiv1.Pod)
	if ok {
		return pod, nil
	}
	return nil, fmt.Errorf("received unexpected object: %v", o)
}

const (
	podNameLabel					= metaLabelPrefix + "pod_name"
	podIPLabel						= metaLabelPrefix + "pod_ip"
	podContainerNameLabel			= metaLabelPrefix + "pod_container_name"
	podContainerPortNameLabel		= metaLabelPrefix + "pod_container_port_name"
	podContainerPortNumberLabel		= metaLabelPrefix + "pod_container_port_number"
	podContainerPortProtocolLabel	= metaLabelPrefix + "pod_container_port_protocol"
	podReadyLabel					= metaLabelPrefix + "pod_ready"
	podPhaseLabel					= metaLabelPrefix + "pod_phase"
	podLabelPrefix					= metaLabelPrefix + "pod_label_"
	podAnnotationPrefix				= metaLabelPrefix + "pod_annotation_"
	podNodeNameLabel				= metaLabelPrefix + "pod_node_name"
	podHostIPLabel					= metaLabelPrefix + "pod_host_ip"
	podUID							= metaLabelPrefix + "pod_uid"
	podControllerKind				= metaLabelPrefix + "pod_controller_kind"
	podControllerName				= metaLabelPrefix + "pod_controller_name"
)

func GetControllerOf(controllee metav1.Object) *metav1.OwnerReference {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, ref := range controllee.GetOwnerReferences() {
		if ref.Controller != nil && *ref.Controller {
			return &ref
		}
	}
	return nil
}
func podLabels(pod *apiv1.Pod) model.LabelSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ls := model.LabelSet{podNameLabel: lv(pod.ObjectMeta.Name), podIPLabel: lv(pod.Status.PodIP), podReadyLabel: podReady(pod), podPhaseLabel: lv(string(pod.Status.Phase)), podNodeNameLabel: lv(pod.Spec.NodeName), podHostIPLabel: lv(pod.Status.HostIP), podUID: lv(string(pod.ObjectMeta.UID))}
	createdBy := GetControllerOf(pod)
	if createdBy != nil {
		if createdBy.Kind != "" {
			ls[podControllerKind] = lv(createdBy.Kind)
		}
		if createdBy.Name != "" {
			ls[podControllerName] = lv(createdBy.Name)
		}
	}
	for k, v := range pod.Labels {
		ln := strutil.SanitizeLabelName(podLabelPrefix + k)
		ls[model.LabelName(ln)] = lv(v)
	}
	for k, v := range pod.Annotations {
		ln := strutil.SanitizeLabelName(podAnnotationPrefix + k)
		ls[model.LabelName(ln)] = lv(v)
	}
	return ls
}
func (p *Pod) buildPod(pod *apiv1.Pod) *targetgroup.Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tg := &targetgroup.Group{Source: podSource(pod)}
	if len(pod.Status.PodIP) == 0 {
		return tg
	}
	tg.Labels = podLabels(pod)
	tg.Labels[namespaceLabel] = lv(pod.Namespace)
	for _, c := range pod.Spec.Containers {
		if len(c.Ports) == 0 {
			tg.Targets = append(tg.Targets, model.LabelSet{model.AddressLabel: lv(pod.Status.PodIP), podContainerNameLabel: lv(c.Name)})
			continue
		}
		for _, port := range c.Ports {
			ports := strconv.FormatUint(uint64(port.ContainerPort), 10)
			addr := net.JoinHostPort(pod.Status.PodIP, ports)
			tg.Targets = append(tg.Targets, model.LabelSet{model.AddressLabel: lv(addr), podContainerNameLabel: lv(c.Name), podContainerPortNumberLabel: lv(ports), podContainerPortNameLabel: lv(port.Name), podContainerPortProtocolLabel: lv(string(port.Protocol))})
		}
	}
	return tg
}
func podSource(pod *apiv1.Pod) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return podSourceFromNamespaceAndName(pod.Namespace, pod.Name)
}
func podSourceFromNamespaceAndName(namespace, name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "pod/" + namespace + "/" + name
}
func podReady(pod *apiv1.Pod) model.LabelValue {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, cond := range pod.Status.Conditions {
		if cond.Type == apiv1.PodReady {
			return lv(strings.ToLower(string(cond.Status)))
		}
	}
	return lv(strings.ToLower(string(apiv1.ConditionUnknown)))
}
