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
	"github.com/prometheus/prometheus/util/strutil"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	NodeLegacyHostIP = "LegacyHostIP"
)

type Node struct {
	logger		log.Logger
	informer	cache.SharedInformer
	store		cache.Store
	queue		*workqueue.Type
}

func NewNode(l log.Logger, inf cache.SharedInformer) *Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if l == nil {
		l = log.NewNopLogger()
	}
	n := &Node{logger: l, informer: inf, store: inf.GetStore(), queue: workqueue.NewNamed("node")}
	n.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{AddFunc: func(o interface{}) {
		eventCount.WithLabelValues("node", "add").Inc()
		n.enqueue(o)
	}, DeleteFunc: func(o interface{}) {
		eventCount.WithLabelValues("node", "delete").Inc()
		n.enqueue(o)
	}, UpdateFunc: func(_, o interface{}) {
		eventCount.WithLabelValues("node", "update").Inc()
		n.enqueue(o)
	}})
	return n
}
func (n *Node) enqueue(obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	n.queue.Add(key)
}
func (n *Node) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer n.queue.ShutDown()
	if !cache.WaitForCacheSync(ctx.Done(), n.informer.HasSynced) {
		level.Error(n.logger).Log("msg", "node informer unable to sync cache")
		return
	}
	go func() {
		for n.process(ctx, ch) {
		}
	}()
	<-ctx.Done()
}
func (n *Node) process(ctx context.Context, ch chan<- []*targetgroup.Group) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	keyObj, quit := n.queue.Get()
	if quit {
		return false
	}
	defer n.queue.Done(keyObj)
	key := keyObj.(string)
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return true
	}
	o, exists, err := n.store.GetByKey(key)
	if err != nil {
		return true
	}
	if !exists {
		send(ctx, n.logger, RoleNode, ch, &targetgroup.Group{Source: nodeSourceFromName(name)})
		return true
	}
	node, err := convertToNode(o)
	if err != nil {
		level.Error(n.logger).Log("msg", "converting to Node object failed", "err", err)
		return true
	}
	send(ctx, n.logger, RoleNode, ch, n.buildNode(node))
	return true
}
func convertToNode(o interface{}) (*apiv1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	node, ok := o.(*apiv1.Node)
	if ok {
		return node, nil
	}
	return nil, fmt.Errorf("received unexpected object: %v", o)
}
func nodeSource(n *apiv1.Node) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nodeSourceFromName(n.Name)
}
func nodeSourceFromName(name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "node/" + name
}

const (
	nodeNameLabel		= metaLabelPrefix + "node_name"
	nodeLabelPrefix		= metaLabelPrefix + "node_label_"
	nodeAnnotationPrefix	= metaLabelPrefix + "node_annotation_"
	nodeAddressPrefix	= metaLabelPrefix + "node_address_"
)

func nodeLabels(n *apiv1.Node) model.LabelSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ls := make(model.LabelSet, len(n.Labels)+len(n.Annotations)+1)
	ls[nodeNameLabel] = lv(n.Name)
	for k, v := range n.Labels {
		ln := strutil.SanitizeLabelName(nodeLabelPrefix + k)
		ls[model.LabelName(ln)] = lv(v)
	}
	for k, v := range n.Annotations {
		ln := strutil.SanitizeLabelName(nodeAnnotationPrefix + k)
		ls[model.LabelName(ln)] = lv(v)
	}
	return ls
}
func (n *Node) buildNode(node *apiv1.Node) *targetgroup.Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tg := &targetgroup.Group{Source: nodeSource(node)}
	tg.Labels = nodeLabels(node)
	addr, addrMap, err := nodeAddress(node)
	if err != nil {
		level.Warn(n.logger).Log("msg", "No node address found", "err", err)
		return nil
	}
	addr = net.JoinHostPort(addr, strconv.FormatInt(int64(node.Status.DaemonEndpoints.KubeletEndpoint.Port), 10))
	t := model.LabelSet{model.AddressLabel: lv(addr), model.InstanceLabel: lv(node.Name)}
	for ty, a := range addrMap {
		ln := strutil.SanitizeLabelName(nodeAddressPrefix + string(ty))
		t[model.LabelName(ln)] = lv(a[0])
	}
	tg.Targets = append(tg.Targets, t)
	return tg
}
func nodeAddress(node *apiv1.Node) (string, map[apiv1.NodeAddressType][]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := map[apiv1.NodeAddressType][]string{}
	for _, a := range node.Status.Addresses {
		m[a.Type] = append(m[a.Type], a.Address)
	}
	if addresses, ok := m[apiv1.NodeInternalIP]; ok {
		return addresses[0], m, nil
	}
	if addresses, ok := m[apiv1.NodeExternalIP]; ok {
		return addresses[0], m, nil
	}
	if addresses, ok := m[apiv1.NodeAddressType(NodeLegacyHostIP)]; ok {
		return addresses[0], m, nil
	}
	if addresses, ok := m[apiv1.NodeHostName]; ok {
		return addresses[0], m, nil
	}
	return "", m, fmt.Errorf("host address unknown")
}
