package tsdb_test

import (
	"testing"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage/tsdb"
	"github.com/prometheus/prometheus/util/testutil"
)

func TestMetrics(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	db := testutil.NewStorage(t)
	defer db.Close()
	metrics := &dto.Metric{}
	startTime := *tsdb.StartTime
	headMinTime := *tsdb.HeadMinTime
	headMaxTime := *tsdb.HeadMaxTime
	testutil.Ok(t, startTime.Write(metrics))
	testutil.Equals(t, float64(model.Latest)/1000, metrics.Gauge.GetValue())
	testutil.Ok(t, headMinTime.Write(metrics))
	testutil.Equals(t, float64(model.Latest)/1000, metrics.Gauge.GetValue())
	testutil.Ok(t, headMaxTime.Write(metrics))
	testutil.Equals(t, float64(model.Earliest)/1000, metrics.Gauge.GetValue())
	app, err := db.Appender()
	testutil.Ok(t, err)
	app.Add(labels.FromStrings(model.MetricNameLabel, "a"), 1, 1)
	app.Add(labels.FromStrings(model.MetricNameLabel, "a"), 2, 1)
	app.Add(labels.FromStrings(model.MetricNameLabel, "a"), 3, 1)
	testutil.Ok(t, app.Commit())
	testutil.Ok(t, startTime.Write(metrics))
	testutil.Equals(t, 0.001, metrics.Gauge.GetValue())
	testutil.Ok(t, headMinTime.Write(metrics))
	testutil.Equals(t, 0.001, metrics.Gauge.GetValue())
	testutil.Ok(t, headMaxTime.Write(metrics))
	testutil.Equals(t, 0.003, metrics.Gauge.GetValue())
}
