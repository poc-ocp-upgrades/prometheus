package scrape

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"time"
	"unsafe"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/common/version"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/pool"
	"github.com/prometheus/prometheus/pkg/relabel"
	"github.com/prometheus/prometheus/pkg/textparse"
	"github.com/prometheus/prometheus/pkg/timestamp"
	"github.com/prometheus/prometheus/pkg/value"
	"github.com/prometheus/prometheus/storage"
)

var (
	targetIntervalLength			= prometheus.NewSummaryVec(prometheus.SummaryOpts{Name: "prometheus_target_interval_length_seconds", Help: "Actual intervals between scrapes.", Objectives: map[float64]float64{0.01: 0.001, 0.05: 0.005, 0.5: 0.05, 0.90: 0.01, 0.99: 0.001}}, []string{"interval"})
	targetReloadIntervalLength		= prometheus.NewSummaryVec(prometheus.SummaryOpts{Name: "prometheus_target_reload_length_seconds", Help: "Actual interval to reload the scrape pool with a given configuration.", Objectives: map[float64]float64{0.01: 0.001, 0.05: 0.005, 0.5: 0.05, 0.90: 0.01, 0.99: 0.001}}, []string{"interval"})
	targetSyncIntervalLength		= prometheus.NewSummaryVec(prometheus.SummaryOpts{Name: "prometheus_target_sync_length_seconds", Help: "Actual interval to sync the scrape pool.", Objectives: map[float64]float64{0.01: 0.001, 0.05: 0.005, 0.5: 0.05, 0.90: 0.01, 0.99: 0.001}}, []string{"scrape_job"})
	targetScrapePoolSyncsCounter	= prometheus.NewCounterVec(prometheus.CounterOpts{Name: "prometheus_target_scrape_pool_sync_total", Help: "Total number of syncs that were executed on a scrape pool."}, []string{"scrape_job"})
	targetScrapeSampleLimit			= prometheus.NewCounter(prometheus.CounterOpts{Name: "prometheus_target_scrapes_exceeded_sample_limit_total", Help: "Total number of scrapes that hit the sample limit and were rejected."})
	targetScrapeSampleDuplicate		= prometheus.NewCounter(prometheus.CounterOpts{Name: "prometheus_target_scrapes_sample_duplicate_timestamp_total", Help: "Total number of samples rejected due to duplicate timestamps but different values"})
	targetScrapeSampleOutOfOrder	= prometheus.NewCounter(prometheus.CounterOpts{Name: "prometheus_target_scrapes_sample_out_of_order_total", Help: "Total number of samples rejected due to not being out of the expected order"})
	targetScrapeSampleOutOfBounds	= prometheus.NewCounter(prometheus.CounterOpts{Name: "prometheus_target_scrapes_sample_out_of_bounds_total", Help: "Total number of samples rejected due to timestamp falling outside of the time bounds"})
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(targetIntervalLength)
	prometheus.MustRegister(targetReloadIntervalLength)
	prometheus.MustRegister(targetSyncIntervalLength)
	prometheus.MustRegister(targetScrapePoolSyncsCounter)
	prometheus.MustRegister(targetScrapeSampleLimit)
	prometheus.MustRegister(targetScrapeSampleDuplicate)
	prometheus.MustRegister(targetScrapeSampleOutOfOrder)
	prometheus.MustRegister(targetScrapeSampleOutOfBounds)
}

type scrapePool struct {
	appendable		Appendable
	logger			log.Logger
	mtx				sync.RWMutex
	config			*config.ScrapeConfig
	client			*http.Client
	activeTargets	map[uint64]*Target
	droppedTargets	[]*Target
	loops			map[uint64]loop
	cancel			context.CancelFunc
	newLoop			func(*Target, scraper, int, bool, []*relabel.Config) loop
}

const maxAheadTime = 10 * time.Minute

type labelsMutator func(labels.Labels) labels.Labels

