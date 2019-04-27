package kubernetes

import (
	"fmt"
	"testing"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func makeOptionalBool(v bool) *bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &v
}
func makeMultiPortPods() *v1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "testpod", Namespace: "default", Labels: map[string]string{"testlabel": "testvalue"}, Annotations: map[string]string{"testannotation": "testannotationvalue"}, UID: types.UID("abc123"), OwnerReferences: []metav1.OwnerReference{{Kind: "testcontrollerkind", Name: "testcontrollername", Controller: makeOptionalBool(true)}}}, Spec: v1.PodSpec{NodeName: "testnode", Containers: []v1.Container{{Name: "testcontainer0", Ports: []v1.ContainerPort{{Name: "testport0", Protocol: v1.ProtocolTCP, ContainerPort: int32(9000)}, {Name: "testport1", Protocol: v1.ProtocolUDP, ContainerPort: int32(9001)}}}, {Name: "testcontainer1"}}}, Status: v1.PodStatus{PodIP: "1.2.3.4", HostIP: "2.3.4.5", Phase: "Running", Conditions: []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionTrue}}}}
}
func makePods() *v1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "testpod", Namespace: "default", UID: types.UID("abc123")}, Spec: v1.PodSpec{NodeName: "testnode", Containers: []v1.Container{{Name: "testcontainer", Ports: []v1.ContainerPort{{Name: "testport", Protocol: v1.ProtocolTCP, ContainerPort: int32(9000)}}}}}, Status: v1.PodStatus{PodIP: "1.2.3.4", HostIP: "2.3.4.5", Phase: "Running", Conditions: []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionTrue}}}}
}
func expectedPodTargetGroups(ns string) map[string]*targetgroup.Group {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	key := fmt.Sprintf("pod/%s/testpod", ns)
	return map[string]*targetgroup.Group{key: {Targets: []model.LabelSet{{"__address__": "1.2.3.4:9000", "__meta_kubernetes_pod_container_name": "testcontainer", "__meta_kubernetes_pod_container_port_name": "testport", "__meta_kubernetes_pod_container_port_number": "9000", "__meta_kubernetes_pod_container_port_protocol": "TCP"}}, Labels: model.LabelSet{"__meta_kubernetes_pod_name": "testpod", "__meta_kubernetes_namespace": lv(ns), "__meta_kubernetes_pod_node_name": "testnode", "__meta_kubernetes_pod_ip": "1.2.3.4", "__meta_kubernetes_pod_host_ip": "2.3.4.5", "__meta_kubernetes_pod_ready": "true", "__meta_kubernetes_pod_phase": "Running", "__meta_kubernetes_pod_uid": "abc123"}, Source: key}}
}
func TestPodDiscoveryBeforeRun(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, c, w := makeDiscovery(RolePod, NamespaceDiscovery{})
	k8sDiscoveryTest{discovery: n, beforeRun: func() {
		obj := makeMultiPortPods()
		c.CoreV1().Pods(obj.Namespace).Create(obj)
		w.Pods().Add(obj)
	}, expectedMaxItems: 1, expectedRes: map[string]*targetgroup.Group{"pod/default/testpod": {Targets: []model.LabelSet{{"__address__": "1.2.3.4:9000", "__meta_kubernetes_pod_container_name": "testcontainer0", "__meta_kubernetes_pod_container_port_name": "testport0", "__meta_kubernetes_pod_container_port_number": "9000", "__meta_kubernetes_pod_container_port_protocol": "TCP"}, {"__address__": "1.2.3.4:9001", "__meta_kubernetes_pod_container_name": "testcontainer0", "__meta_kubernetes_pod_container_port_name": "testport1", "__meta_kubernetes_pod_container_port_number": "9001", "__meta_kubernetes_pod_container_port_protocol": "UDP"}, {"__address__": "1.2.3.4", "__meta_kubernetes_pod_container_name": "testcontainer1"}}, Labels: model.LabelSet{"__meta_kubernetes_pod_name": "testpod", "__meta_kubernetes_namespace": "default", "__meta_kubernetes_pod_label_testlabel": "testvalue", "__meta_kubernetes_pod_annotation_testannotation": "testannotationvalue", "__meta_kubernetes_pod_node_name": "testnode", "__meta_kubernetes_pod_ip": "1.2.3.4", "__meta_kubernetes_pod_host_ip": "2.3.4.5", "__meta_kubernetes_pod_ready": "true", "__meta_kubernetes_pod_phase": "Running", "__meta_kubernetes_pod_uid": "abc123", "__meta_kubernetes_pod_controller_kind": "testcontrollerkind", "__meta_kubernetes_pod_controller_name": "testcontrollername"}, Source: "pod/default/testpod"}}}.Run(t)
}
func TestPodDiscoveryAdd(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, c, w := makeDiscovery(RolePod, NamespaceDiscovery{})
	k8sDiscoveryTest{discovery: n, afterStart: func() {
		obj := makePods()
		c.CoreV1().Pods(obj.Namespace).Create(obj)
		w.Pods().Add(obj)
	}, expectedMaxItems: 1, expectedRes: expectedPodTargetGroups("default")}.Run(t)
}
func TestPodDiscoveryDelete(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj := makePods()
	n, c, w := makeDiscovery(RolePod, NamespaceDiscovery{}, obj)
	k8sDiscoveryTest{discovery: n, afterStart: func() {
		obj := makePods()
		c.CoreV1().Pods(obj.Namespace).Delete(obj.Name, &metav1.DeleteOptions{})
		w.Pods().Delete(obj)
	}, expectedMaxItems: 2, expectedRes: map[string]*targetgroup.Group{"pod/default/testpod": {Source: "pod/default/testpod"}}}.Run(t)
}
func TestPodDiscoveryUpdate(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "testpod", Namespace: "default", UID: "xyz321"}, Spec: v1.PodSpec{NodeName: "testnode", Containers: []v1.Container{{Name: "testcontainer", Ports: []v1.ContainerPort{{Name: "testport", Protocol: v1.ProtocolTCP, ContainerPort: int32(9000)}}}}}, Status: v1.PodStatus{PodIP: "1.2.3.4", HostIP: "2.3.4.5"}}
	n, c, w := makeDiscovery(RolePod, NamespaceDiscovery{}, obj)
	k8sDiscoveryTest{discovery: n, afterStart: func() {
		obj := makePods()
		c.CoreV1().Pods(obj.Namespace).Create(obj)
		w.Pods().Modify(obj)
	}, expectedMaxItems: 2, expectedRes: expectedPodTargetGroups("default")}.Run(t)
}
func TestPodDiscoveryUpdateEmptyPodIP(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, c, w := makeDiscovery(RolePod, NamespaceDiscovery{})
	initialPod := makePods()
	updatedPod := makePods()
	updatedPod.Status.PodIP = ""
	k8sDiscoveryTest{discovery: n, beforeRun: func() {
		c.CoreV1().Pods(initialPod.Namespace).Create(initialPod)
		w.Pods().Add(initialPod)
	}, afterStart: func() {
		c.CoreV1().Pods(updatedPod.Namespace).Create(updatedPod)
		w.Pods().Modify(updatedPod)
	}, expectedMaxItems: 2, expectedRes: map[string]*targetgroup.Group{"pod/default/testpod": {Source: "pod/default/testpod"}}}.Run(t)
}
func TestPodDiscoveryNamespaces(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, c, w := makeDiscovery(RolePod, NamespaceDiscovery{Names: []string{"ns1", "ns2"}})
	expected := expectedPodTargetGroups("ns1")
	for k, v := range expectedPodTargetGroups("ns2") {
		expected[k] = v
	}
	k8sDiscoveryTest{discovery: n, beforeRun: func() {
		for _, ns := range []string{"ns1", "ns2"} {
			pod := makePods()
			pod.Namespace = ns
			c.CoreV1().Pods(pod.Namespace).Create(pod)
			w.Pods().Add(pod)
		}
	}, expectedMaxItems: 2, expectedRes: expected}.Run(t)
}
