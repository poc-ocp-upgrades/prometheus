package scrape

import (
	"errors"
	"fmt"
	"hash/fnv"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/relabel"
	"github.com/prometheus/prometheus/pkg/textparse"
	"github.com/prometheus/prometheus/pkg/value"
	"github.com/prometheus/prometheus/storage"
)

type TargetHealth string

const (
	HealthUnknown	TargetHealth	= "unknown"
	HealthGood	TargetHealth	= "up"
	HealthBad	TargetHealth	= "down"
)

type Target struct {
	discoveredLabels	labels.Labels
	labels			labels.Labels
	params			url.Values
	mtx			sync.RWMutex
	lastError		error
	lastScrape		time.Time
	lastScrapeDuration	time.Duration
	health			TargetHealth
	metadata		metricMetadataStore
}

func NewTarget(labels, discoveredLabels labels.Labels, params url.Values) *Target {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Target{labels: labels, discoveredLabels: discoveredLabels, params: params, health: HealthUnknown}
}
func (t *Target) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.URL().String()
}

type metricMetadataStore interface {
	listMetadata() []MetricMetadata
	getMetadata(metric string) (MetricMetadata, bool)
}
type MetricMetadata struct {
	Metric	string
	Type	textparse.MetricType
	Help	string
	Unit	string
}

func (t *Target) MetadataList() []MetricMetadata {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	if t.metadata == nil {
		return nil
	}
	return t.metadata.listMetadata()
}
func (t *Target) Metadata(metric string) (MetricMetadata, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	if t.metadata == nil {
		return MetricMetadata{}, false
	}
	return t.metadata.getMetadata(metric)
}
func (t *Target) setMetadataStore(s metricMetadataStore) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.metadata = s
}
func (t *Target) hash() uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%016d", t.labels.Hash())))
	h.Write([]byte(t.URL().String()))
	return h.Sum64()
}
func (t *Target) offset(interval time.Duration) time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	now := time.Now().UnixNano()
	var (
		base	= int64(interval) - now%int64(interval)
		offset	= t.hash() % uint64(interval)
		next	= base + int64(offset)
	)
	if next > int64(interval) {
		next -= int64(interval)
	}
	return time.Duration(next)
}
func (t *Target) Labels() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	lset := make(labels.Labels, 0, len(t.labels))
	for _, l := range t.labels {
		if !strings.HasPrefix(l.Name, model.ReservedLabelPrefix) {
			lset = append(lset, l)
		}
	}
	return lset
}
func (t *Target) DiscoveredLabels() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.mtx.Lock()
	defer t.mtx.Unlock()
	lset := make(labels.Labels, len(t.discoveredLabels))
	copy(lset, t.discoveredLabels)
	return lset
}
func (t *Target) SetDiscoveredLabels(l labels.Labels) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.discoveredLabels = l
}
func (t *Target) URL() *url.URL {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	params := url.Values{}
	for k, v := range t.params {
		params[k] = make([]string, len(v))
		copy(params[k], v)
	}
	for _, l := range t.labels {
		if !strings.HasPrefix(l.Name, model.ParamLabelPrefix) {
			continue
		}
		ks := l.Name[len(model.ParamLabelPrefix):]
		if len(params[ks]) > 0 {
			params[ks][0] = l.Value
		} else {
			params[ks] = []string{l.Value}
		}
	}
	return &url.URL{Scheme: t.labels.Get(model.SchemeLabel), Host: t.labels.Get(model.AddressLabel), Path: t.labels.Get(model.MetricsPathLabel), RawQuery: params.Encode()}
}
func (t *Target) report(start time.Time, dur time.Duration, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if err == nil {
		t.health = HealthGood
	} else {
		t.health = HealthBad
	}
	t.lastError = err
	t.lastScrape = start
	t.lastScrapeDuration = dur
}
func (t *Target) LastError() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	return t.lastError
}
func (t *Target) LastScrape() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	return t.lastScrape
}
func (t *Target) LastScrapeDuration() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	return t.lastScrapeDuration
}
func (t *Target) Health() TargetHealth {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	return t.health
}

type Targets []*Target

func (ts Targets) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(ts)
}
func (ts Targets) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ts[i].URL().String() < ts[j].URL().String()
}
func (ts Targets) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ts[i], ts[j] = ts[j], ts[i]
}

var errSampleLimit = errors.New("sample limit exceeded")

type limitAppender struct {
	storage.Appender
	limit	int
	i	int
}