func newScrapePool(cfg *config.ScrapeConfig, app Appendable, logger log.Logger) *scrapePool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if logger == nil {
		logger = log.NewNopLogger()
	}
	client, err := config_util.NewClientFromConfig(cfg.HTTPClientConfig, cfg.JobName)
	if err != nil {
		level.Error(logger).Log("msg", "Error creating HTTP client", "err", err)
	}
	buffers := pool.New(1e3, 100e6, 3, func(sz int) interface{} {
		return make([]byte, 0, sz)
	})
	ctx, cancel := context.WithCancel(context.Background())
	sp := &scrapePool{cancel: cancel, appendable: app, config: cfg, client: client, activeTargets: map[uint64]*Target{}, loops: map[uint64]loop{}, logger: logger}
	sp.newLoop = func(t *Target, s scraper, limit int, honor bool, mrc []*relabel.Config) loop {
		cache := newScrapeCache()
		t.setMetadataStore(cache)
		return newScrapeLoop(ctx, s, log.With(logger, "target", t), buffers, func(l labels.Labels) labels.Labels {
			return mutateSampleLabels(l, t, honor, mrc)
		}, func(l labels.Labels) labels.Labels {
			return mutateReportSampleLabels(l, t)
		}, func() storage.Appender {
			app, err := app.Appender()
			if err != nil {
				panic(err)
			}
			return appender(app, limit)
		}, cache)
	}
	return sp
}
func (sp *scrapePool) ActiveTargets() []*Target {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sp.mtx.Lock()
	defer sp.mtx.Unlock()
	var tActive []*Target
	for _, t := range sp.activeTargets {
		tActive = append(tActive, t)
	}
	return tActive
}
func (sp *scrapePool) DroppedTargets() []*Target {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sp.mtx.Lock()
	defer sp.mtx.Unlock()
	return sp.droppedTargets
}
func (sp *scrapePool) stop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sp.cancel()
	var wg sync.WaitGroup
	sp.mtx.Lock()
	defer sp.mtx.Unlock()
	for fp, l := range sp.loops {
		wg.Add(1)
		go func(l loop) {
			l.stop()
			wg.Done()
		}(l)
		delete(sp.loops, fp)
		delete(sp.activeTargets, fp)
	}
	wg.Wait()
}
func (sp *scrapePool) reload(cfg *config.ScrapeConfig) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	start := time.Now()
	sp.mtx.Lock()
	defer sp.mtx.Unlock()
	client, err := config_util.NewClientFromConfig(cfg.HTTPClientConfig, cfg.JobName)
	if err != nil {
		level.Error(sp.logger).Log("msg", "Error creating HTTP client", "err", err)
	}
	sp.config = cfg
	sp.client = client
	var (
		wg			sync.WaitGroup
		interval	= time.Duration(sp.config.ScrapeInterval)
		timeout		= time.Duration(sp.config.ScrapeTimeout)
		limit		= int(sp.config.SampleLimit)
		honor		= sp.config.HonorLabels
		mrc			= sp.config.MetricRelabelConfigs
	)
	for fp, oldLoop := range sp.loops {
		var (
			t		= sp.activeTargets[fp]
			s		= &targetScraper{Target: t, client: sp.client, timeout: timeout}
			newLoop	= sp.newLoop(t, s, limit, honor, mrc)
		)
		wg.Add(1)
		go func(oldLoop, newLoop loop) {
			oldLoop.stop()
			wg.Done()
			go newLoop.run(interval, timeout, nil)
		}(oldLoop, newLoop)
		sp.loops[fp] = newLoop
	}
	wg.Wait()
	targetReloadIntervalLength.WithLabelValues(interval.String()).Observe(time.Since(start).Seconds())
}
func (sp *scrapePool) Sync(tgs []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	start := time.Now()
	var all []*Target
	sp.mtx.Lock()
	sp.droppedTargets = []*Target{}
	for _, tg := range tgs {
		targets, err := targetsFromGroup(tg, sp.config)
		if err != nil {
			level.Error(sp.logger).Log("msg", "creating targets failed", "err", err)
			continue
		}
		for _, t := range targets {
			if t.Labels().Len() > 0 {
				all = append(all, t)
			} else if t.DiscoveredLabels().Len() > 0 {
				sp.droppedTargets = append(sp.droppedTargets, t)
			}
		}
	}
	sp.mtx.Unlock()
	sp.sync(all)
	targetSyncIntervalLength.WithLabelValues(sp.config.JobName).Observe(time.Since(start).Seconds())
	targetScrapePoolSyncsCounter.WithLabelValues(sp.config.JobName).Inc()
}
func (sp *scrapePool) sync(targets []*Target) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sp.mtx.Lock()
	defer sp.mtx.Unlock()
	var (
		uniqueTargets	= map[uint64]struct{}{}
		interval		= time.Duration(sp.config.ScrapeInterval)
		timeout			= time.Duration(sp.config.ScrapeTimeout)
		limit			= int(sp.config.SampleLimit)
		honor			= sp.config.HonorLabels
		mrc				= sp.config.MetricRelabelConfigs
	)
	for _, t := range targets {
		t := t
		hash := t.hash()
		uniqueTargets[hash] = struct{}{}
		if _, ok := sp.activeTargets[hash]; !ok {
			s := &targetScraper{Target: t, client: sp.client, timeout: timeout}
			l := sp.newLoop(t, s, limit, honor, mrc)
			sp.activeTargets[hash] = t
			sp.loops[hash] = l
			go l.run(interval, timeout, nil)
		} else {
			sp.activeTargets[hash].SetDiscoveredLabels(t.DiscoveredLabels())
		}
	}
	var wg sync.WaitGroup
	for hash := range sp.activeTargets {
		if _, ok := uniqueTargets[hash]; !ok {
			wg.Add(1)
			go func(l loop) {
				l.stop()
				wg.Done()
			}(sp.loops[hash])
			delete(sp.loops, hash)
			delete(sp.activeTargets, hash)
		}
	}
	wg.Wait()
}
func mutateSampleLabels(lset labels.Labels, target *Target, honor bool, rc []*relabel.Config) labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lb := labels.NewBuilder(lset)
	if honor {
		for _, l := range target.Labels() {
			if !lset.Has(l.Name) {
				lb.Set(l.Name, l.Value)
			}
		}
	} else {
		for _, l := range target.Labels() {
			lv := lset.Get(l.Name)
			if lv != "" {
				lb.Set(model.ExportedLabelPrefix+l.Name, lv)
			}
			lb.Set(l.Name, l.Value)
		}
	}
	for _, l := range lb.Labels() {
		if l.Value == "" {
			lb.Del(l.Name)
		}
	}
	res := lb.Labels()
	if len(rc) > 0 {
		res = relabel.Process(res, rc...)
	}
	return res
}
func mutateReportSampleLabels(lset labels.Labels, target *Target) labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lb := labels.NewBuilder(lset)
	for _, l := range target.Labels() {
		lv := lset.Get(l.Name)
		if lv != "" {
			lb.Set(model.ExportedLabelPrefix+l.Name, lv)
		}
		lb.Set(l.Name, l.Value)
	}
	return lb.Labels()
}
func appender(app storage.Appender, limit int) storage.Appender {
	_logClusterCodePath()
	defer _logClusterCodePath()
	app = &timeLimitAppender{Appender: app, maxTime: timestamp.FromTime(time.Now().Add(maxAheadTime))}
	if limit > 0 {
		app = &limitAppender{Appender: app, limit: limit}
	}
	return app
}

