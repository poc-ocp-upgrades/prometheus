package remote

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"golang.org/x/time/rate"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	pkgrelabel "github.com/prometheus/prometheus/pkg/relabel"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/relabel"
)

const (
	namespace		= "prometheus"
	subsystem		= "remote_storage"
	queue			= "queue"
	ewmaWeight		= 0.2
	shardUpdateDuration	= 10 * time.Second
	shardToleranceFraction	= 0.3
	logRateLimit		= 0.1
	logBurst		= 10
)

var (
	succeededSamplesTotal	= prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "succeeded_samples_total", Help: "Total number of samples successfully sent to remote storage."}, []string{queue})
	failedSamplesTotal	= prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "failed_samples_total", Help: "Total number of samples which failed on send to remote storage."}, []string{queue})
	droppedSamplesTotal	= prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "dropped_samples_total", Help: "Total number of samples which were dropped due to the queue being full."}, []string{queue})
	sentBatchDuration	= prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem, Name: "sent_batch_duration_seconds", Help: "Duration of sample batch send calls to the remote storage.", Buckets: prometheus.DefBuckets}, []string{queue})
	queueLength		= prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem, Name: "queue_length", Help: "The number of processed samples queued to be sent to the remote storage."}, []string{queue})
	shardCapacity		= prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem, Name: "shard_capacity", Help: "The capacity of each shard of the queue used for parallel sending to the remote storage."}, []string{queue})
	numShards		= prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem, Name: "shards", Help: "The number of shards used for parallel sending to the remote storage."}, []string{queue})
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(succeededSamplesTotal)
	prometheus.MustRegister(failedSamplesTotal)
	prometheus.MustRegister(droppedSamplesTotal)
	prometheus.MustRegister(sentBatchDuration)
	prometheus.MustRegister(queueLength)
	prometheus.MustRegister(shardCapacity)
	prometheus.MustRegister(numShards)
}

type StorageClient interface {
	Store(context.Context, *prompb.WriteRequest) error
	Name() string
}
type QueueManager struct {
	logger						log.Logger
	flushDeadline					time.Duration
	cfg						config.QueueConfig
	externalLabels					model.LabelSet
	relabelConfigs					[]*pkgrelabel.Config
	client						StorageClient
	queueName					string
	logLimiter					*rate.Limiter
	shardsMtx					sync.RWMutex
	shards						*shards
	numShards					int
	reshardChan					chan int
	quit						chan struct{}
	wg						sync.WaitGroup
	samplesIn, samplesOut, samplesOutDuration	*ewmaRate
	integralAccumulator				float64
}

