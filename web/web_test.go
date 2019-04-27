package web

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
	"github.com/prometheus/prometheus/storage/tsdb"
	"github.com/prometheus/prometheus/util/testutil"
	libtsdb "github.com/prometheus/tsdb"
)

func TestMain(m *testing.M) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	os.Setenv("no_proxy", "localhost,127.0.0.1,0.0.0.0,:")
	os.Exit(m.Run())
}
func TestGlobalURL(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	opts := &Options{ListenAddress: ":9090", ExternalURL: &url.URL{Scheme: "https", Host: "externalhost:80", Path: "/path/prefix"}}
	tests := []struct {
		inURL	string
		outURL	string
	}{{inURL: "http://somehost:9090/metrics", outURL: "http://somehost:9090/metrics"}, {inURL: "http://localhost:9090/metrics", outURL: "https://externalhost:80/metrics"}, {inURL: "http://localhost:8000/metrics", outURL: "http://externalhost:8000/metrics"}, {inURL: "http://127.0.0.1:9090/metrics", outURL: "https://externalhost:80/metrics"}}
	for _, test := range tests {
		inURL, err := url.Parse(test.inURL)
		testutil.Ok(t, err)
		globalURL := tmplFuncs("", opts)["globalURL"].(func(u *url.URL) *url.URL)
		outURL := globalURL(inURL)
		testutil.Equals(t, test.outURL, outURL.String())
	}
}
func TestReadyAndHealthy(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Parallel()
	dbDir, err := ioutil.TempDir("", "tsdb-ready")
	testutil.Ok(t, err)
	defer os.RemoveAll(dbDir)
	db, err := libtsdb.Open(dbDir, nil, nil, nil)
	testutil.Ok(t, err)
	opts := &Options{ListenAddress: ":9090", ReadTimeout: 30 * time.Second, MaxConnections: 512, Context: nil, Storage: &tsdb.ReadyStorage{}, QueryEngine: nil, ScrapeManager: nil, RuleManager: nil, Notifier: nil, RoutePrefix: "/", EnableAdminAPI: true, TSDB: func() *libtsdb.DB {
		return db
	}, ExternalURL: &url.URL{Scheme: "http", Host: "localhost:9090", Path: "/"}, Version: &PrometheusVersion{}}
	opts.Flags = map[string]string{}
	webHandler := New(nil, opts)
	go func() {
		err := webHandler.Run(context.Background())
		if err != nil {
			panic(fmt.Sprintf("Can't start web handler:%s", err))
		}
	}()
	time.Sleep(5 * time.Second)
	resp, err := http.Get("http://localhost:9090/-/healthy")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
	resp, err = http.Get("http://localhost:9090/-/ready")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusServiceUnavailable, resp.StatusCode)
	resp, err = http.Get("http://localhost:9090/version")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusServiceUnavailable, resp.StatusCode)
	resp, err = http.Get("http://localhost:9090/graph")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusServiceUnavailable, resp.StatusCode)
	resp, err = http.Post("http://localhost:9090/api/v2/admin/tsdb/snapshot", "", strings.NewReader(""))
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusServiceUnavailable, resp.StatusCode)
	resp, err = http.Post("http://localhost:9090/api/v2/admin/tsdb/delete_series", "", strings.NewReader("{}"))
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusServiceUnavailable, resp.StatusCode)
	webHandler.Ready()
	resp, err = http.Get("http://localhost:9090/-/healthy")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
	resp, err = http.Get("http://localhost:9090/-/ready")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
	resp, err = http.Get("http://localhost:9090/version")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
	resp, err = http.Get("http://localhost:9090/graph")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
	resp, err = http.Post("http://localhost:9090/api/v2/admin/tsdb/snapshot", "", strings.NewReader(""))
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
	resp, err = http.Post("http://localhost:9090/api/v2/admin/tsdb/delete_series", "", strings.NewReader("{}"))
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
}
func TestRoutePrefix(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Parallel()
	dbDir, err := ioutil.TempDir("", "tsdb-ready")
	testutil.Ok(t, err)
	defer os.RemoveAll(dbDir)
	db, err := libtsdb.Open(dbDir, nil, nil, nil)
	testutil.Ok(t, err)
	opts := &Options{ListenAddress: ":9091", ReadTimeout: 30 * time.Second, MaxConnections: 512, Context: nil, Storage: &tsdb.ReadyStorage{}, QueryEngine: nil, ScrapeManager: nil, RuleManager: nil, Notifier: nil, RoutePrefix: "/prometheus", EnableAdminAPI: true, TSDB: func() *libtsdb.DB {
		return db
	}}
	opts.Flags = map[string]string{}
	webHandler := New(nil, opts)
	go func() {
		err := webHandler.Run(context.Background())
		if err != nil {
			panic(fmt.Sprintf("Can't start web handler:%s", err))
		}
	}()
	time.Sleep(5 * time.Second)
	resp, err := http.Get("http://localhost:9091" + opts.RoutePrefix + "/-/healthy")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
	resp, err = http.Get("http://localhost:9091" + opts.RoutePrefix + "/-/ready")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusServiceUnavailable, resp.StatusCode)
	resp, err = http.Get("http://localhost:9091" + opts.RoutePrefix + "/version")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusServiceUnavailable, resp.StatusCode)
	resp, err = http.Post("http://localhost:9091"+opts.RoutePrefix+"/api/v2/admin/tsdb/snapshot", "", strings.NewReader(""))
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusServiceUnavailable, resp.StatusCode)
	resp, err = http.Post("http://localhost:9091"+opts.RoutePrefix+"/api/v2/admin/tsdb/delete_series", "", strings.NewReader("{}"))
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusServiceUnavailable, resp.StatusCode)
	webHandler.Ready()
	resp, err = http.Get("http://localhost:9091" + opts.RoutePrefix + "/-/healthy")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
	resp, err = http.Get("http://localhost:9091" + opts.RoutePrefix + "/-/ready")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
	resp, err = http.Get("http://localhost:9091" + opts.RoutePrefix + "/version")
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
	resp, err = http.Post("http://localhost:9091"+opts.RoutePrefix+"/api/v2/admin/tsdb/snapshot", "", strings.NewReader(""))
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
	resp, err = http.Post("http://localhost:9091"+opts.RoutePrefix+"/api/v2/admin/tsdb/delete_series", "", strings.NewReader("{}"))
	testutil.Ok(t, err)
	testutil.Equals(t, http.StatusOK, resp.StatusCode)
}
func TestDebugHandler(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, tc := range []struct {
		prefix, url	string
		code		int
	}{{"/", "/debug/pprof/cmdline", 200}, {"/foo", "/foo/debug/pprof/cmdline", 200}, {"/", "/debug/pprof/goroutine", 200}, {"/foo", "/foo/debug/pprof/goroutine", 200}, {"/", "/debug/pprof/foo", 404}, {"/foo", "/bar/debug/pprof/goroutine", 404}} {
		opts := &Options{RoutePrefix: tc.prefix}
		handler := New(nil, opts)
		handler.Ready()
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", tc.url, nil)
		testutil.Ok(t, err)
		handler.router.ServeHTTP(w, req)
		testutil.Equals(t, tc.code, w.Code)
	}
}
