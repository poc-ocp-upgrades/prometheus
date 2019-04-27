package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"math"
	"net"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"google.golang.org/grpc"
	template_text "text/template"
	"github.com/cockroachdb/cmux"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/mwitkow/go-conntrack"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
	"github.com/prometheus/common/route"
	"github.com/prometheus/tsdb"
	"golang.org/x/net/netutil"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/notifier"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/rules"
	"github.com/prometheus/prometheus/scrape"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/template"
	"github.com/prometheus/prometheus/util/httputil"
	api_v1 "github.com/prometheus/prometheus/web/api/v1"
	api_v2 "github.com/prometheus/prometheus/web/api/v2"
	"github.com/prometheus/prometheus/web/ui"
)

var localhostRepresentations = []string{"127.0.0.1", "localhost"}

func withStackTracer(h http.Handler, l log.Logger) http.Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				level.Error(l).Log("msg", "panic while serving request", "client", r.RemoteAddr, "url", r.URL, "err", err, "stack", buf)
				panic(err)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

var (
	requestDuration	= prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "prometheus_http_request_duration_seconds", Help: "Histogram of latencies for HTTP requests.", Buckets: []float64{.1, .2, .4, 1, 3, 8, 20, 60, 120}}, []string{"handler"})
	responseSize	= prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "prometheus_http_response_size_bytes", Help: "Histogram of response size for HTTP requests.", Buckets: prometheus.ExponentialBuckets(100, 10, 8)}, []string{"handler"})
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(requestDuration, responseSize)
}

type Handler struct {
	logger		log.Logger
	scrapeManager	*scrape.Manager
	ruleManager	*rules.Manager
	queryEngine	*promql.Engine
	context		context.Context
	tsdb		func() *tsdb.DB
	storage		storage.Storage
	notifier	*notifier.Manager
	apiV1		*api_v1.API
	router		*route.Router
	quitCh		chan struct{}
	reloadCh	chan chan error
	options		*Options
	config		*config.Config
	versionInfo	*PrometheusVersion
	birth		time.Time
	cwd		string
	flagsMap	map[string]string
	mtx		sync.RWMutex
	now		func() model.Time
	ready		uint32
}

func (h *Handler) ApplyConfig(conf *config.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.mtx.Lock()
	defer h.mtx.Unlock()
	h.config = conf
	return nil
}

type PrometheusVersion struct {
	Version		string	`json:"version"`
	Revision	string	`json:"revision"`
	Branch		string	`json:"branch"`
	BuildUser	string	`json:"buildUser"`
	BuildDate	string	`json:"buildDate"`
	GoVersion	string	`json:"goVersion"`
}
type Options struct {
	Context				context.Context
	TSDB				func() *tsdb.DB
	Storage				storage.Storage
	QueryEngine			*promql.Engine
	ScrapeManager			*scrape.Manager
	RuleManager			*rules.Manager
	Notifier			*notifier.Manager
	Version				*PrometheusVersion
	Flags				map[string]string
	ListenAddress			string
	CORSOrigin			*regexp.Regexp
	ReadTimeout			time.Duration
	MaxConnections			int
	ExternalURL			*url.URL
	RoutePrefix			string
	UseLocalAssets			bool
	UserAssetsPath			string
	ConsoleTemplatesPath		string
	ConsoleLibrariesPath		string
	EnableLifecycle			bool
	EnableAdminAPI			bool
	PageTitle			string
	RemoteReadSampleLimit		int
	RemoteReadConcurrencyLimit	int
}

