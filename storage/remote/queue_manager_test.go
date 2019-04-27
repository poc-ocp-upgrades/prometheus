package remote

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/prompb"
)

const defaultFlushDeadline = 1 * time.Minute

type TestStorageClient struct {
	receivedSamples	map[string][]prompb.Sample
	expectedSamples	map[string][]prompb.Sample
	wg		sync.WaitGroup
	mtx		sync.Mutex
}

func NewTestStorageClient() *TestStorageClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TestStorageClient{receivedSamples: map[string][]prompb.Sample{}, expectedSamples: map[string][]prompb.Sample{}}
}
func (c *TestStorageClient) expectSamples(ss model.Samples) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.expectedSamples = map[string][]prompb.Sample{}
	c.receivedSamples = map[string][]prompb.Sample{}
	for _, s := range ss {
		ts := labelProtosToLabels(MetricToLabelProtos(s.Metric)).String()
		c.expectedSamples[ts] = append(c.expectedSamples[ts], prompb.Sample{Timestamp: int64(s.Timestamp), Value: float64(s.Value)})
	}
	c.wg.Add(len(ss))
}
func (c *TestStorageClient) waitForExpectedSamples(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.wg.Wait()
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for ts, expectedSamples := range c.expectedSamples {
		if !reflect.DeepEqual(expectedSamples, c.receivedSamples[ts]) {
			t.Fatalf("%s: Expected %v, got %v", ts, expectedSamples, c.receivedSamples[ts])
		}
	}
}
func (c *TestStorageClient) Store(_ context.Context, req *prompb.WriteRequest) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.mtx.Lock()
	defer c.mtx.Unlock()
	count := 0
	for _, ts := range req.Timeseries {
		labels := labelProtosToLabels(ts.Labels).String()
		for _, sample := range ts.Samples {
			count++
			c.receivedSamples[labels] = append(c.receivedSamples[labels], sample)
		}
	}
	c.wg.Add(-count)
	return nil
}
func (c *TestStorageClient) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "teststorageclient"
}
func TestSampleDelivery(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n := config.DefaultQueueConfig.Capacity * 2
	samples := make(model.Samples, 0, n)
	for i := 0; i < n; i++ {
		name := model.LabelValue(fmt.Sprintf("test_metric_%d", i))
		samples = append(samples, &model.Sample{Metric: model.Metric{model.MetricNameLabel: name}, Value: model.SampleValue(i)})
	}
	c := NewTestStorageClient()
	c.expectSamples(samples[:len(samples)/2])
	cfg := config.DefaultQueueConfig
	cfg.MaxShards = 1
	m := NewQueueManager(nil, cfg, nil, nil, c, defaultFlushDeadline)
	for _, s := range samples[:len(samples)/2] {
		m.Append(s)
	}
	for _, s := range samples[len(samples)/2:] {
		m.Append(s)
	}
	m.Start()
	defer m.Stop()
	c.waitForExpectedSamples(t)
}
func TestSampleDeliveryTimeout(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n := config.DefaultQueueConfig.Capacity - 1
	samples := make(model.Samples, 0, n)
	for i := 0; i < n; i++ {
		name := model.LabelValue(fmt.Sprintf("test_metric_%d", i))
		samples = append(samples, &model.Sample{Metric: model.Metric{model.MetricNameLabel: name}, Value: model.SampleValue(i)})
	}
	c := NewTestStorageClient()
	cfg := config.DefaultQueueConfig
	cfg.MaxShards = 1
	cfg.BatchSendDeadline = model.Duration(100 * time.Millisecond)
	m := NewQueueManager(nil, cfg, nil, nil, c, defaultFlushDeadline)
	m.Start()
	defer m.Stop()
	c.expectSamples(samples)
	for _, s := range samples {
		m.Append(s)
	}
	c.waitForExpectedSamples(t)
	c.expectSamples(samples)
	for _, s := range samples {
		m.Append(s)
	}
	c.waitForExpectedSamples(t)
}
func TestSampleDeliveryOrder(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ts := 10
	n := config.DefaultQueueConfig.MaxSamplesPerSend * ts
	samples := make(model.Samples, 0, n)
	for i := 0; i < n; i++ {
		name := model.LabelValue(fmt.Sprintf("test_metric_%d", i%ts))
		samples = append(samples, &model.Sample{Metric: model.Metric{model.MetricNameLabel: name}, Value: model.SampleValue(i), Timestamp: model.Time(i)})
	}
	c := NewTestStorageClient()
	c.expectSamples(samples)
	m := NewQueueManager(nil, config.DefaultQueueConfig, nil, nil, c, defaultFlushDeadline)
	for _, s := range samples {
		m.Append(s)
	}
	m.Start()
	defer m.Stop()
	c.waitForExpectedSamples(t)
}