func NewQueueManager(logger log.Logger, cfg config.QueueConfig, externalLabels model.LabelSet, relabelConfigs []*pkgrelabel.Config, client StorageClient, flushDeadline time.Duration) *QueueManager {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if logger == nil {
		logger = log.NewNopLogger()
	} else {
		logger = log.With(logger, "queue", client.Name())
	}
	t := &QueueManager{logger: logger, flushDeadline: flushDeadline, cfg: cfg, externalLabels: externalLabels, relabelConfigs: relabelConfigs, client: client, queueName: client.Name(), logLimiter: rate.NewLimiter(logRateLimit, logBurst), numShards: cfg.MinShards, reshardChan: make(chan int), quit: make(chan struct{}), samplesIn: newEWMARate(ewmaWeight, shardUpdateDuration), samplesOut: newEWMARate(ewmaWeight, shardUpdateDuration), samplesOutDuration: newEWMARate(ewmaWeight, shardUpdateDuration)}
	t.shards = t.newShards(t.numShards)
	numShards.WithLabelValues(t.queueName).Set(float64(t.numShards))
	shardCapacity.WithLabelValues(t.queueName).Set(float64(t.cfg.Capacity))
	sentBatchDuration.WithLabelValues(t.queueName)
	succeededSamplesTotal.WithLabelValues(t.queueName)
	failedSamplesTotal.WithLabelValues(t.queueName)
	droppedSamplesTotal.WithLabelValues(t.queueName)
	return t
}
func (t *QueueManager) Append(s *model.Sample) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	snew := *s
	snew.Metric = s.Metric.Clone()
	for ln, lv := range t.externalLabels {
		if _, ok := s.Metric[ln]; !ok {
			snew.Metric[ln] = lv
		}
	}
	snew.Metric = model.Metric(relabel.Process(model.LabelSet(snew.Metric), t.relabelConfigs...))
	if snew.Metric == nil {
		return nil
	}
	t.shardsMtx.RLock()
	enqueued := t.shards.enqueue(&snew)
	t.shardsMtx.RUnlock()
	if enqueued {
		queueLength.WithLabelValues(t.queueName).Inc()
	} else {
		droppedSamplesTotal.WithLabelValues(t.queueName).Inc()
		if t.logLimiter.Allow() {
			level.Warn(t.logger).Log("msg", "Remote storage queue full, discarding sample. Multiple subsequent messages of this kind may be suppressed.")
		}
	}
	return nil
}
func (*QueueManager) NeedsThrottling() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (t *QueueManager) Start() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.wg.Add(2)
	go t.updateShardsLoop()
	go t.reshardLoop()
	t.shardsMtx.Lock()
	defer t.shardsMtx.Unlock()
	t.shards.start()
}
func (t *QueueManager) Stop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	level.Info(t.logger).Log("msg", "Stopping remote storage...")
	close(t.quit)
	t.wg.Wait()
	t.shardsMtx.Lock()
	defer t.shardsMtx.Unlock()
	t.shards.stop(t.flushDeadline)
	level.Info(t.logger).Log("msg", "Remote storage stopped.")
}
func (t *QueueManager) updateShardsLoop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer t.wg.Done()
	ticker := time.NewTicker(shardUpdateDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			t.calculateDesiredShards()
		case <-t.quit:
			return
		}
	}
}
func (t *QueueManager) calculateDesiredShards() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.samplesIn.tick()
	t.samplesOut.tick()
	t.samplesOutDuration.tick()
	var (
		samplesIn		= t.samplesIn.rate()
		samplesOut		= t.samplesOut.rate()
		samplesPending		= samplesIn - samplesOut
		samplesOutDuration	= t.samplesOutDuration.rate()
	)
	t.integralAccumulator = t.integralAccumulator + (samplesPending * 0.1)
	if samplesOut <= 0 {
		return
	}
	var (
		timePerSample	= samplesOutDuration / samplesOut
		desiredShards	= (timePerSample * (samplesIn + samplesPending + t.integralAccumulator)) / float64(time.Second)
	)
	level.Debug(t.logger).Log("msg", "QueueManager.caclulateDesiredShards", "samplesIn", samplesIn, "samplesOut", samplesOut, "samplesPending", samplesPending, "desiredShards", desiredShards)
	var (
		lowerBound	= float64(t.numShards) * (1. - shardToleranceFraction)
		upperBound	= float64(t.numShards) * (1. + shardToleranceFraction)
	)
	level.Debug(t.logger).Log("msg", "QueueManager.updateShardsLoop", "lowerBound", lowerBound, "desiredShards", desiredShards, "upperBound", upperBound)
	if lowerBound <= desiredShards && desiredShards <= upperBound {
		return
	}
	numShards := int(math.Ceil(desiredShards))
	if numShards > t.cfg.MaxShards {
		numShards = t.cfg.MaxShards
	} else if numShards < t.cfg.MinShards {
		numShards = t.cfg.MinShards
	}
	if numShards == t.numShards {
		return
	}
	select {
	case t.reshardChan <- numShards:
		level.Info(t.logger).Log("msg", "Remote storage resharding", "from", t.numShards, "to", numShards)
		t.numShards = numShards
	default:
		level.Info(t.logger).Log("msg", "Currently resharding, skipping.")
	}
}
func (t *QueueManager) reshardLoop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer t.wg.Done()
	for {
		select {
		case numShards := <-t.reshardChan:
			t.reshard(numShards)
		case <-t.quit:
			return
		}
	}
}
func (t *QueueManager) reshard(n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	numShards.WithLabelValues(t.queueName).Set(float64(n))
	t.shardsMtx.Lock()
	newShards := t.newShards(n)
	oldShards := t.shards
	t.shards = newShards
	t.shardsMtx.Unlock()
	oldShards.stop(t.flushDeadline)
	newShards.start()
}