type scraper interface {
	scrape(ctx context.Context, w io.Writer) (string, error)
	report(start time.Time, dur time.Duration, err error)
	offset(interval time.Duration) time.Duration
}
type targetScraper struct {
	*Target
	client	*http.Client
	req		*http.Request
	timeout	time.Duration
	gzipr	*gzip.Reader
	buf		*bufio.Reader
}

const acceptHeader = `application/openmetrics-text; version=0.0.1,text/plain;version=0.0.4;q=0.5,*/*;q=0.1`

var userAgentHeader = fmt.Sprintf("Prometheus/%s", version.Version)

func (s *targetScraper) scrape(ctx context.Context, w io.Writer) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.req == nil {
		req, err := http.NewRequest("GET", s.URL().String(), nil)
		if err != nil {
			return "", err
		}
		req.Header.Add("Accept", acceptHeader)
		req.Header.Add("Accept-Encoding", "gzip")
		req.Header.Set("User-Agent", userAgentHeader)
		req.Header.Set("X-Prometheus-Scrape-Timeout-Seconds", fmt.Sprintf("%f", s.timeout.Seconds()))
		s.req = req
	}
	resp, err := s.client.Do(s.req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned HTTP status %s", resp.Status)
	}
	if resp.Header.Get("Content-Encoding") != "gzip" {
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			return "", err
		}
		return resp.Header.Get("Content-Type"), nil
	}
	if s.gzipr == nil {
		s.buf = bufio.NewReader(resp.Body)
		s.gzipr, err = gzip.NewReader(s.buf)
		if err != nil {
			return "", err
		}
	} else {
		s.buf.Reset(resp.Body)
		if err = s.gzipr.Reset(s.buf); err != nil {
			return "", err
		}
	}
	_, err = io.Copy(w, s.gzipr)
	s.gzipr.Close()
	if err != nil {
		return "", err
	}
	return resp.Header.Get("Content-Type"), nil
}

