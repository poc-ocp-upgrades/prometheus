package stats

import (
	"bytes"
	"fmt"
	"sort"
	"time"
)

type Timer struct {
	name		fmt.Stringer
	created		time.Time
	start		time.Time
	duration	time.Duration
}

func (t *Timer) Start() *Timer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.start = time.Now()
	return t
}
func (t *Timer) Stop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.duration += time.Since(t.start)
}
func (t *Timer) ElapsedTime() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return time.Since(t.start)
}
func (t *Timer) Duration() float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.duration.Seconds()
}
func (t *Timer) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s: %s", t.name, t.duration)
}

type TimerGroup struct{ timers map[fmt.Stringer]*Timer }

func NewTimerGroup() *TimerGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TimerGroup{timers: map[fmt.Stringer]*Timer{}}
}
func (t *TimerGroup) GetTimer(name fmt.Stringer) *Timer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if timer, exists := t.timers[name]; exists {
		return timer
	}
	timer := &Timer{name: name, created: time.Now()}
	t.timers[name] = timer
	return timer
}

type Timers []*Timer
type byCreationTimeSorter struct{ Timers }

func (t Timers) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(t)
}
func (t Timers) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t[i], t[j] = t[j], t[i]
}
func (s byCreationTimeSorter) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.Timers[i].created.Before(s.Timers[j].created)
}
func (t *TimerGroup) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	timers := byCreationTimeSorter{}
	for _, timer := range t.timers {
		timers.Timers = append(timers.Timers, timer)
	}
	sort.Sort(timers)
	result := &bytes.Buffer{}
	for _, timer := range timers.Timers {
		fmt.Fprintf(result, "%s\n", timer)
	}
	return result.String()
}
