package promql

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/testutil"
)

var (
	minNormal	= math.Float64frombits(0x0010000000000000)
	patSpace	= regexp.MustCompile("[\t ]+")
	patLoad		= regexp.MustCompile(`^load\s+(.+?)$`)
	patEvalInstant	= regexp.MustCompile(`^eval(?:_(fail|ordered))?\s+instant\s+(?:at\s+(.+?))?\s+(.+)$`)
)

const (
	epsilon = 0.000001
)

var testStartTime = time.Unix(0, 0)

type Test struct {
	testutil.T
	cmds		[]testCommand
	storage		storage.Storage
	queryEngine	*Engine
	context		context.Context
	cancelCtx	context.CancelFunc
}

func NewTest(t testutil.T, input string) (*Test, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := &Test{T: t, cmds: []testCommand{}}
	err := test.parse(input)
	test.clear()
	return test, err
}
func newTestFromFile(t testutil.T, filename string) (*Test, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewTest(t, string(content))
}
func (t *Test) QueryEngine() *Engine {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.queryEngine
}
func (t *Test) Queryable() storage.Queryable {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.storage
}
func (t *Test) Context() context.Context {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.context
}
func (t *Test) Storage() storage.Storage {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.storage
}
func raise(line int, format string, v ...interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &ParseErr{Line: line + 1, Err: fmt.Errorf(format, v...)}
}
func parseLoad(lines []string, i int) (int, *loadCmd, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !patLoad.MatchString(lines[i]) {
		return i, nil, raise(i, "invalid load command. (load <step:duration>)")
	}
	parts := patLoad.FindStringSubmatch(lines[i])
	gap, err := model.ParseDuration(parts[1])
	if err != nil {
		return i, nil, raise(i, "invalid step definition %q: %s", parts[1], err)
	}
	cmd := newLoadCmd(time.Duration(gap))
	for i+1 < len(lines) {
		i++
		defLine := lines[i]
		if len(defLine) == 0 {
			i--
			break
		}
		metric, vals, err := parseSeriesDesc(defLine)
		if err != nil {
			if perr, ok := err.(*ParseErr); ok {
				perr.Line = i + 1
			}
			return i, nil, err
		}
		cmd.set(metric, vals...)
	}
	return i, cmd, nil
}
func (t *Test) parseEval(lines []string, i int) (int, *evalCmd, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !patEvalInstant.MatchString(lines[i]) {
		return i, nil, raise(i, "invalid evaluation command. (eval[_fail|_ordered] instant [at <offset:duration>] <query>")
	}
	parts := patEvalInstant.FindStringSubmatch(lines[i])
	var (
		mod	= parts[1]
		at	= parts[2]
		expr	= parts[3]
	)
	_, err := ParseExpr(expr)
	if err != nil {
		if perr, ok := err.(*ParseErr); ok {
			perr.Line = i + 1
			perr.Pos += strings.Index(lines[i], expr)
		}
		return i, nil, err
	}
	offset, err := model.ParseDuration(at)
	if err != nil {
		return i, nil, raise(i, "invalid step definition %q: %s", parts[1], err)
	}
	ts := testStartTime.Add(time.Duration(offset))
	cmd := newEvalCmd(expr, ts, i+1)
	switch mod {
	case "ordered":
		cmd.ordered = true
	case "fail":
		cmd.fail = true
	}
	for j := 1; i+1 < len(lines); j++ {
		i++
		defLine := lines[i]
		if len(defLine) == 0 {
			i--
			break
		}
		if f, err := parseNumber(defLine); err == nil {
			cmd.expect(0, nil, sequenceValue{value: f})
			break
		}
		metric, vals, err := parseSeriesDesc(defLine)
		if err != nil {
			if perr, ok := err.(*ParseErr); ok {
				perr.Line = i + 1
			}
			return i, nil, err
		}
		if len(vals) > 1 {
			return i, nil, raise(i, "expecting multiple values in instant evaluation not allowed")
		}
		cmd.expect(j, metric, vals...)
	}
	return i, cmd, nil
}
func getLines(input string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	lines := strings.Split(input, "\n")
	for i, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "#") {
			l = ""
		}
		lines[i] = l
	}
	return lines
}
func (t *Test) parse(input string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	lines := getLines(input)
	var err error
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if len(l) == 0 {
			continue
		}
		var cmd testCommand
		switch c := strings.ToLower(patSpace.Split(l, 2)[0]); {
		case c == "clear":
			cmd = &clearCmd{}
		case c == "load":
			i, cmd, err = parseLoad(lines, i)
		case strings.HasPrefix(c, "eval"):
			i, cmd, err = t.parseEval(lines, i)
		default:
			return raise(i, "invalid command %q", l)
		}
		if err != nil {
			return err
		}
		t.cmds = append(t.cmds, cmd)
	}
	return nil
}

type testCommand interface{ testCmd() }

func (*clearCmd) testCmd() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*loadCmd) testCmd() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*evalCmd) testCmd() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}

type loadCmd struct {
	gap	time.Duration
	metrics	map[uint64]labels.Labels
	defs	map[uint64][]Point
}