type loop interface {
	run(interval, timeout time.Duration, errc chan<- error)
	stop()
}
type cacheEntry struct {
	ref			uint64
	lastIter	uint64
	hash		uint64
	lset		labels.Labels
}
type scrapeLoop struct {
	scraper				scraper
	l					log.Logger
	cache				*scrapeCache
	lastScrapeSize		int
	buffers				*pool.Pool
	appender			func() storage.Appender
	sampleMutator		labelsMutator
	reportSampleMutator	labelsMutator
	ctx					context.Context
	scrapeCtx			context.Context
	cancel				func()
	stopped				chan struct{}
}
type scrapeCache struct {
	iter			uint64
	series			map[string]*cacheEntry
	droppedSeries	map[string]*uint64
	seriesCur		map[uint64]labels.Labels
	seriesPrev		map[uint64]labels.Labels
	metaMtx			sync.Mutex
	metadata		map[string]*metaEntry
}
type metaEntry struct {
	lastIter	uint64
	typ			textparse.MetricType
	help		string
	unit		string
}

func newScrapeCache() *scrapeCache {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &scrapeCache{series: map[string]*cacheEntry{}, droppedSeries: map[string]*uint64{}, seriesCur: map[uint64]labels.Labels{}, seriesPrev: map[uint64]labels.Labels{}, metadata: map[string]*metaEntry{}}
}
func (c *scrapeCache) iterDone() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for s, e := range c.series {
		if c.iter-e.lastIter > 2 {
			delete(c.series, s)
		}
	}
	for s, iter := range c.droppedSeries {
		if c.iter-*iter > 2 {
			delete(c.droppedSeries, s)
		}
	}
	c.metaMtx.Lock()
	for m, e := range c.metadata {
		if c.iter-e.lastIter > 10 {
			delete(c.metadata, m)
		}
	}
	c.metaMtx.Unlock()
	c.seriesPrev, c.seriesCur = c.seriesCur, c.seriesPrev
	for k := range c.seriesCur {
		delete(c.seriesCur, k)
	}
	c.iter++
}
func (c *scrapeCache) get(met string) (*cacheEntry, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	e, ok := c.series[met]
	if !ok {
		return nil, false
	}
	e.lastIter = c.iter
	return e, true
}
func (c *scrapeCache) addRef(met string, ref uint64, lset labels.Labels, hash uint64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if ref == 0 {
		return
	}
	c.series[met] = &cacheEntry{ref: ref, lastIter: c.iter, lset: lset, hash: hash}
}
func (c *scrapeCache) addDropped(met string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	iter := c.iter
	c.droppedSeries[met] = &iter
}
func (c *scrapeCache) getDropped(met string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	iterp, ok := c.droppedSeries[met]
	if ok {
		*iterp = c.iter
	}
	return ok
}
func (c *scrapeCache) trackStaleness(hash uint64, lset labels.Labels) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.seriesCur[hash] = lset
}
func (c *scrapeCache) forEachStale(f func(labels.Labels) bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for h, lset := range c.seriesPrev {
		if _, ok := c.seriesCur[h]; !ok {
			if !f(lset) {
				break
			}
		}
	}
}
func (c *scrapeCache) setType(metric []byte, t textparse.MetricType) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.metaMtx.Lock()
	e, ok := c.metadata[yoloString(metric)]
	if !ok {
		e = &metaEntry{typ: textparse.MetricTypeUnknown}
		c.metadata[string(metric)] = e
	}
	e.typ = t
	e.lastIter = c.iter
	c.metaMtx.Unlock()
}
func (c *scrapeCache) setHelp(metric, help []byte) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.metaMtx.Lock()
	e, ok := c.metadata[yoloString(metric)]
	if !ok {
		e = &metaEntry{typ: textparse.MetricTypeUnknown}
		c.metadata[string(metric)] = e
	}
	if e.help != yoloString(help) {
		e.help = string(help)
	}
	e.lastIter = c.iter
	c.metaMtx.Unlock()
}
func (c *scrapeCache) setUnit(metric, unit []byte) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.metaMtx.Lock()
	e, ok := c.metadata[yoloString(metric)]
	if !ok {
		e = &metaEntry{typ: textparse.MetricTypeUnknown}
		c.metadata[string(metric)] = e
	}
	if e.unit != yoloString(unit) {
		e.unit = string(unit)
	}
	e.lastIter = c.iter
	c.metaMtx.Unlock()
}
func (c *scrapeCache) getMetadata(metric string) (MetricMetadata, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.metaMtx.Lock()
	defer c.metaMtx.Unlock()
	m, ok := c.metadata[metric]
	if !ok {
		return MetricMetadata{}, false
	}
	return MetricMetadata{Metric: metric, Type: m.typ, Help: m.help, Unit: m.unit}, true
}
func (c *scrapeCache) listMetadata() []MetricMetadata {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.metaMtx.Lock()
	defer c.metaMtx.Unlock()
	res := make([]MetricMetadata, 0, len(c.metadata))
	for m, e := range c.metadata {
		res = append(res, MetricMetadata{Metric: m, Type: e.typ, Help: e.help, Unit: e.unit})
	}
	return res
}
func newScrapeLoop(ctx context.Context, sc scraper, l log.Logger, buffers *pool.Pool, sampleMutator labelsMutator, reportSampleMutator labelsMutator, appender func() storage.Appender, cache *scrapeCache) *scrapeLoop {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if l == nil {
		l = log.NewNopLogger()
	}
	if buffers == nil {
		buffers = pool.New(1e3, 1e6, 3, func(sz int) interface{} {
			return make([]byte, 0, sz)
		})
	}
	if cache == nil {
		cache = newScrapeCache()
	}
	sl := &scrapeLoop{scraper: sc, buffers: buffers, cache: cache, appender: appender, sampleMutator: sampleMutator, reportSampleMutator: reportSampleMutator, stopped: make(chan struct{}), l: l, ctx: ctx}
	sl.scrapeCtx, sl.cancel = context.WithCancel(ctx)
	return sl
}
func (sl *scrapeLoop) run(interval, timeout time.Duration, errc chan<- error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	select {
	case <-time.After(sl.scraper.offset(interval)):
	case <-sl.scrapeCtx.Done():
		close(sl.stopped)
		return
	}
	var last time.Time
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
mainLoop:
	for {
		select {
		case <-sl.ctx.Done():
			close(sl.stopped)
			return
		case <-sl.scrapeCtx.Done():
			break mainLoop
		default:
		}
		var (
			start				= time.Now()
			scrapeCtx, cancel	= context.WithTimeout(sl.ctx, timeout)
		)
		if !last.IsZero() {
			targetIntervalLength.WithLabelValues(interval.String()).Observe(time.Since(last).Seconds())
		}
		b := sl.buffers.Get(sl.lastScrapeSize).([]byte)
		buf := bytes.NewBuffer(b)
		contentType, scrapeErr := sl.scraper.scrape(scrapeCtx, buf)
		cancel()
		if scrapeErr == nil {
			b = buf.Bytes()
			if len(b) > 0 {
				sl.lastScrapeSize = len(b)
			}
		} else {
			level.Debug(sl.l).Log("msg", "Scrape failed", "err", scrapeErr.Error())
			if errc != nil {
				errc <- scrapeErr
			}
		}
		total, added, appErr := sl.append(b, contentType, start)
		if appErr != nil {
			level.Warn(sl.l).Log("msg", "append failed", "err", appErr)
			if _, _, err := sl.append([]byte{}, "", start); err != nil {
				level.Warn(sl.l).Log("msg", "append failed", "err", err)
			}
		}
		sl.buffers.Put(b)
		if scrapeErr == nil {
			scrapeErr = appErr
		}
		if err := sl.report(start, time.Since(start), total, added, scrapeErr); err != nil {
			level.Warn(sl.l).Log("msg", "appending scrape report failed", "err", err)
		}
		last = start
		select {
		case <-sl.ctx.Done():
			close(sl.stopped)
			return
		case <-sl.scrapeCtx.Done():
			break mainLoop
		case <-ticker.C:
		}
	}
	close(sl.stopped)
	sl.endOfRunStaleness(last, ticker, interval)
}
func (sl *scrapeLoop) endOfRunStaleness(last time.Time, ticker *time.Ticker, interval time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if last.IsZero() {
		return
	}
	var staleTime time.Time
	select {
	case <-sl.ctx.Done():
		return
	case <-ticker.C:
		staleTime = time.Now()
	}
	select {
	case <-sl.ctx.Done():
		return
	case <-ticker.C:
	}
	select {
	case <-sl.ctx.Done():
		return
	case <-time.After(interval / 10):
	}
	if _, _, err := sl.append([]byte{}, "", staleTime); err != nil {
		level.Error(sl.l).Log("msg", "stale append failed", "err", err)
	}
	if err := sl.reportStale(staleTime); err != nil {
		level.Error(sl.l).Log("msg", "stale report failed", "err", err)
	}
}
func (sl *scrapeLoop) stop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sl.cancel()
	<-sl.stopped
}

