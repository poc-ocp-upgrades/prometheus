package influxdb

import (
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/prometheus/common/model"
)

func TestClient(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	samples := model.Samples{{Metric: model.Metric{model.MetricNameLabel: "testmetric", "test_label": "test_label_value1"}, Timestamp: model.Time(123456789123), Value: 1.23}, {Metric: model.Metric{model.MetricNameLabel: "testmetric", "test_label": "test_label_value2"}, Timestamp: model.Time(123456789123), Value: 5.1234}, {Metric: model.Metric{model.MetricNameLabel: "nan_value"}, Timestamp: model.Time(123456789123), Value: model.SampleValue(math.NaN())}, {Metric: model.Metric{model.MetricNameLabel: "pos_inf_value"}, Timestamp: model.Time(123456789123), Value: model.SampleValue(math.Inf(1))}, {Metric: model.Metric{model.MetricNameLabel: "neg_inf_value"}, Timestamp: model.Time(123456789123), Value: model.SampleValue(math.Inf(-1))}}
	expectedBody := `testmetric,test_label=test_label_value1 value=1.23 123456789123
testmetric,test_label=test_label_value2 value=5.1234 123456789123
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("Unexpected method; expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/write" {
			t.Fatalf("Unexpected path; expected %s, got %s", "/write", r.URL.Path)
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Error reading body: %s", err)
		}
		if string(b) != expectedBody {
			t.Fatalf("Unexpected request body; expected:\n\n%s\n\ngot:\n\n%s", expectedBody, string(b))
		}
	}))
	defer server.Close()
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Unable to parse server URL %s: %s", server.URL, err)
	}
	conf := influx.HTTPConfig{Addr: serverURL.String(), Username: "testuser", Password: "testpass", Timeout: time.Minute}
	c := NewClient(nil, conf, "test_db", "default")
	if err := c.Write(samples); err != nil {
		t.Fatalf("Error sending samples: %s", err)
	}
}
