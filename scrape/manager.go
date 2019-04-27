package scrape

import (
	"reflect"
	"sync"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/storage"
)

type Appendable interface {
	Appender() (storage.Appender, error)
}

func NewManager(logger log.Logger, app Appendable) *Manager {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &Manager{append: app, logger: logger, scrapeConfigs: make(map[string]*config.ScrapeConfig), scrapePools: make(map[string]*scrapePool), graceShut: make(chan struct{}), triggerReload: make(chan struct{}, 1)}
}

type Manager struct {
	logger		log.Logger
	append		Appendable
	graceShut	chan struct{}
	mtxScrape	sync.Mutex
	scrapeConfigs	map[string]*config.ScrapeConfig
	scrapePools	map[string]*scrapePool
	targetSets	map[string][]*targetgroup.Group
	triggerReload	chan struct{}
}

func (m *Manager) Run(tsets <-chan map[string][]*targetgroup.Group) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	go m.reloader()
	for {
		select {
		case ts := <-tsets:
			m.updateTsets(ts)
			select {
			case m.triggerReload <- struct{}{}:
			default:
			}
		case <-m.graceShut:
			return nil
		}
	}
}
func (m *Manager) reloader() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-m.graceShut:
			return
		case <-ticker.C:
			select {
			case <-m.triggerReload:
				m.reload()
			case <-m.graceShut:
				return
			}
		}
	}
}
func (m *Manager) reload() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtxScrape.Lock()
	var wg sync.WaitGroup
	for setName, groups := range m.targetSets {
		var sp *scrapePool
		existing, ok := m.scrapePools[setName]
		if !ok {
			scrapeConfig, ok := m.scrapeConfigs[setName]
			if !ok {
				level.Error(m.logger).Log("msg", "error reloading target set", "err", "invalid config id:"+setName)
				continue
			}
			sp = newScrapePool(scrapeConfig, m.append, log.With(m.logger, "scrape_pool", setName))
			m.scrapePools[setName] = sp
		} else {
			sp = existing
		}
		wg.Add(1)
		go func(sp *scrapePool, groups []*targetgroup.Group) {
			sp.Sync(groups)
			wg.Done()
		}(sp, groups)
	}
	m.mtxScrape.Unlock()
	wg.Wait()
}
func (m *Manager) Stop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtxScrape.Lock()
	defer m.mtxScrape.Unlock()
	for _, sp := range m.scrapePools {
		sp.stop()
	}
	close(m.graceShut)
}
func (m *Manager) updateTsets(tsets map[string][]*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtxScrape.Lock()
	m.targetSets = tsets
	m.mtxScrape.Unlock()
}
func (m *Manager) ApplyConfig(cfg *config.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtxScrape.Lock()
	defer m.mtxScrape.Unlock()
	c := make(map[string]*config.ScrapeConfig)
	for _, scfg := range cfg.ScrapeConfigs {
		c[scfg.JobName] = scfg
	}
	m.scrapeConfigs = c
	for name, sp := range m.scrapePools {
		if cfg, ok := m.scrapeConfigs[name]; !ok {
			sp.stop()
			delete(m.scrapePools, name)
		} else if !reflect.DeepEqual(sp.config, cfg) {
			sp.reload(cfg)
		}
	}
	return nil
}
func (m *Manager) TargetsAll() map[string][]*Target {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtxScrape.Lock()
	defer m.mtxScrape.Unlock()
	targets := make(map[string][]*Target, len(m.scrapePools))
	for tset, sp := range m.scrapePools {
		targets[tset] = append(sp.ActiveTargets(), sp.DroppedTargets()...)
	}
	return targets
}
func (m *Manager) TargetsActive() map[string][]*Target {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtxScrape.Lock()
	defer m.mtxScrape.Unlock()
	targets := make(map[string][]*Target, len(m.scrapePools))
	for tset, sp := range m.scrapePools {
		targets[tset] = sp.ActiveTargets()
	}
	return targets
}
func (m *Manager) TargetsDropped() map[string][]*Target {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mtxScrape.Lock()
	defer m.mtxScrape.Unlock()
	targets := make(map[string][]*Target, len(m.scrapePools))
	for tset, sp := range m.scrapePools {
		targets[tset] = sp.DroppedTargets()
	}
	return targets
}