type sample struct {
	metric	labels.Labels
	t		int64
	v		float64
}
type samples []sample

func (s samples) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(s)
}
func (s samples) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s[i], s[j] = s[j], s[i]
}
func (s samples) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d := labels.Compare(s[i].metric, s[j].metric)
	if d < 0 {
		return true
	} else if d > 0 {
		return false
	}
	return s[i].t < s[j].t
}
func (sl *scrapeLoop) append(b []byte, contentType string, ts time.Time) (total, added int, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		app				= sl.appender()
		p				= textparse.New(b, contentType)
		defTime			= timestamp.FromTime(ts)
		numOutOfOrder	= 0
		numDuplicates	= 0
		numOutOfBounds	= 0
	)
	var sampleLimitErr error
loop:
	for {
		var et textparse.Entry
		if et, err = p.Next(); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		switch et {
		case textparse.EntryType:
			sl.cache.setType(p.Type())
			continue
		case textparse.EntryHelp:
			sl.cache.setHelp(p.Help())
			continue
		case textparse.EntryUnit:
			sl.cache.setUnit(p.Unit())
			continue
		case textparse.EntryComment:
			continue
		default:
		}
		total++
		t := defTime
		met, tp, v := p.Series()
		if tp != nil {
			t = *tp
		}
		if sl.cache.getDropped(yoloString(met)) {
			continue
		}
		ce, ok := sl.cache.get(yoloString(met))
		if ok {
			switch err = app.AddFast(ce.lset, ce.ref, t, v); err {
			case nil:
				if tp == nil {
					sl.cache.trackStaleness(ce.hash, ce.lset)
				}
			case storage.ErrNotFound:
				ok = false
			case storage.ErrOutOfOrderSample:
				numOutOfOrder++
				level.Debug(sl.l).Log("msg", "Out of order sample", "series", string(met))
				targetScrapeSampleOutOfOrder.Inc()
				continue
			case storage.ErrDuplicateSampleForTimestamp:
				numDuplicates++
				level.Debug(sl.l).Log("msg", "Duplicate sample for timestamp", "series", string(met))
				targetScrapeSampleDuplicate.Inc()
				continue
			case storage.ErrOutOfBounds:
				numOutOfBounds++
				level.Debug(sl.l).Log("msg", "Out of bounds metric", "series", string(met))
				targetScrapeSampleOutOfBounds.Inc()
				continue
			case errSampleLimit:
				sampleLimitErr = err
				added++
				continue
			default:
				break loop
			}
		}
		if !ok {
			var lset labels.Labels
			mets := p.Metric(&lset)
			hash := lset.Hash()
			lset = sl.sampleMutator(lset)
			if lset == nil {
				sl.cache.addDropped(mets)
				continue
			}
			var ref uint64
			ref, err = app.Add(lset, t, v)
			switch err {
			case nil:
			case storage.ErrOutOfOrderSample:
				err = nil
				numOutOfOrder++
				level.Debug(sl.l).Log("msg", "Out of order sample", "series", string(met))
				targetScrapeSampleOutOfOrder.Inc()
				continue
			case storage.ErrDuplicateSampleForTimestamp:
				err = nil
				numDuplicates++
				level.Debug(sl.l).Log("msg", "Duplicate sample for timestamp", "series", string(met))
				targetScrapeSampleDuplicate.Inc()
				continue
			case storage.ErrOutOfBounds:
				err = nil
				numOutOfBounds++
				level.Debug(sl.l).Log("msg", "Out of bounds metric", "series", string(met))
				targetScrapeSampleOutOfBounds.Inc()
				continue
			case errSampleLimit:
				sampleLimitErr = err
				added++
				continue
			default:
				level.Debug(sl.l).Log("msg", "unexpected error", "series", string(met), "err", err)
				break loop
			}
			if tp == nil {
				sl.cache.trackStaleness(hash, lset)
			}
			sl.cache.addRef(mets, ref, lset, hash)
		}
		added++
	}
	if sampleLimitErr != nil {
		if err == nil {
			err = sampleLimitErr
		}
		targetScrapeSampleLimit.Inc()
	}
	if numOutOfOrder > 0 {
		level.Warn(sl.l).Log("msg", "Error on ingesting out-of-order samples", "num_dropped", numOutOfOrder)
	}
	if numDuplicates > 0 {
		level.Warn(sl.l).Log("msg", "Error on ingesting samples with different value but same timestamp", "num_dropped", numDuplicates)
	}
	if numOutOfBounds > 0 {
		level.Warn(sl.l).Log("msg", "Error on ingesting samples that are too old or are too far into the future", "num_dropped", numOutOfBounds)
	}
	if err == nil {
		sl.cache.forEachStale(func(lset labels.Labels) bool {
			_, err = app.Add(lset, defTime, math.Float64frombits(value.StaleNaN))
			switch err {
			case storage.ErrOutOfOrderSample, storage.ErrDuplicateSampleForTimestamp:
				err = nil
			}
			return err == nil
		})
	}
	if err != nil {
		app.Rollback()
		return total, added, err
	}
	if err := app.Commit(); err != nil {
		return total, added, err
	}
	sl.cache.iterDone()
	return total, added, nil
}
func yoloString(b []byte) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return *((*string)(unsafe.Pointer(&b)))
}

