package kubernetes

import (
	"fmt"
	"testing"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeNode(name, address string, labels map[string]string, annotations map[string]string) *v1.Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &v1.Node{ObjectMeta: metav1.ObjectMeta{Name: name, Labels: labels, Annotations: annotations}, Status: v1.NodeStatus{Addresses: []v1.NodeAddress{{Type: v1.NodeInternalIP, Address: address}}, DaemonEndpoints: v1.NodeDaemonEndpoints{KubeletEndpoint: v1.DaemonEndpoint{Port: 10250}}}}
}
func makeEnumeratedNode(i int) *v1.Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return makeNode(fmt.Sprintf("test%d", i), "1.2.3.4", map[string]string{}, map[string]string{})
}
func TestNodeDiscoveryBeforeStart(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, c, w := makeDiscovery(RoleNode, NamespaceDiscovery{})
	k8sDiscoveryTest{discovery: n, beforeRun: func() {
		obj := makeNode("test", "1.2.3.4", map[string]string{"testlabel": "testvalue"}, map[string]string{"testannotation": "testannotationvalue"})
		c.CoreV1().Nodes().Create(obj)
		w.Nodes().Add(obj)
	}, expectedMaxItems: 1, expectedRes: map[string]*targetgroup.Group{"node/test": {Targets: []model.LabelSet{{"__address__": "1.2.3.4:10250", "instance": "test", "__meta_kubernetes_node_address_InternalIP": "1.2.3.4"}}, Labels: model.LabelSet{"__meta_kubernetes_node_name": "test", "__meta_kubernetes_node_label_testlabel": "testvalue", "__meta_kubernetes_node_annotation_testannotation": "testannotationvalue"}, Source: "node/test"}}}.Run(t)
}
func TestNodeDiscoveryAdd(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, c, w := makeDiscovery(RoleNode, NamespaceDiscovery{})
	k8sDiscoveryTest{discovery: n, afterStart: func() {
		obj := makeEnumeratedNode(1)
		c.CoreV1().Nodes().Create(obj)
		w.Nodes().Add(obj)
	}, expectedMaxItems: 1, expectedRes: map[string]*targetgroup.Group{"node/test1": {Targets: []model.LabelSet{{"__address__": "1.2.3.4:10250", "instance": "test1", "__meta_kubernetes_node_address_InternalIP": "1.2.3.4"}}, Labels: model.LabelSet{"__meta_kubernetes_node_name": "test1"}, Source: "node/test1"}}}.Run(t)
}
func TestNodeDiscoveryDelete(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj := makeEnumeratedNode(0)
	n, c, w := makeDiscovery(RoleNode, NamespaceDiscovery{}, obj)
	k8sDiscoveryTest{discovery: n, afterStart: func() {
		c.CoreV1().Nodes().Delete(obj.Name, &metav1.DeleteOptions{})
		w.Nodes().Delete(obj)
	}, expectedMaxItems: 2, expectedRes: map[string]*targetgroup.Group{"node/test0": {Source: "node/test0"}}}.Run(t)
}
func TestNodeDiscoveryUpdate(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, c, w := makeDiscovery(RoleNode, NamespaceDiscovery{})
	k8sDiscoveryTest{discovery: n, afterStart: func() {
		obj1 := makeEnumeratedNode(0)
		c.CoreV1().Nodes().Create(obj1)
		w.Nodes().Add(obj1)
		obj2 := makeNode("test0", "1.2.3.4", map[string]string{"Unschedulable": "true"}, map[string]string{})
		c.CoreV1().Nodes().Update(obj2)
		w.Nodes().Modify(obj2)
	}, expectedMaxItems: 2, expectedRes: map[string]*targetgroup.Group{"node/test0": {Targets: []model.LabelSet{{"__address__": "1.2.3.4:10250", "instance": "test0", "__meta_kubernetes_node_address_InternalIP": "1.2.3.4"}}, Labels: model.LabelSet{"__meta_kubernetes_node_label_Unschedulable": "true", "__meta_kubernetes_node_name": "test0"}, Source: "node/test0"}}}.Run(t)
}