type TestBlockingStorageClient struct {
	numCalls	uint64
	block		chan bool
}

func NewTestBlockedStorageClient() *TestBlockingStorageClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TestBlockingStorageClient{block: make(chan bool), numCalls: 0}
}
func (c *TestBlockingStorageClient) Store(ctx context.Context, _ *prompb.WriteRequest) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	atomic.AddUint64(&c.numCalls, 1)
	select {
	case <-c.block:
	case <-ctx.Done():
	}
	return nil
}
func (c *TestBlockingStorageClient) NumCalls() uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return atomic.LoadUint64(&c.numCalls)
}
func (c *TestBlockingStorageClient) unlock() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	close(c.block)
}
func (c *TestBlockingStorageClient) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "testblockingstorageclient"
}
func (t *QueueManager) queueLen() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.shardsMtx.Lock()
	defer t.shardsMtx.Unlock()
	queueLength := 0
	for _, shard := range t.shards.queues {
		queueLength += len(shard)
	}
	return queueLength
}
func TestSpawnNotMoreThanMaxConcurrentSendsGoroutines(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n := config.DefaultQueueConfig.MaxSamplesPerSend * 2
	samples := make(model.Samples, 0, n)
	for i := 0; i < n; i++ {
		name := model.LabelValue(fmt.Sprintf("test_metric_%d", i))
		samples = append(samples, &model.Sample{Metric: model.Metric{model.MetricNameLabel: name}, Value: model.SampleValue(i)})
	}
	c := NewTestBlockedStorageClient()
	cfg := config.DefaultQueueConfig
	cfg.MaxShards = 1
	cfg.Capacity = n
	m := NewQueueManager(nil, cfg, nil, nil, c, defaultFlushDeadline)
	m.Start()
	defer func() {
		c.unlock()
		m.Stop()
	}()
	for _, s := range samples {
		m.Append(s)
	}
	for i := 0; i < 100 && m.queueLen() > 0; i++ {
		time.Sleep(10 * time.Millisecond)
	}
	if m.queueLen() != config.DefaultQueueConfig.MaxSamplesPerSend {
		t.Fatalf("Failed to drain QueueManager queue, %d elements left", m.queueLen())
	}
	numCalls := c.NumCalls()
	if numCalls != uint64(1) {
		t.Errorf("Saw %d concurrent sends, expected 1", numCalls)
	}
}
func TestShutdown(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	deadline := 10 * time.Second
	c := NewTestBlockedStorageClient()
	m := NewQueueManager(nil, config.DefaultQueueConfig, nil, nil, c, deadline)
	for i := 0; i < config.DefaultQueueConfig.MaxSamplesPerSend; i++ {
		m.Append(&model.Sample{Metric: model.Metric{model.MetricNameLabel: model.LabelValue(fmt.Sprintf("test_metric_%d", i))}, Value: model.SampleValue(i), Timestamp: model.Time(i)})
	}
	m.Start()
	start := time.Now()
	m.Stop()
	duration := time.Since(start)
	if duration > deadline+(deadline/10) {
		t.Errorf("Took too long to shutdown: %s > %s", duration, deadline)
	}
}
