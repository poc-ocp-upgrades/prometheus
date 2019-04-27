package stats

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/util/testutil"
	"regexp"
	"testing"
	"time"
)

func TestTimerGroupNewTimer(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tg := NewTimerGroup()
	timer := tg.GetTimer(ExecTotalTime)
	if duration := timer.Duration(); duration != 0 {
		t.Fatalf("Expected duration of 0, but it was %f instead.", duration)
	}
	minimum := 2 * time.Millisecond
	timer.Start()
	time.Sleep(minimum)
	timer.Stop()
	if duration := timer.Duration(); duration == 0 {
		t.Fatalf("Expected duration greater than 0, but it was %f instead.", duration)
	}
	if elapsed := timer.ElapsedTime(); elapsed < minimum {
		t.Fatalf("Expected elapsed time to be greater than time slept, elapsed was %d, and time slept was %d.", elapsed.Nanoseconds(), minimum)
	}
}
func TestQueryStatsWithTimers(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	qt := NewQueryTimers()
	timer := qt.GetTimer(ExecTotalTime)
	timer.Start()
	time.Sleep(2 * time.Millisecond)
	timer.Stop()
	qs := NewQueryStats(qt)
	actual, err := json.Marshal(qs)
	if err != nil {
		t.Fatalf("Unexpected error during serialization: %v", err)
	}
	match, err := regexp.MatchString(`[,{]"execTotalTime":\d+\.\d+[,}]`, string(actual))
	if err != nil {
		t.Fatalf("Unexpected error while matching string: %v", err)
	}
	if !match {
		t.Fatalf("Expected timings with one non-zero entry, but got %s.", actual)
	}
}
func TestQueryStatsWithSpanTimers(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	qt := NewQueryTimers()
	ctx := &testutil.MockContext{DoneCh: make(chan struct{})}
	qst, _ := qt.GetSpanTimer(ctx, ExecQueueTime, prometheus.NewSummary(prometheus.SummaryOpts{}))
	time.Sleep(5 * time.Millisecond)
	qst.Finish()
	qs := NewQueryStats(qt)
	actual, err := json.Marshal(qs)
	if err != nil {
		t.Fatalf("Unexpected error during serialization: %v", err)
	}
	match, err := regexp.MatchString(`[,{]"execQueueTime":\d+\.\d+[,}]`, string(actual))
	if err != nil {
		t.Fatalf("Unexpected error while matching string: %v", err)
	}
	if !match {
		t.Fatalf("Expected timings with one non-zero entry, but got %s.", actual)
	}
}
func TestTimerGroup(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tg := NewTimerGroup()
	execTotalTimer := tg.GetTimer(ExecTotalTime)
	if tg.GetTimer(ExecTotalTime).String() != "Exec total time: 0s" {
		t.Fatalf("Expected string %s, but got %s", "", execTotalTimer.String())
	}
	execQueueTimer := tg.GetTimer(ExecQueueTime)
	if tg.GetTimer(ExecQueueTime).String() != "Exec queue wait time: 0s" {
		t.Fatalf("Expected string %s, but got %s", "", execQueueTimer.String())
	}
	innerEvalTimer := tg.GetTimer(InnerEvalTime)
	if tg.GetTimer(InnerEvalTime).String() != "Inner eval time: 0s" {
		t.Fatalf("Expected string %s, but got %s", "", innerEvalTimer.String())
	}
	queryPreparationTimer := tg.GetTimer(QueryPreparationTime)
	if tg.GetTimer(QueryPreparationTime).String() != "Query preparation time: 0s" {
		t.Fatalf("Expected string %s, but got %s", "", queryPreparationTimer.String())
	}
	resultSortTimer := tg.GetTimer(ResultSortTime)
	if tg.GetTimer(ResultSortTime).String() != "Result sorting time: 0s" {
		t.Fatalf("Expected string %s, but got %s", "", resultSortTimer.String())
	}
	evalTotalTimer := tg.GetTimer(EvalTotalTime)
	if tg.GetTimer(EvalTotalTime).String() != "Eval total time: 0s" {
		t.Fatalf("Expected string %s, but got %s", "", evalTotalTimer.String())
	}
	actual := tg.String()
	expected := "Exec total time: 0s\nExec queue wait time: 0s\nInner eval time: 0s\nQuery preparation time: 0s\nResult sorting time: 0s\nEval total time: 0s\n"
	if actual != expected {
		t.Fatalf("Expected timerGroup string %s, but got %s.", expected, actual)
	}
}