const (
	scrapeHealthMetricName			= "up" + "\xff"
	scrapeDurationMetricName		= "scrape_duration_seconds" + "\xff"
	scrapeSamplesMetricName			= "scrape_samples_scraped" + "\xff"
	samplesPostRelabelMetricName	= "scrape_samples_post_metric_relabeling" + "\xff"
)

func (sl *scrapeLoop) report(start time.Time, duration time.Duration, scraped, appended int, err error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sl.scraper.report(start, duration, err)
	ts := timestamp.FromTime(start)
	var health float64
	if err == nil {
		health = 1
	}
	app := sl.appender()
	if err := sl.addReportSample(app, scrapeHealthMetricName, ts, health); err != nil {
		app.Rollback()
		return err
	}
	if err := sl.addReportSample(app, scrapeDurationMetricName, ts, duration.Seconds()); err != nil {
		app.Rollback()
		return err
	}
	if err := sl.addReportSample(app, scrapeSamplesMetricName, ts, float64(scraped)); err != nil {
		app.Rollback()
		return err
	}
	if err := sl.addReportSample(app, samplesPostRelabelMetricName, ts, float64(appended)); err != nil {
		app.Rollback()
		return err
	}
	return app.Commit()
}
func (sl *scrapeLoop) reportStale(start time.Time) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ts := timestamp.FromTime(start)
	app := sl.appender()
	stale := math.Float64frombits(value.StaleNaN)
	if err := sl.addReportSample(app, scrapeHealthMetricName, ts, stale); err != nil {
		app.Rollback()
		return err
	}
	if err := sl.addReportSample(app, scrapeDurationMetricName, ts, stale); err != nil {
		app.Rollback()
		return err
	}
	if err := sl.addReportSample(app, scrapeSamplesMetricName, ts, stale); err != nil {
		app.Rollback()
		return err
	}
	if err := sl.addReportSample(app, samplesPostRelabelMetricName, ts, stale); err != nil {
		app.Rollback()
		return err
	}
	return app.Commit()
}
func (sl *scrapeLoop) addReportSample(app storage.Appender, s string, t int64, v float64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ce, ok := sl.cache.get(s)
	if ok {
		err := app.AddFast(ce.lset, ce.ref, t, v)
		switch err {
		case nil:
			return nil
		case storage.ErrNotFound:
		case storage.ErrOutOfOrderSample, storage.ErrDuplicateSampleForTimestamp:
			return nil
		default:
			return err
		}
	}
	lset := labels.Labels{labels.Label{Name: labels.MetricName, Value: s[:len(s)-1]}}
	hash := lset.Hash()
	lset = sl.reportSampleMutator(lset)
	ref, err := app.Add(lset, t, v)
	switch err {
	case nil:
		sl.cache.addRef(s, ref, lset, hash)
		return nil
	case storage.ErrOutOfOrderSample, storage.ErrDuplicateSampleForTimestamp:
		return nil
	default:
		return err
	}
}