func instrumentHandlerWithPrefix(prefix string) func(handlerName string, handler http.HandlerFunc) http.HandlerFunc {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(handlerName string, handler http.HandlerFunc) http.HandlerFunc {
		return instrumentHandler(prefix+handlerName, handler)
	}
}
func instrumentHandler(handlerName string, handler http.HandlerFunc) http.HandlerFunc {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return promhttp.InstrumentHandlerDuration(requestDuration.MustCurryWith(prometheus.Labels{"handler": handlerName}), promhttp.InstrumentHandlerResponseSize(responseSize.MustCurryWith(prometheus.Labels{"handler": handlerName}), handler))
}
func New(logger log.Logger, o *Options) *Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	router := route.New().WithInstrumentation(instrumentHandler)
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "<error retrieving current working directory>"
	}
	if logger == nil {
		logger = log.NewNopLogger()
	}
	h := &Handler{logger: logger, router: router, quitCh: make(chan struct{}), reloadCh: make(chan chan error), options: o, versionInfo: o.Version, birth: time.Now(), cwd: cwd, flagsMap: o.Flags, context: o.Context, scrapeManager: o.ScrapeManager, ruleManager: o.RuleManager, queryEngine: o.QueryEngine, tsdb: o.TSDB, storage: o.Storage, notifier: o.Notifier, now: model.Now, ready: 0}
	h.apiV1 = api_v1.NewAPI(h.queryEngine, h.storage, h.scrapeManager, h.notifier, func() config.Config {
		h.mtx.RLock()
		defer h.mtx.RUnlock()
		return *h.config
	}, o.Flags, h.testReady, func() api_v1.TSDBAdmin {
		if db := h.options.TSDB(); db != nil {
			return db
		}
		return nil
	}, h.options.EnableAdminAPI, logger, h.ruleManager, h.options.RemoteReadSampleLimit, h.options.RemoteReadConcurrencyLimit, h.options.CORSOrigin)
	if o.RoutePrefix != "/" {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, o.RoutePrefix, http.StatusFound)
		})
		router = router.WithPrefix(o.RoutePrefix)
	}
	readyf := h.testReady
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path.Join(o.ExternalURL.Path, "/graph"), http.StatusFound)
	})
	router.Get("/alerts", readyf(h.alerts))
	router.Get("/graph", readyf(h.graph))
	router.Get("/status", readyf(h.status))
	router.Get("/flags", readyf(h.flags))
	router.Get("/config", readyf(h.serveConfig))
	router.Get("/rules", readyf(h.rules))
	router.Get("/targets", readyf(h.targets))
	router.Get("/version", readyf(h.version))
	router.Get("/service-discovery", readyf(h.serviceDiscovery))
	router.Get("/metrics", promhttp.Handler().ServeHTTP)
	router.Get("/federate", readyf(httputil.CompressionHandler{Handler: http.HandlerFunc(h.federation)}.ServeHTTP))
	router.Get("/consoles/*filepath", readyf(h.consoles))
	router.Get("/static/*filepath", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = path.Join("/static", route.Param(r.Context(), "filepath"))
		fs := http.FileServer(ui.Assets)
		fs.ServeHTTP(w, r)
	})
	if o.UserAssetsPath != "" {
		router.Get("/user/*filepath", route.FileServe(o.UserAssetsPath))
	}
	if o.EnableLifecycle {
		router.Post("/-/quit", h.quit)
		router.Post("/-/reload", h.reload)
	} else {
		router.Post("/-/quit", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Lifecycle APIs are not enabled"))
		})
		router.Post("/-/reload", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Lifecycle APIs are not enabled"))
		})
	}
	router.Get("/-/quit", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST requests allowed"))
	})
	router.Get("/-/reload", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST requests allowed"))
	})
	router.Get("/debug/*subpath", serveDebug)
	router.Post("/debug/*subpath", serveDebug)
	router.Get("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Prometheus is Healthy.\n")
	})
	router.Get("/-/ready", readyf(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Prometheus is Ready.\n")
	}))
	return h
}
func serveDebug(w http.ResponseWriter, req *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx := req.Context()
	subpath := route.Param(ctx, "subpath")
	if subpath == "/pprof" {
		http.Redirect(w, req, req.URL.Path+"/", http.StatusMovedPermanently)
		return
	}
	if !strings.HasPrefix(subpath, "/pprof/") {
		http.NotFound(w, req)
		return
	}
	subpath = strings.TrimPrefix(subpath, "/pprof/")
	switch subpath {
	case "cmdline":
		pprof.Cmdline(w, req)
	case "profile":
		pprof.Profile(w, req)
	case "symbol":
		pprof.Symbol(w, req)
	case "trace":
		pprof.Trace(w, req)
	default:
		req.URL.Path = "/debug/pprof/" + subpath
		pprof.Index(w, req)
	}
}
func (h *Handler) Ready() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	atomic.StoreUint32(&h.ready, 1)
}
func (h *Handler) isReady() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ready := atomic.LoadUint32(&h.ready)
	return ready > 0
}
func (h *Handler) testReady(f http.HandlerFunc) http.HandlerFunc {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(w http.ResponseWriter, r *http.Request) {
		if h.isReady() {
			f(w, r)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Service Unavailable")
		}
	}
}
func (h *Handler) testReadyHandler(f http.Handler) http.HandlerFunc {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return h.testReady(f.ServeHTTP)
}
func (h *Handler) Quit() <-chan struct{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return h.quitCh
}
func (h *Handler) Reload() <-chan chan error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return h.reloadCh
}
func (h *Handler) Run(ctx context.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	level.Info(h.logger).Log("msg", "Start listening for connections", "address", h.options.ListenAddress)
	listener, err := net.Listen("tcp", h.options.ListenAddress)
	if err != nil {
		return err
	}
	listener = netutil.LimitListener(listener, h.options.MaxConnections)
	listener = conntrack.NewListener(listener, conntrack.TrackWithName("http"), conntrack.TrackWithTracing())
	var (
		m	= cmux.New(listener)
		grpcl	= m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
		httpl	= m.Match(cmux.HTTP1Fast())
		grpcSrv	= grpc.NewServer()
	)
	av2 := api_v2.New(h.options.TSDB, h.options.EnableAdminAPI)
	av2.RegisterGRPC(grpcSrv)
	hh, err := av2.HTTPHandler(ctx, h.options.ListenAddress)
	if err != nil {
		return err
	}
	hhFunc := h.testReadyHandler(hh)
	operationName := nethttp.OperationNameFunc(func(r *http.Request) string {
		return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
	})
	mux := http.NewServeMux()
	mux.Handle("/", h.router)
	av1 := route.New().WithInstrumentation(instrumentHandlerWithPrefix("/api/v1"))
	h.apiV1.Register(av1)
	apiPath := "/api"
	if h.options.RoutePrefix != "/" {
		apiPath = h.options.RoutePrefix + apiPath
		level.Info(h.logger).Log("msg", "router prefix", "prefix", h.options.RoutePrefix)
	}
	mux.Handle(apiPath+"/v1/", http.StripPrefix(apiPath+"/v1", av1))
	mux.Handle(apiPath+"/", http.StripPrefix(apiPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httputil.SetCORS(w, h.options.CORSOrigin, r)
		hhFunc(w, r)
	})))
	errlog := stdlog.New(log.NewStdlibAdapter(level.Error(h.logger)), "", 0)
	httpSrv := &http.Server{Handler: withStackTracer(nethttp.Middleware(opentracing.GlobalTracer(), mux, operationName), h.logger), ErrorLog: errlog, ReadTimeout: h.options.ReadTimeout}
	errCh := make(chan error)
	go func() {
		errCh <- httpSrv.Serve(httpl)
	}()
	go func() {
		errCh <- grpcSrv.Serve(grpcl)
	}()
	go func() {
		errCh <- m.Serve()
	}()
	select {
	case e := <-errCh:
		return e
	case <-ctx.Done():
		httpSrv.Shutdown(ctx)
		grpcSrv.GracefulStop()
		return nil
	}
}
func (h *Handler) alerts(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	alerts := h.ruleManager.AlertingRules()
	alertsSorter := byAlertStateAndNameSorter{alerts: alerts}
	sort.Sort(alertsSorter)
	alertStatus := AlertStatus{AlertingRules: alertsSorter.alerts, AlertStateToRowClass: map[rules.AlertState]string{rules.StateInactive: "success", rules.StatePending: "warning", rules.StateFiring: "danger"}}
	h.executeTemplate(w, "alerts.html", alertStatus)
}
func (h *Handler) consoles(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx := r.Context()
	name := route.Param(ctx, "filepath")
	file, err := http.Dir(h.options.ConsoleTemplatesPath).Open(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()
	text, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rawParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := map[string]string{}
	for k, v := range rawParams {
		params[k] = v[0]
	}
	data := struct {
		RawParams	url.Values
		Params		map[string]string
		Path		string
	}{RawParams: rawParams, Params: params, Path: strings.TrimLeft(name, "/")}
	tmpl := template.NewTemplateExpander(h.context, string(text), "__console_"+name, data, h.now(), template.QueryFunc(rules.EngineQueryFunc(h.queryEngine, h.storage)), h.options.ExternalURL)
	filenames, err := filepath.Glob(h.options.ConsoleLibrariesPath + "/*.lib")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := tmpl.ExpandHTML(filenames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, result)
}
func (h *Handler) graph(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.executeTemplate(w, "graph.html", nil)
}
func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	status := struct {
		Birth			time.Time
		CWD			string
		Version			*PrometheusVersion
		Alertmanagers		[]*url.URL
		GoroutineCount		int
		GOMAXPROCS		int
		GOGC			string
		CorruptionCount		int64
		ChunkCount		int64
		TimeSeriesCount		int64
		LastConfigTime		time.Time
		ReloadConfigSuccess	bool
	}{Birth: h.birth, CWD: h.cwd, Version: h.versionInfo, Alertmanagers: h.notifier.Alertmanagers(), GoroutineCount: runtime.NumGoroutine(), GOMAXPROCS: runtime.GOMAXPROCS(0), GOGC: os.Getenv("GOGC")}
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		http.Error(w, fmt.Sprintf("error gathering runtime status: %s", err), http.StatusInternalServerError)
		return
	}
	for _, mF := range metrics {
		switch *mF.Name {
		case "prometheus_tsdb_head_chunks":
			status.ChunkCount = int64(toFloat64(mF))
		case "prometheus_tsdb_head_series":
			status.TimeSeriesCount = int64(toFloat64(mF))
		case "prometheus_tsdb_wal_corruptions_total":
			status.CorruptionCount = int64(toFloat64(mF))
		case "prometheus_config_last_reload_successful":
			status.ReloadConfigSuccess = toFloat64(mF) != 0
		case "prometheus_config_last_reload_success_timestamp_seconds":
			status.LastConfigTime = time.Unix(int64(toFloat64(mF)), 0)
		}
	}
	h.executeTemplate(w, "status.html", status)
}
func toFloat64(f *io_prometheus_client.MetricFamily) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := *f.Metric[0]
	if m.Gauge != nil {
		return m.Gauge.GetValue()
	}
	if m.Counter != nil {
		return m.Counter.GetValue()
	}
	if m.Untyped != nil {
		return m.Untyped.GetValue()
	}
	return math.NaN()
}
func (h *Handler) flags(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.executeTemplate(w, "flags.html", h.flagsMap)
}
func (h *Handler) serveConfig(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.mtx.RLock()
	defer h.mtx.RUnlock()
	h.executeTemplate(w, "config.html", h.config.String())
}
func (h *Handler) rules(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.executeTemplate(w, "rules.html", h.ruleManager)
}
func (h *Handler) serviceDiscovery(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var index []string
	targets := h.scrapeManager.TargetsAll()
	for job := range targets {
		index = append(index, job)
	}
	sort.Strings(index)
	scrapeConfigData := struct {
		Index	[]string
		Targets	map[string][]*scrape.Target
		Active	[]int
		Dropped	[]int
		Total	[]int
	}{Index: index, Targets: make(map[string][]*scrape.Target), Active: make([]int, len(index)), Dropped: make([]int, len(index)), Total: make([]int, len(index))}
	for i, job := range scrapeConfigData.Index {
		scrapeConfigData.Targets[job] = make([]*scrape.Target, 0, len(targets[job]))
		scrapeConfigData.Total[i] = len(targets[job])
		for _, target := range targets[job] {
			if target.Labels().Len() == 0 {
				scrapeConfigData.Dropped[i]++
				if scrapeConfigData.Dropped[i] > 100 {
					continue
				}
			} else {
				scrapeConfigData.Active[i]++
			}
			scrapeConfigData.Targets[job] = append(scrapeConfigData.Targets[job], target)
		}
	}
	h.executeTemplate(w, "service-discovery.html", scrapeConfigData)
}
func (h *Handler) targets(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tps := h.scrapeManager.TargetsActive()
	for _, targets := range tps {
		sort.Slice(targets, func(i, j int) bool {
			iJobLabel := targets[i].Labels().Get(model.JobLabel)
			jJobLabel := targets[j].Labels().Get(model.JobLabel)
			if iJobLabel == jJobLabel {
				return targets[i].Labels().Get(model.InstanceLabel) < targets[j].Labels().Get(model.InstanceLabel)
			}
			return iJobLabel < jJobLabel
		})
	}
	h.executeTemplate(w, "targets.html", struct{ TargetPools map[string][]*scrape.Target }{TargetPools: tps})
}
func (h *Handler) version(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dec := json.NewEncoder(w)
	if err := dec.Encode(h.versionInfo); err != nil {
		http.Error(w, fmt.Sprintf("error encoding JSON: %s", err), http.StatusInternalServerError)
	}
}
func (h *Handler) quit(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	fmt.Fprintf(w, "Requesting termination... Goodbye!")
	close(h.quitCh)
}
func (h *Handler) reload(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	rc := make(chan error)
	h.reloadCh <- rc
	if err := <-rc; err != nil {
		http.Error(w, fmt.Sprintf("failed to reload config: %s", err), http.StatusInternalServerError)
	}
}
func (h *Handler) consolesPath() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if _, err := os.Stat(h.options.ConsoleTemplatesPath + "/index.html"); !os.IsNotExist(err) {
		return h.options.ExternalURL.Path + "/consoles/index.html"
	}
	if h.options.UserAssetsPath != "" {
		if _, err := os.Stat(h.options.UserAssetsPath + "/index.html"); !os.IsNotExist(err) {
			return h.options.ExternalURL.Path + "/user/index.html"
		}
	}
	return ""
}
func tmplFuncs(consolesPath string, opts *Options) template_text.FuncMap {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return template_text.FuncMap{"since": func(t time.Time) time.Duration {
		return time.Since(t) / time.Millisecond * time.Millisecond
	}, "consolesPath": func() string {
		return consolesPath
	}, "pathPrefix": func() string {
		return opts.ExternalURL.Path
	}, "pageTitle": func() string {
		return opts.PageTitle
	}, "buildVersion": func() string {
		return opts.Version.Revision
	}, "stripLabels": func(lset map[string]string, labels ...string) map[string]string {
		for _, ln := range labels {
			delete(lset, ln)
		}
		return lset
	}, "globalURL": func(u *url.URL) *url.URL {
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return u
		}
		for _, lhr := range localhostRepresentations {
			if host == lhr {
				_, ownPort, err := net.SplitHostPort(opts.ListenAddress)
				if err != nil {
					return u
				}
				if port == ownPort {
					u.Scheme = opts.ExternalURL.Scheme
					u.Host = opts.ExternalURL.Host
				} else {
					host, _, err := net.SplitHostPort(opts.ExternalURL.Host)
					if err != nil {
						return u
					}
					u.Host = host + ":" + port
				}
				break
			}
		}
		return u
	}, "numHealthy": func(pool []*scrape.Target) int {
		alive := len(pool)
		for _, p := range pool {
			if p.Health() != scrape.HealthGood {
				alive--
			}
		}
		return alive
	}, "targetHealthToClass": func(th scrape.TargetHealth) string {
		switch th {
		case scrape.HealthUnknown:
			return "warning"
		case scrape.HealthGood:
			return "success"
		default:
			return "danger"
		}
	}, "ruleHealthToClass": func(rh rules.RuleHealth) string {
		switch rh {
		case rules.HealthUnknown:
			return "warning"
		case rules.HealthGood:
			return "success"
		default:
			return "danger"
		}
	}, "alertStateToClass": func(as rules.AlertState) string {
		switch as {
		case rules.StateInactive:
			return "success"
		case rules.StatePending:
			return "warning"
		case rules.StateFiring:
			return "danger"
		default:
			panic("unknown alert state")
		}
	}}
}
func (h *Handler) getTemplate(name string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var tmpl string
	appendf := func(name string) error {
		f, err := ui.Assets.Open(path.Join("/templates", name))
		if err != nil {
			return err
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		tmpl += string(b)
		return nil
	}
	err := appendf("_base.html")
	if err != nil {
		return "", fmt.Errorf("error reading base template: %s", err)
	}
	err = appendf(name)
	if err != nil {
		return "", fmt.Errorf("error reading page template %s: %s", name, err)
	}
	return tmpl, nil
}
func (h *Handler) executeTemplate(w http.ResponseWriter, name string, data interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	text, err := h.getTemplate(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tmpl := template.NewTemplateExpander(h.context, text, name, data, h.now(), template.QueryFunc(rules.EngineQueryFunc(h.queryEngine, h.storage)), h.options.ExternalURL)
	tmpl.Funcs(tmplFuncs(h.consolesPath(), h.options))
	result, err := tmpl.ExpandHTML(nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, result)
}

type AlertStatus struct {
	AlertingRules		[]*rules.AlertingRule
	AlertStateToRowClass	map[rules.AlertState]string
}
type byAlertStateAndNameSorter struct{ alerts []*rules.AlertingRule }

func (s byAlertStateAndNameSorter) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(s.alerts)
}
func (s byAlertStateAndNameSorter) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.alerts[i].State() > s.alerts[j].State() || (s.alerts[i].State() == s.alerts[j].State() && s.alerts[i].Name() < s.alerts[j].Name())
}
func (s byAlertStateAndNameSorter) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.alerts[i], s.alerts[j] = s.alerts[j], s.alerts[i]
}
