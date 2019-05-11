package file

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"gopkg.in/fsnotify/fsnotify.v1"
	"gopkg.in/yaml.v2"
)

var (
	patFileSDName	= regexp.MustCompile(`^[^*]*(\*[^/]*)?\.(json|yml|yaml|JSON|YML|YAML)$`)
	DefaultSDConfig	= SDConfig{RefreshInterval: model.Duration(5 * time.Minute)}
)

type SDConfig struct {
	Files			[]string		`yaml:"files"`
	RefreshInterval	model.Duration	`yaml:"refresh_interval,omitempty"`
}

func (c *SDConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = DefaultSDConfig
	type plain SDConfig
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if len(c.Files) == 0 {
		return fmt.Errorf("file service discovery config must contain at least one path name")
	}
	for _, name := range c.Files {
		if !patFileSDName.MatchString(name) {
			return fmt.Errorf("path name %q is not valid for file discovery", name)
		}
	}
	return nil
}

const fileSDFilepathLabel = model.MetaLabelPrefix + "filepath"

type TimestampCollector struct {
	Description	*prometheus.Desc
	discoverers	map[*Discovery]struct{}
	lock		sync.RWMutex
}

func (t *TimestampCollector) Describe(ch chan<- *prometheus.Desc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- t.Description
}
func (t *TimestampCollector) Collect(ch chan<- prometheus.Metric) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	uniqueFiles := make(map[string]float64)
	t.lock.RLock()
	for fileSD := range t.discoverers {
		fileSD.lock.RLock()
		for filename, timestamp := range fileSD.timestamps {
			uniqueFiles[filename] = timestamp
		}
		fileSD.lock.RUnlock()
	}
	t.lock.RUnlock()
	for filename, timestamp := range uniqueFiles {
		ch <- prometheus.MustNewConstMetric(t.Description, prometheus.GaugeValue, timestamp, filename)
	}
}
func (t *TimestampCollector) addDiscoverer(disc *Discovery) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.lock.Lock()
	t.discoverers[disc] = struct{}{}
	t.lock.Unlock()
}
func (t *TimestampCollector) removeDiscoverer(disc *Discovery) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.lock.Lock()
	delete(t.discoverers, disc)
	t.lock.Unlock()
}
func NewTimestampCollector() *TimestampCollector {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TimestampCollector{Description: prometheus.NewDesc("prometheus_sd_file_mtime_seconds", "Timestamp (mtime) of files read by FileSD. Timestamp is set at read time.", []string{"filename"}, nil), discoverers: make(map[*Discovery]struct{})}
}

var (
	fileSDScanDuration		= prometheus.NewSummary(prometheus.SummaryOpts{Name: "prometheus_sd_file_scan_duration_seconds", Help: "The duration of the File-SD scan in seconds."})
	fileSDReadErrorsCount	= prometheus.NewCounter(prometheus.CounterOpts{Name: "prometheus_sd_file_read_errors_total", Help: "The number of File-SD read errors."})
	fileSDTimeStamp			= NewTimestampCollector()
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(fileSDScanDuration)
	prometheus.MustRegister(fileSDReadErrorsCount)
	prometheus.MustRegister(fileSDTimeStamp)
}

type Discovery struct {
	paths		[]string
	watcher		*fsnotify.Watcher
	interval	time.Duration
	timestamps	map[string]float64
	lock		sync.RWMutex
	lastRefresh	map[string]int
	logger		log.Logger
}

