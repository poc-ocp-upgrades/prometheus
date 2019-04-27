package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/notifier"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/rules"
	"github.com/prometheus/prometheus/util/testutil"
)

var promPath string
var promConfig = filepath.Join("..", "..", "documentation", "examples", "prometheus.yml")
var promData = filepath.Join(os.TempDir(), "data")

func TestMain(m *testing.M) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}
	os.Setenv("no_proxy", "localhost,127.0.0.1,0.0.0.0,:")
	var err error
	promPath, err = os.Getwd()
	if err != nil {
		fmt.Printf("can't get current dir :%s \n", err)
		os.Exit(1)
	}
	promPath = filepath.Join(promPath, "prometheus")
	build := exec.Command("go", "build", "-o", promPath)
	output, err := build.CombinedOutput()
	if err != nil {
		fmt.Printf("compilation error :%s \n", output)
		os.Exit(1)
	}
	exitCode := m.Run()
	os.Remove(promPath)
	os.RemoveAll(promData)
	os.Exit(exitCode)
}
func TestStartupInterrupt(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	prom := exec.Command(promPath, "--config.file="+promConfig, "--storage.tsdb.path="+promData)
	err := prom.Start()
	if err != nil {
		t.Errorf("execution error: %v", err)
		return
	}
	done := make(chan error)
	go func() {
		done <- prom.Wait()
	}()
	var startedOk bool
	var stoppedErr error
Loop:
	for x := 0; x < 10; x++ {
		if _, err := http.Get("http://localhost:9090/graph"); err == nil {
			startedOk = true
			prom.Process.Signal(os.Interrupt)
			select {
			case stoppedErr = <-done:
				break Loop
			case <-time.After(10 * time.Second):
			}
			break Loop
		}
		time.Sleep(500 * time.Millisecond)
	}
	if !startedOk {
		t.Errorf("prometheus didn't start in the specified timeout")
		return
	}
	if err := prom.Process.Kill(); err == nil {
		t.Errorf("prometheus didn't shutdown gracefully after sending the Interrupt signal")
	} else if stoppedErr != nil && stoppedErr.Error() != "signal: interrupt" {
		t.Errorf("prometheus exited with an unexpected error:%v", stoppedErr)
	}
}
func TestComputeExternalURL(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		input	string
		valid	bool
	}{{input: "", valid: true}, {input: "http://proxy.com/prometheus", valid: true}, {input: "'https://url/prometheus'", valid: false}, {input: "'relative/path/with/quotes'", valid: false}, {input: "http://alertmanager.company.com", valid: true}, {input: "https://double--dash.de", valid: true}, {input: "'http://starts/with/quote", valid: false}, {input: "ends/with/quote\"", valid: false}}
	for _, test := range tests {
		_, err := computeExternalURL(test.input, "0.0.0.0:9090")
		if test.valid {
			testutil.Ok(t, err)
		} else {
			testutil.NotOk(t, err, "input=%q", test.input)
		}
	}
}
func TestFailedStartupExitCode(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	fakeInputFile := "fake-input-file"
	expectedExitStatus := 1
	prom := exec.Command(promPath, "--config.file="+fakeInputFile)
	err := prom.Run()
	testutil.NotOk(t, err, "")
	if exitError, ok := err.(*exec.ExitError); ok {
		status := exitError.Sys().(syscall.WaitStatus)
		testutil.Equals(t, expectedExitStatus, status.ExitStatus())
	} else {
		t.Errorf("unable to retrieve the exit status for prometheus: %v", err)
	}
}

type senderFunc func(alerts ...*notifier.Alert)

func (s senderFunc) Send(alerts ...*notifier.Alert) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s(alerts...)
}
func TestSendAlerts(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	testCases := []struct {
		in	[]*rules.Alert
		exp	[]*notifier.Alert
	}{{in: []*rules.Alert{{Labels: []labels.Label{{Name: "l1", Value: "v1"}}, Annotations: []labels.Label{{Name: "a2", Value: "v2"}}, ActiveAt: time.Unix(1, 0), FiredAt: time.Unix(2, 0), ValidUntil: time.Unix(3, 0)}}, exp: []*notifier.Alert{{Labels: []labels.Label{{Name: "l1", Value: "v1"}}, Annotations: []labels.Label{{Name: "a2", Value: "v2"}}, StartsAt: time.Unix(2, 0), EndsAt: time.Unix(3, 0), GeneratorURL: "http://localhost:9090/graph?g0.expr=up&g0.tab=1"}}}, {in: []*rules.Alert{{Labels: []labels.Label{{Name: "l1", Value: "v1"}}, Annotations: []labels.Label{{Name: "a2", Value: "v2"}}, ActiveAt: time.Unix(1, 0), FiredAt: time.Unix(2, 0), ResolvedAt: time.Unix(4, 0)}}, exp: []*notifier.Alert{{Labels: []labels.Label{{Name: "l1", Value: "v1"}}, Annotations: []labels.Label{{Name: "a2", Value: "v2"}}, StartsAt: time.Unix(2, 0), EndsAt: time.Unix(4, 0), GeneratorURL: "http://localhost:9090/graph?g0.expr=up&g0.tab=1"}}}, {in: []*rules.Alert{}}}
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			senderFunc := senderFunc(func(alerts ...*notifier.Alert) {
				if len(tc.in) == 0 {
					t.Fatalf("sender called with 0 alert")
				}
				testutil.Equals(t, tc.exp, alerts)
			})
			sendAlerts(senderFunc, "http://localhost:9090")(context.TODO(), "up", tc.in...)
		})
	}
}
func TestWALSegmentSizeBounds(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	for size, expectedExitStatus := range map[string]int{"9MB": 1, "257MB": 1, "10": 2, "1GB": 1, "12MB": 0} {
		prom := exec.Command(promPath, "--storage.tsdb.wal-segment-size="+size, "--config.file="+promConfig)
		err := prom.Start()
		testutil.Ok(t, err)
		if expectedExitStatus == 0 {
			done := make(chan error, 1)
			go func() {
				done <- prom.Wait()
			}()
			select {
			case err := <-done:
				t.Errorf("prometheus should be still running: %v", err)
			case <-time.After(5 * time.Second):
				prom.Process.Signal(os.Interrupt)
			}
			continue
		}
		err = prom.Wait()
		testutil.NotOk(t, err, "")
		if exitError, ok := err.(*exec.ExitError); ok {
			status := exitError.Sys().(syscall.WaitStatus)
			testutil.Equals(t, expectedExitStatus, status.ExitStatus())
		} else {
			t.Errorf("unable to retrieve the exit status for prometheus: %v", err)
		}
	}
}
func TestChooseRetention(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	retention1, err := model.ParseDuration("20d")
	testutil.Ok(t, err)
	retention2, err := model.ParseDuration("30d")
	testutil.Ok(t, err)
	cases := []struct {
		oldFlagRetention	model.Duration
		newFlagRetention	model.Duration
		chosen			model.Duration
	}{{defaultRetentionDuration, defaultRetentionDuration, defaultRetentionDuration}, {retention1, defaultRetentionDuration, retention1}, {defaultRetentionDuration, retention2, retention2}, {retention1, retention2, retention2}}
	for _, tc := range cases {
		retention := chooseRetention(tc.oldFlagRetention, tc.newFlagRetention)
		testutil.Equals(t, tc.chosen, retention)
	}
}
