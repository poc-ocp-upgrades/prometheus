package remote

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage"
)

func TestValidateLabelsAndMetricName(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		input		labels.Labels
		expectedErr	string
		shouldPass	bool
	}{{input: labels.FromStrings("__name__", "name", "labelName", "labelValue"), expectedErr: "", shouldPass: true}, {input: labels.FromStrings("__name__", "name", "_labelName", "labelValue"), expectedErr: "", shouldPass: true}, {input: labels.FromStrings("__name__", "name", "@labelName", "labelValue"), expectedErr: "Invalid label name: @labelName", shouldPass: false}, {input: labels.FromStrings("__name__", "name", "123labelName", "labelValue"), expectedErr: "Invalid label name: 123labelName", shouldPass: false}, {input: labels.FromStrings("__name__", "name", "", "labelValue"), expectedErr: "Invalid label name: ", shouldPass: false}, {input: labels.FromStrings("__name__", "name", "labelName", string([]byte{0xff})), expectedErr: "Invalid label value: " + string([]byte{0xff}), shouldPass: false}, {input: labels.FromStrings("__name__", "@invalid_name"), expectedErr: "Invalid metric name: @invalid_name", shouldPass: false}}
	for _, test := range tests {
		err := validateLabelsAndMetricName(test.input)
		if test.shouldPass != (err == nil) {
			if test.shouldPass {
				t.Fatalf("Test should pass, got unexpected error: %v", err)
			} else {
				t.Fatalf("Test should fail, unexpected error, got: %v, expected: %v", err, test.expectedErr)
			}
		}
	}
}
func TestConcreteSeriesSet(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	series1 := &concreteSeries{labels: labels.FromStrings("foo", "bar"), samples: []prompb.Sample{{Value: 1, Timestamp: 2}}}
	series2 := &concreteSeries{labels: labels.FromStrings("foo", "baz"), samples: []prompb.Sample{{Value: 3, Timestamp: 4}}}
	c := &concreteSeriesSet{series: []storage.Series{series1, series2}}
	if !c.Next() {
		t.Fatalf("Expected Next() to be true.")
	}
	if c.At() != series1 {
		t.Fatalf("Unexpected series returned.")
	}
	if !c.Next() {
		t.Fatalf("Expected Next() to be true.")
	}
	if c.At() != series2 {
		t.Fatalf("Unexpected series returned.")
	}
	if c.Next() {
		t.Fatalf("Expected Next() to be false.")
	}
}
func TestConcreteSeriesClonesLabels(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lbls := labels.Labels{labels.Label{Name: "a", Value: "b"}, labels.Label{Name: "c", Value: "d"}}
	cs := concreteSeries{labels: labels.New(lbls...)}
	gotLabels := cs.Labels()
	require.Equal(t, lbls, gotLabels)
	gotLabels[0].Value = "foo"
	gotLabels[1].Value = "bar"
	gotLabels = cs.Labels()
	require.Equal(t, lbls, gotLabels)
}