type shards struct {
	qm	*QueueManager
	queues	[]chan *model.Sample
	done	chan struct{}
	running	int32
	ctx	context.Context
	cancel	context.CancelFunc
}

func (t *QueueManager) newShards(numShards int) *shards {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	queues := make([]chan *model.Sample, numShards)
	for i := 0; i < numShards; i++ {
		queues[i] = make(chan *model.Sample, t.cfg.Capacity)
	}
	ctx, cancel := context.WithCancel(context.Background())
	s := &shards{qm: t, queues: queues, done: make(chan struct{}), running: int32(numShards), ctx: ctx, cancel: cancel}
	return s
}
func (s *shards) start() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := 0; i < len(s.queues); i++ {
		go s.runShard(i)
	}
}
func (s *shards) stop(deadline time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, shard := range s.queues {
		close(shard)
	}
	select {
	case <-s.done:
		return
	case <-time.After(deadline):
		level.Error(s.qm.logger).Log("msg", "Failed to flush all samples on shutdown")
	}
	s.cancel()
	<-s.done
}
func (s *shards) enqueue(sample *model.Sample) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.qm.samplesIn.incr(1)
	fp := sample.Metric.FastFingerprint()
	shard := uint64(fp) % uint64(len(s.queues))
	select {
	case s.queues[shard] <- sample:
		return true
	default:
		return false
	}
}
func (s *shards) runShard(i int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer func() {
		if atomic.AddInt32(&s.running, -1) == 0 {
			close(s.done)
		}
	}()
	queue := s.queues[i]
	pendingSamples := model.Samples{}
	timer := time.NewTimer(time.Duration(s.qm.cfg.BatchSendDeadline))
	stop := func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}
	defer stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case sample, ok := <-queue:
			if !ok {
				if len(pendingSamples) > 0 {
					level.Debug(s.qm.logger).Log("msg", "Flushing samples to remote storage...", "count", len(pendingSamples))
					s.sendSamples(pendingSamples)
					level.Debug(s.qm.logger).Log("msg", "Done flushing.")
				}
				return
			}
			queueLength.WithLabelValues(s.qm.queueName).Dec()
			pendingSamples = append(pendingSamples, sample)
			if len(pendingSamples) >= s.qm.cfg.MaxSamplesPerSend {
				s.sendSamples(pendingSamples[:s.qm.cfg.MaxSamplesPerSend])
				pendingSamples = pendingSamples[s.qm.cfg.MaxSamplesPerSend:]
				stop()
				timer.Reset(time.Duration(s.qm.cfg.BatchSendDeadline))
			}
		case <-timer.C:
			if len(pendingSamples) > 0 {
				s.sendSamples(pendingSamples)
				pendingSamples = pendingSamples[:0]
			}
			timer.Reset(time.Duration(s.qm.cfg.BatchSendDeadline))
		}
	}
}
func (s *shards) sendSamples(samples model.Samples) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	begin := time.Now()
	s.sendSamplesWithBackoff(samples)
	s.qm.samplesOut.incr(int64(len(samples)))
	s.qm.samplesOutDuration.incr(int64(time.Since(begin)))
}
func (s *shards) sendSamplesWithBackoff(samples model.Samples) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	backoff := s.qm.cfg.MinBackoff
	req := ToWriteRequest(samples)
	for retries := s.qm.cfg.MaxRetries; retries > 0; retries-- {
		begin := time.Now()
		err := s.qm.client.Store(s.ctx, req)
		sentBatchDuration.WithLabelValues(s.qm.queueName).Observe(time.Since(begin).Seconds())
		if err == nil {
			succeededSamplesTotal.WithLabelValues(s.qm.queueName).Add(float64(len(samples)))
			return
		}
		level.Warn(s.qm.logger).Log("msg", "Error sending samples to remote storage", "count", len(samples), "err", err)
		if _, ok := err.(recoverableError); !ok {
			break
		}
		time.Sleep(time.Duration(backoff))
		backoff = backoff * 2
		if backoff > s.qm.cfg.MaxBackoff {
			backoff = s.qm.cfg.MaxBackoff
		}
	}
	failedSamplesTotal.WithLabelValues(s.qm.queueName).Add(float64(len(samples)))
}