func newLoadCmd(gap time.Duration) *loadCmd {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &loadCmd{gap: gap, metrics: map[uint64]labels.Labels{}, defs: map[uint64][]Point{}}
}
func (cmd loadCmd) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "load"
}
func (cmd *loadCmd) set(m labels.Labels, vals ...sequenceValue) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := m.Hash()
	samples := make([]Point, 0, len(vals))
	ts := testStartTime
	for _, v := range vals {
		if !v.omitted {
			samples = append(samples, Point{T: ts.UnixNano() / int64(time.Millisecond/time.Nanosecond), V: v.value})
		}
		ts = ts.Add(cmd.gap)
	}
	cmd.defs[h] = samples
	cmd.metrics[h] = m
}
func (cmd *loadCmd) append(a storage.Appender) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for h, smpls := range cmd.defs {
		m := cmd.metrics[h]
		for _, s := range smpls {
			if _, err := a.Add(m, s.T, s.V); err != nil {
				return err
			}
		}
	}
	return nil
}

type evalCmd struct {
	expr		string
	start		time.Time
	line		int
	fail, ordered	bool
	metrics		map[uint64]labels.Labels
	expected	map[uint64]entry
}
type entry struct {
	pos	int
	vals	[]sequenceValue
}

func (e entry) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%d: %s", e.pos, e.vals)
}
func newEvalCmd(expr string, start time.Time, line int) *evalCmd {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &evalCmd{expr: expr, start: start, line: line, metrics: map[uint64]labels.Labels{}, expected: map[uint64]entry{}}
}
func (ev *evalCmd) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "eval"
}
func (ev *evalCmd) expect(pos int, m labels.Labels, vals ...sequenceValue) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		ev.expected[0] = entry{pos: pos, vals: vals}
		return
	}
	h := m.Hash()
	ev.metrics[h] = m
	ev.expected[h] = entry{pos: pos, vals: vals}
}
func (ev *evalCmd) compareResult(result Value) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch val := result.(type) {
	case Matrix:
		return fmt.Errorf("received range result on instant evaluation")
	case Vector:
		seen := map[uint64]bool{}
		for pos, v := range val {
			fp := v.Metric.Hash()
			if _, ok := ev.metrics[fp]; !ok {
				return fmt.Errorf("unexpected metric %s in result", v.Metric)
			}
			exp := ev.expected[fp]
			if ev.ordered && exp.pos != pos+1 {
				return fmt.Errorf("expected metric %s with %v at position %d but was at %d", v.Metric, exp.vals, exp.pos, pos+1)
			}
			if !almostEqual(exp.vals[0].value, v.V) {
				return fmt.Errorf("expected %v for %s but got %v", exp.vals[0].value, v.Metric, v.V)
			}
			seen[fp] = true
		}
		for fp, expVals := range ev.expected {
			if !seen[fp] {
				fmt.Println("vector result", len(val), ev.expr)
				for _, ss := range val {
					fmt.Println("    ", ss.Metric, ss.Point)
				}
				return fmt.Errorf("expected metric %s with %v not found", ev.metrics[fp], expVals)
			}
		}
	case Scalar:
		if !almostEqual(ev.expected[0].vals[0].value, val.V) {
			return fmt.Errorf("expected Scalar %v but got %v", val.V, ev.expected[0].vals[0].value)
		}
	default:
		panic(fmt.Errorf("promql.Test.compareResult: unexpected result type %T", result))
	}
	return nil
}

type clearCmd struct{}