func (app *limitAppender) Add(lset labels.Labels, t int64, v float64) (uint64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !value.IsStaleNaN(v) {
		app.i++
		if app.i > app.limit {
			return 0, errSampleLimit
		}
	}
	ref, err := app.Appender.Add(lset, t, v)
	if err != nil {
		return 0, err
	}
	return ref, nil
}
func (app *limitAppender) AddFast(lset labels.Labels, ref uint64, t int64, v float64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !value.IsStaleNaN(v) {
		app.i++
		if app.i > app.limit {
			return errSampleLimit
		}
	}
	err := app.Appender.AddFast(lset, ref, t, v)
	return err
}

type timeLimitAppender struct {
	storage.Appender
	maxTime	int64
}

func (app *timeLimitAppender) Add(lset labels.Labels, t int64, v float64) (uint64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t > app.maxTime {
		return 0, storage.ErrOutOfBounds
	}
	ref, err := app.Appender.Add(lset, t, v)
	if err != nil {
		return 0, err
	}
	return ref, nil
}
func (app *timeLimitAppender) AddFast(lset labels.Labels, ref uint64, t int64, v float64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t > app.maxTime {
		return storage.ErrOutOfBounds
	}
	err := app.Appender.AddFast(lset, ref, t, v)
	return err
}
func populateLabels(lset labels.Labels, cfg *config.ScrapeConfig) (res, orig labels.Labels, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	scrapeLabels := []labels.Label{{Name: model.JobLabel, Value: cfg.JobName}, {Name: model.MetricsPathLabel, Value: cfg.MetricsPath}, {Name: model.SchemeLabel, Value: cfg.Scheme}}
	lb := labels.NewBuilder(lset)
	for _, l := range scrapeLabels {
		if lv := lset.Get(l.Name); lv == "" {
			lb.Set(l.Name, l.Value)
		}
	}
	for k, v := range cfg.Params {
		if len(v) > 0 {
			lb.Set(model.ParamLabelPrefix+k, v[0])
		}
	}
	preRelabelLabels := lb.Labels()
	lset = relabel.Process(preRelabelLabels, cfg.RelabelConfigs...)
	if lset == nil {
		return nil, preRelabelLabels, nil
	}
	if v := lset.Get(model.AddressLabel); v == "" {
		return nil, nil, fmt.Errorf("no address")
	}
	lb = labels.NewBuilder(lset)
	addPort := func(s string) bool {
		if _, _, err := net.SplitHostPort(s); err == nil {
			return false
		}
		_, _, err := net.SplitHostPort(s + ":1234")
		return err == nil
	}
	addr := lset.Get(model.AddressLabel)
	if addPort(addr) {
		switch lset.Get(model.SchemeLabel) {
		case "http", "":
			addr = addr + ":80"
		case "https":
			addr = addr + ":443"
		default:
			return nil, nil, fmt.Errorf("invalid scheme: %q", cfg.Scheme)
		}
		lb.Set(model.AddressLabel, addr)
	}
	if err := config.CheckTargetAddress(model.LabelValue(addr)); err != nil {
		return nil, nil, err
	}
	for _, l := range lset {
		if strings.HasPrefix(l.Name, model.MetaLabelPrefix) {
			lb.Del(l.Name)
		}
	}
	if v := lset.Get(model.InstanceLabel); v == "" {
		lb.Set(model.InstanceLabel, addr)
	}
	res = lb.Labels()
	for _, l := range res {
		if !model.LabelValue(l.Value).IsValid() {
			return nil, nil, fmt.Errorf("invalid label value for %q: %q", l.Name, l.Value)
		}
	}
	return res, preRelabelLabels, nil
}
func targetsFromGroup(tg *targetgroup.Group, cfg *config.ScrapeConfig) ([]*Target, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	targets := make([]*Target, 0, len(tg.Targets))
	for i, tlset := range tg.Targets {
		lbls := make([]labels.Label, 0, len(tlset)+len(tg.Labels))
		for ln, lv := range tlset {
			lbls = append(lbls, labels.Label{Name: string(ln), Value: string(lv)})
		}
		for ln, lv := range tg.Labels {
			if _, ok := tlset[ln]; !ok {
				lbls = append(lbls, labels.Label{Name: string(ln), Value: string(lv)})
			}
		}
		lset := labels.New(lbls...)
		lbls, origLabels, err := populateLabels(lset, cfg)
		if err != nil {
			return nil, fmt.Errorf("instance %d in group %s: %s", i, tg, err)
		}
		if lbls != nil || origLabels != nil {
			targets = append(targets, NewTarget(lbls, origLabels, cfg.Params))
		}
	}
	return targets, nil
}
