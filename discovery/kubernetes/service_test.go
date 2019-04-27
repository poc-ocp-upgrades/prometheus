package kubernetes

import (
	"fmt"
	"testing"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeMultiPortService() *v1.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: "testservice", Namespace: "default", Labels: map[string]string{"testlabel": "testvalue"}, Annotations: map[string]string{"testannotation": "testannotationvalue"}}, Spec: v1.ServiceSpec{Ports: []v1.ServicePort{{Name: "testport0", Protocol: v1.ProtocolTCP, Port: int32(30900)}, {Name: "testport1", Protocol: v1.ProtocolUDP, Port: int32(30901)}}, Type: v1.ServiceTypeClusterIP, ClusterIP: "10.0.0.1"}}
}
func makeSuffixedService(suffix string) *v1.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("testservice%s", suffix), Namespace: "default"}, Spec: v1.ServiceSpec{Ports: []v1.ServicePort{{Name: "testport", Protocol: v1.ProtocolTCP, Port: int32(30900)}}, Type: v1.ServiceTypeClusterIP, ClusterIP: "10.0.0.1"}}
}
func makeService() *v1.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return makeSuffixedService("")
}
func makeExternalService() *v1.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: "testservice-external", Namespace: "default"}, Spec: v1.ServiceSpec{Ports: []v1.ServicePort{{Name: "testport", Protocol: v1.ProtocolTCP, Port: int32(31900)}}, Type: v1.ServiceTypeExternalName, ExternalName: "FooExternalName"}}
}
func TestServiceDiscoveryAdd(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, c, w := makeDiscovery(RoleService, NamespaceDiscovery{})
	k8sDiscoveryTest{discovery: n, afterStart: func() {
		obj := makeService()
		c.CoreV1().Services(obj.Namespace).Create(obj)
		w.Services().Add(obj)
		obj = makeExternalService()
		c.CoreV1().Services(obj.Namespace).Create(obj)
		w.Services().Add(obj)
	}, expectedMaxItems: 2, expectedRes: map[string]*targetgroup.Group{"svc/default/testservice": {Targets: []model.LabelSet{{"__meta_kubernetes_service_port_protocol": "TCP", "__address__": "testservice.default.svc:30900", "__meta_kubernetes_service_cluster_ip": "10.0.0.1", "__meta_kubernetes_service_port_name": "testport"}}, Labels: model.LabelSet{"__meta_kubernetes_service_name": "testservice", "__meta_kubernetes_namespace": "default"}, Source: "svc/default/testservice"}, "svc/default/testservice-external": {Targets: []model.LabelSet{{"__meta_kubernetes_service_port_protocol": "TCP", "__address__": "testservice-external.default.svc:31900", "__meta_kubernetes_service_port_name": "testport", "__meta_kubernetes_service_external_name": "FooExternalName"}}, Labels: model.LabelSet{"__meta_kubernetes_service_name": "testservice-external", "__meta_kubernetes_namespace": "default"}, Source: "svc/default/testservice-external"}}}.Run(t)
}
func TestServiceDiscoveryDelete(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, c, w := makeDiscovery(RoleService, NamespaceDiscovery{}, makeService())
	k8sDiscoveryTest{discovery: n, afterStart: func() {
		obj := makeService()
		c.CoreV1().Services(obj.Namespace).Delete(obj.Name, &metav1.DeleteOptions{})
		w.Services().Delete(obj)
	}, expectedMaxItems: 2, expectedRes: map[string]*targetgroup.Group{"svc/default/testservice": {Source: "svc/default/testservice"}}}.Run(t)
}
func TestServiceDiscoveryUpdate(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, c, w := makeDiscovery(RoleService, NamespaceDiscovery{}, makeService())
	k8sDiscoveryTest{discovery: n, afterStart: func() {
		obj := makeMultiPortService()
		c.CoreV1().Services(obj.Namespace).Update(obj)
		w.Services().Modify(obj)
	}, expectedMaxItems: 2, expectedRes: map[string]*targetgroup.Group{"svc/default/testservice": {Targets: []model.LabelSet{{"__meta_kubernetes_service_port_protocol": "TCP", "__address__": "testservice.default.svc:30900", "__meta_kubernetes_service_cluster_ip": "10.0.0.1", "__meta_kubernetes_service_port_name": "testport0"}, {"__meta_kubernetes_service_port_protocol": "UDP", "__address__": "testservice.default.svc:30901", "__meta_kubernetes_service_cluster_ip": "10.0.0.1", "__meta_kubernetes_service_port_name": "testport1"}}, Labels: model.LabelSet{"__meta_kubernetes_service_name": "testservice", "__meta_kubernetes_namespace": "default", "__meta_kubernetes_service_label_testlabel": "testvalue", "__meta_kubernetes_service_annotation_testannotation": "testannotationvalue"}, Source: "svc/default/testservice"}}}.Run(t)
}
func TestServiceDiscoveryNamespaces(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, c, w := makeDiscovery(RoleService, NamespaceDiscovery{Names: []string{"ns1", "ns2"}})
	k8sDiscoveryTest{discovery: n, afterStart: func() {
		for _, ns := range []string{"ns1", "ns2"} {
			obj := makeService()
			obj.Namespace = ns
			c.CoreV1().Services(obj.Namespace).Create(obj)
			w.Services().Add(obj)
		}
	}, expectedMaxItems: 2, expectedRes: map[string]*targetgroup.Group{"svc/ns1/testservice": {Targets: []model.LabelSet{{"__meta_kubernetes_service_port_protocol": "TCP", "__address__": "testservice.ns1.svc:30900", "__meta_kubernetes_service_cluster_ip": "10.0.0.1", "__meta_kubernetes_service_port_name": "testport"}}, Labels: model.LabelSet{"__meta_kubernetes_service_name": "testservice", "__meta_kubernetes_namespace": "ns1"}, Source: "svc/ns1/testservice"}, "svc/ns2/testservice": {Targets: []model.LabelSet{{"__meta_kubernetes_service_port_protocol": "TCP", "__address__": "testservice.ns2.svc:30900", "__meta_kubernetes_service_cluster_ip": "10.0.0.1", "__meta_kubernetes_service_port_name": "testport"}}, Labels: model.LabelSet{"__meta_kubernetes_service_name": "testservice", "__meta_kubernetes_namespace": "ns2"}, Source: "svc/ns2/testservice"}}}.Run(t)
}
