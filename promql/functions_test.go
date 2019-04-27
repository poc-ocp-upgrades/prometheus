package promql

import (
	"context"
	"testing"
	"time"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/timestamp"
	"github.com/prometheus/prometheus/util/testutil"
)

func TestDeriv(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	storage := testutil.NewStorage(t)
	defer storage.Close()
	opts := EngineOpts{Logger: nil, Reg: nil, MaxConcurrent: 10, MaxSamples: 10000, Timeout: 10 * time.Second}
	engine := NewEngine(opts)
	a, err := storage.Appender()
	testutil.Ok(t, err)
	metric := labels.FromStrings("__name__", "foo")
	a.Add(metric, 1493712816939, 1.0)
	a.Add(metric, 1493712846939, 1.0)
	err = a.Commit()
	testutil.Ok(t, err)
	query, err := engine.NewInstantQuery(storage, "deriv(foo[30m])", timestamp.Time(1493712846939))
	testutil.Ok(t, err)
	result := query.Exec(context.Background())
	testutil.Ok(t, result.Err)
	vec, _ := result.Vector()
	testutil.Assert(t, len(vec) == 1, "Expected 1 result, got %d", len(vec))
	testutil.Assert(t, vec[0].V == 0.0, "Expected 0.0 as value, got %f", vec[0].V)
}
