package adapter

import (
	"reflect"
	"testing"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

func TestGenerateTargetGroups(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	testCases := []struct {
		title			string
		targetGroup		map[string][]*targetgroup.Group
		expectedCustomSD	map[string]*customSD
	}{{title: "Empty targetGroup", targetGroup: map[string][]*targetgroup.Group{"customSD": {{Source: "Consul"}, {Source: "Kubernetes"}}}, expectedCustomSD: map[string]*customSD{"customSD:Consul:0000000000000000": {Targets: []string{}, Labels: map[string]string{}}, "customSD:Kubernetes:0000000000000000": {Targets: []string{}, Labels: map[string]string{}}}}, {title: "targetGroup filled", targetGroup: map[string][]*targetgroup.Group{"customSD": {{Source: "Azure", Targets: []model.LabelSet{{model.AddressLabel: "host1"}, {model.AddressLabel: "host2"}}, Labels: model.LabelSet{model.LabelName("__meta_test_label"): model.LabelValue("label_test_1")}}, {Source: "Openshift", Targets: []model.LabelSet{{model.AddressLabel: "host3"}, {model.AddressLabel: "host4"}}, Labels: model.LabelSet{model.LabelName("__meta_test_label"): model.LabelValue("label_test_2")}}}}, expectedCustomSD: map[string]*customSD{"customSD:Azure:282a007a18fadbbb": {Targets: []string{"host1", "host2"}, Labels: map[string]string{"__meta_test_label": "label_test_1"}}, "customSD:Openshift:281c007a18ea2ad0": {Targets: []string{"host3", "host4"}, Labels: map[string]string{"__meta_test_label": "label_test_2"}}}}, {title: "Mixed between empty targetGroup and targetGroup filled", targetGroup: map[string][]*targetgroup.Group{"customSD": {{Source: "GCE", Targets: []model.LabelSet{{model.AddressLabel: "host1"}, {model.AddressLabel: "host2"}}, Labels: model.LabelSet{model.LabelName("__meta_test_label"): model.LabelValue("label_test_1")}}, {Source: "Kubernetes", Labels: model.LabelSet{model.LabelName("__meta_test_label"): model.LabelValue("label_test_2")}}}}, expectedCustomSD: map[string]*customSD{"customSD:GCE:282a007a18fadbbb": {Targets: []string{"host1", "host2"}, Labels: map[string]string{"__meta_test_label": "label_test_1"}}, "customSD:Kubernetes:282e007a18fad483": {Targets: []string{}, Labels: map[string]string{"__meta_test_label": "label_test_2"}}}}}
	for _, testCase := range testCases {
		result := generateTargetGroups(testCase.targetGroup)
		if !reflect.DeepEqual(result, testCase.expectedCustomSD) {
			t.Errorf("%q failed\ngot: %#v\nexpected: %v", testCase.title, result, testCase.expectedCustomSD)
		}
	}
}