func (cmd clearCmd) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "clear"
}
func (t *Test) Run() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, cmd := range t.cmds {
		err := t.exec(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
func (t *Test) exec(tc testCommand) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch cmd := tc.(type) {
	case *clearCmd:
		t.clear()
	case *loadCmd:
		app, err := t.storage.Appender()
		if err != nil {
			return err
		}
		if err := cmd.append(app); err != nil {
			app.Rollback()
			return err
		}
		if err := app.Commit(); err != nil {
			return err
		}
	case *evalCmd:
		q, err := t.queryEngine.NewInstantQuery(t.storage, cmd.expr, cmd.start)
		if err != nil {
			return err
		}
		res := q.Exec(t.context)
		if res.Err != nil {
			if cmd.fail {
				return nil
			}
			return fmt.Errorf("error evaluating query %q (line %d): %s", cmd.expr, cmd.line, res.Err)
		}
		defer q.Close()
		if res.Err == nil && cmd.fail {
			return fmt.Errorf("expected error evaluating query %q (line %d) but got none", cmd.expr, cmd.line)
		}
		err = cmd.compareResult(res.Value)
		if err != nil {
			return fmt.Errorf("error in %s %s: %s", cmd, cmd.expr, err)
		}
		q, err = t.queryEngine.NewRangeQuery(t.storage, cmd.expr, cmd.start.Add(-time.Minute), cmd.start.Add(time.Minute), time.Minute)
		if err != nil {
			return err
		}
		rangeRes := q.Exec(t.context)
		if rangeRes.Err != nil {
			return fmt.Errorf("error evaluating query %q (line %d) in range mode: %s", cmd.expr, cmd.line, rangeRes.Err)
		}
		defer q.Close()
		if cmd.ordered {
			return nil
		}
		mat := rangeRes.Value.(Matrix)
		vec := make(Vector, 0, len(mat))
		for _, series := range mat {
			for _, point := range series.Points {
				if point.T == timeMilliseconds(cmd.start) {
					vec = append(vec, Sample{Metric: series.Metric, Point: point})
					break
				}
			}
		}
		if _, ok := res.Value.(Scalar); ok {
			err = cmd.compareResult(Scalar{V: vec[0].Point.V})
		} else {
			err = cmd.compareResult(vec)
		}
		if err != nil {
			return fmt.Errorf("error in %s %s (line %d) rande mode: %s", cmd, cmd.expr, cmd.line, err)
		}
	default:
		panic("promql.Test.exec: unknown test command type")
	}
	return nil
}
func (t *Test) clear() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.storage != nil {
		if err := t.storage.Close(); err != nil {
			t.T.Fatalf("closing test storage: %s", err)
		}
	}
	if t.cancelCtx != nil {
		t.cancelCtx()
	}
	t.storage = testutil.NewStorage(t)
	opts := EngineOpts{Logger: nil, Reg: nil, MaxConcurrent: 20, MaxSamples: 10000, Timeout: 100 * time.Second}
	t.queryEngine = NewEngine(opts)
	t.context, t.cancelCtx = context.WithCancel(context.Background())
}
func (t *Test) Close() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.cancelCtx()
	if err := t.storage.Close(); err != nil {
		t.T.Fatalf("closing test storage: %s", err)
	}
}
func almostEqual(a, b float64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if a == b {
		return true
	}
	diff := math.Abs(a - b)
	if a == 0 || b == 0 || diff < minNormal {
		return diff < epsilon*minNormal
	}
	return diff/(math.Abs(a)+math.Abs(b)) < epsilon
}
func parseNumber(s string) (float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, err := strconv.ParseInt(s, 0, 64)
	f := float64(n)
	if err != nil {
		f, err = strconv.ParseFloat(s, 64)
	}
	if err != nil {
		return 0, fmt.Errorf("error parsing number: %s", err)
	}
	return f, nil
}

type LazyLoader struct {
	testutil.T
	loadCmd		*loadCmd
	storage		storage.Storage
	queryEngine	*Engine
	context		context.Context
	cancelCtx	context.CancelFunc
}

func NewLazyLoader(t testutil.T, input string) (*LazyLoader, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ll := &LazyLoader{T: t}
	err := ll.parse(input)
	ll.clear()
	return ll, err
}
func (ll *LazyLoader) parse(input string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	lines := getLines(input)
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if len(l) == 0 {
			continue
		}
		if strings.ToLower(patSpace.Split(l, 2)[0]) == "load" {
			_, cmd, err := parseLoad(lines, i)
			if err != nil {
				return err
			}
			ll.loadCmd = cmd
			return nil
		}
		return raise(i, "invalid command %q", l)
	}
	return errors.New("no \"load\" command found")
}
func (ll *LazyLoader) clear() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if ll.storage != nil {
		if err := ll.storage.Close(); err != nil {
			ll.T.Fatalf("closing test storage: %s", err)
		}
	}
	if ll.cancelCtx != nil {
		ll.cancelCtx()
	}
	ll.storage = testutil.NewStorage(ll)
	opts := EngineOpts{Logger: nil, Reg: nil, MaxConcurrent: 20, MaxSamples: 10000, Timeout: 100 * time.Second}
	ll.queryEngine = NewEngine(opts)
	ll.context, ll.cancelCtx = context.WithCancel(context.Background())
}
func (ll *LazyLoader) appendTill(ts int64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	app, err := ll.storage.Appender()
	if err != nil {
		return err
	}
	for h, smpls := range ll.loadCmd.defs {
		m := ll.loadCmd.metrics[h]
		for i, s := range smpls {
			if s.T > ts {
				ll.loadCmd.defs[h] = smpls[i:]
				break
			}
			if _, err := app.Add(m, s.T, s.V); err != nil {
				return err
			}
			if i == len(smpls)-1 {
				ll.loadCmd.defs[h] = nil
			}
		}
	}
	return app.Commit()
}
func (ll *LazyLoader) WithSamplesTill(ts time.Time, fn func(error)) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tsMilli := ts.Sub(time.Unix(0, 0)) / time.Millisecond
	fn(ll.appendTill(int64(tsMilli)))
}
func (ll *LazyLoader) QueryEngine() *Engine {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ll.queryEngine
}
func (ll *LazyLoader) Queryable() storage.Queryable {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ll.storage
}
func (ll *LazyLoader) Context() context.Context {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ll.context
}
func (ll *LazyLoader) Storage() storage.Storage {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ll.storage
}
func (ll *LazyLoader) Close() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ll.cancelCtx()
	if err := ll.storage.Close(); err != nil {
		ll.T.Fatalf("closing test storage: %s", err)
	}
}
