package graphite

import (
	"testing"
	"github.com/prometheus/common/model"
)

var (
	metric = model.Metric{model.MetricNameLabel: "test:metric", "testlabel": "test:value", "many_chars": "abc!ABC:012-3!45ö67~89./(){},=.\"\\"}
)

func TestEscape(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	value := "abzABZ019(){},'\"\\"
	expected := "abzABZ019\\(\\)\\{\\}\\,\\'\\\"\\\\"
	actual := escape(model.LabelValue(value))
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
	value = "é/|_;:%."
	expected = "%C3%A9%2F|_;:%25%2E"
	actual = escape(model.LabelValue(value))
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
func TestPathFromMetric(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	expected := ("prefix." + "test:metric" + ".many_chars.abc!ABC:012-3!45%C3%B667~89%2E%2F\\(\\)\\{\\}\\,%3D%2E\\\"\\\\" + ".testlabel.test:value")
	actual := pathFromMetric(metric, "prefix.")
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