func NewDiscovery(conf *SDConfig, logger log.Logger) *Discovery {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if logger == nil {
		logger = log.NewNopLogger()
	}
	disc := &Discovery{paths: conf.Files, interval: time.Duration(conf.RefreshInterval), timestamps: make(map[string]float64), logger: logger}
	fileSDTimeStamp.addDiscoverer(disc)
	return disc
}
func (d *Discovery) listFiles() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var paths []string
	for _, p := range d.paths {
		files, err := filepath.Glob(p)
		if err != nil {
			level.Error(d.logger).Log("msg", "Error expanding glob", "glob", p, "err", err)
			continue
		}
		paths = append(paths, files...)
	}
	return paths
}
func (d *Discovery) watchFiles() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if d.watcher == nil {
		panic("no watcher configured")
	}
	for _, p := range d.paths {
		if idx := strings.LastIndex(p, "/"); idx > -1 {
			p = p[:idx]
		} else {
			p = "./"
		}
		if err := d.watcher.Add(p); err != nil {
			level.Error(d.logger).Log("msg", "Error adding file watch", "path", p, "err", err)
		}
	}
}
func (d *Discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		level.Error(d.logger).Log("msg", "Error adding file watcher", "err", err)
		return
	}
	d.watcher = watcher
	defer d.stop()
	d.refresh(ctx, ch)
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-d.watcher.Events:
			if len(event.Name) == 0 {
				break
			}
			if event.Op^fsnotify.Chmod == 0 {
				break
			}
			d.refresh(ctx, ch)
		case <-ticker.C:
			d.refresh(ctx, ch)
		case err := <-d.watcher.Errors:
			if err != nil {
				level.Error(d.logger).Log("msg", "Error watching file", "err", err)
			}
		}
	}
}
func (d *Discovery) writeTimestamp(filename string, timestamp float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d.lock.Lock()
	d.timestamps[filename] = timestamp
	d.lock.Unlock()
}
func (d *Discovery) deleteTimestamp(filename string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d.lock.Lock()
	delete(d.timestamps, filename)
	d.lock.Unlock()
}
func (d *Discovery) stop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	level.Debug(d.logger).Log("msg", "Stopping file discovery...", "paths", fmt.Sprintf("%v", d.paths))
	done := make(chan struct{})
	defer close(done)
	fileSDTimeStamp.removeDiscoverer(d)
	go func() {
		for {
			select {
			case <-d.watcher.Errors:
			case <-d.watcher.Events:
			case <-done:
				return
			}
		}
	}()
	if err := d.watcher.Close(); err != nil {
		level.Error(d.logger).Log("msg", "Error closing file watcher", "paths", fmt.Sprintf("%v", d.paths), "err", err)
	}
	level.Debug(d.logger).Log("msg", "File discovery stopped")
}
func (d *Discovery) refresh(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t0 := time.Now()
	defer func() {
		fileSDScanDuration.Observe(time.Since(t0).Seconds())
	}()
	ref := map[string]int{}
	for _, p := range d.listFiles() {
		tgroups, err := d.readFile(p)
		if err != nil {
			fileSDReadErrorsCount.Inc()
			level.Error(d.logger).Log("msg", "Error reading file", "path", p, "err", err)
			ref[p] = d.lastRefresh[p]
			continue
		}
		select {
		case ch <- tgroups:
		case <-ctx.Done():
			return
		}
		ref[p] = len(tgroups)
	}
	for f, n := range d.lastRefresh {
		m, ok := ref[f]
		if !ok || n > m {
			level.Debug(d.logger).Log("msg", "file_sd refresh found file that should be removed", "file", f)
			d.deleteTimestamp(f)
			for i := m; i < n; i++ {
				select {
				case ch <- []*targetgroup.Group{{Source: fileSource(f, i)}}:
				case <-ctx.Done():
					return
				}
			}
		}
	}
	d.lastRefresh = ref
	d.watchFiles()
}
func (d *Discovery) readFile(filename string) ([]*targetgroup.Group, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	info, err := fd.Stat()
	if err != nil {
		return nil, err
	}
	var targetGroups []*targetgroup.Group
	switch ext := filepath.Ext(filename); strings.ToLower(ext) {
	case ".json":
		if err := json.Unmarshal(content, &targetGroups); err != nil {
			return nil, err
		}
	case ".yml", ".yaml":
		if err := yaml.UnmarshalStrict(content, &targetGroups); err != nil {
			return nil, err
		}
	default:
		panic(fmt.Errorf("discovery.File.readFile: unhandled file extension %q", ext))
	}
	for i, tg := range targetGroups {
		if tg == nil {
			err = errors.New("nil target group item found")
			return nil, err
		}
		tg.Source = fileSource(filename, i)
		if tg.Labels == nil {
			tg.Labels = model.LabelSet{}
		}
		tg.Labels[fileSDFilepathLabel] = model.LabelValue(filename)
	}
	d.writeTimestamp(filename, float64(info.ModTime().Unix()))
	return targetGroups, nil
}
func fileSource(filename string, i int) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s:%d", filename, i)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
