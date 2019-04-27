package textparse

import (
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/value"
)

type openMetricsLexer struct {
	b	[]byte
	i	int
	start	int
	err	error
	state	int
}

func (l *openMetricsLexer) buf() []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return l.b[l.start:l.i]
}
func (l *openMetricsLexer) cur() byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return l.b[l.i]
}
func (l *openMetricsLexer) next() byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l.i++
	if l.i >= len(l.b) {
		l.err = io.EOF
		return byte(tEOF)
	}
	for l.b[l.i] == 0 && (l.state == sLValue || l.state == sMeta2 || l.state == sComment) {
		l.i++
		if l.i >= len(l.b) {
			l.err = io.EOF
			return byte(tEOF)
		}
	}
	return l.b[l.i]
}
func (l *openMetricsLexer) Error(es string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l.err = errors.New(es)
}

type OpenMetricsParser struct {
	l	*openMetricsLexer
	series	[]byte
	text	[]byte
	mtype	MetricType
	val	float64
	ts	int64
	hasTS	bool
	start	int
	offsets	[]int
}

func NewOpenMetricsParser(b []byte) Parser {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &OpenMetricsParser{l: &openMetricsLexer{b: b}}
}
func (p *OpenMetricsParser) Series() ([]byte, *int64, float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if p.hasTS {
		return p.series, &p.ts, p.val
	}
	return p.series, nil, p.val
}
func (p *OpenMetricsParser) Help() ([]byte, []byte) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := p.l.b[p.offsets[0]:p.offsets[1]]
	if strings.IndexByte(yoloString(p.text), byte('\\')) >= 0 {
		return m, []byte(lvalReplacer.Replace(string(p.text)))
	}
	return m, p.text
}
func (p *OpenMetricsParser) Type() ([]byte, MetricType) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.l.b[p.offsets[0]:p.offsets[1]], p.mtype
}
func (p *OpenMetricsParser) Unit() ([]byte, []byte) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.l.b[p.offsets[0]:p.offsets[1]], p.text
}
func (p *OpenMetricsParser) Comment() []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.text
}
func (p *OpenMetricsParser) Metric(l *labels.Labels) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := string(p.series)
	*l = append(*l, labels.Label{Name: labels.MetricName, Value: s[:p.offsets[0]-p.start]})
	for i := 1; i < len(p.offsets); i += 4 {
		a := p.offsets[i] - p.start
		b := p.offsets[i+1] - p.start
		c := p.offsets[i+2] - p.start
		d := p.offsets[i+3] - p.start
		if strings.IndexByte(s[c:d], byte('\\')) >= 0 {
			*l = append(*l, labels.Label{Name: s[a:b], Value: lvalReplacer.Replace(s[c:d])})
			continue
		}
		*l = append(*l, labels.Label{Name: s[a:b], Value: s[c:d]})
	}
	sort.Sort((*l)[1:])
	return s
}
func (p *OpenMetricsParser) nextToken() token {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tok := p.l.Lex()
	return tok
}
func (p *OpenMetricsParser) Next() (Entry, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	p.start = p.l.i
	p.offsets = p.offsets[:0]
	switch t := p.nextToken(); t {
	case tEofWord:
		if t := p.nextToken(); t != tEOF {
			return EntryInvalid, fmt.Errorf("unexpected data after # EOF")
		}
		return EntryInvalid, io.EOF
	case tEOF:
		return EntryInvalid, parseError("unexpected end of data", t)
	case tHelp, tType, tUnit:
		switch t := p.nextToken(); t {
		case tMName:
			p.offsets = append(p.offsets, p.l.start, p.l.i)
		default:
			return EntryInvalid, parseError("expected metric name after HELP", t)
		}
		switch t := p.nextToken(); t {
		case tText:
			if len(p.l.buf()) > 1 {
				p.text = p.l.buf()[1 : len(p.l.buf())-1]
			} else {
				p.text = []byte{}
			}
		default:
			return EntryInvalid, parseError("expected text in HELP", t)
		}
		switch t {
		case tType:
			switch s := yoloString(p.text); s {
			case "counter":
				p.mtype = MetricTypeCounter
			case "gauge":
				p.mtype = MetricTypeGauge
			case "histogram":
				p.mtype = MetricTypeHistogram
			case "gaugehistogram":
				p.mtype = MetricTypeGaugeHistogram
			case "summary":
				p.mtype = MetricTypeSummary
			case "info":
				p.mtype = MetricTypeInfo
			case "stateset":
				p.mtype = MetricTypeStateset
			case "unknown":
				p.mtype = MetricTypeUnknown
			default:
				return EntryInvalid, fmt.Errorf("invalid metric type %q", s)
			}
		case tHelp:
			if !utf8.Valid(p.text) {
				return EntryInvalid, fmt.Errorf("help text is not a valid utf8 string")
			}
		}
		switch t {
		case tHelp:
			return EntryHelp, nil
		case tType:
			return EntryType, nil
		case tUnit:
			m := yoloString(p.l.b[p.offsets[0]:p.offsets[1]])
			u := yoloString(p.text)
			if len(u) > 0 {
				if !strings.HasSuffix(m, u) || len(m) < len(u)+1 || p.l.b[p.offsets[1]-len(u)-1] != '_' {
					return EntryInvalid, fmt.Errorf("unit not a suffix of metric %q", m)
				}
			}
			return EntryUnit, nil
		}
	case tMName:
		p.offsets = append(p.offsets, p.l.i)
		p.series = p.l.b[p.start:p.l.i]
		t2 := p.nextToken()
		if t2 == tBraceOpen {
			if err := p.parseLVals(); err != nil {
				return EntryInvalid, err
			}
			p.series = p.l.b[p.start:p.l.i]
			t2 = p.nextToken()
		}
		if t2 != tValue {
			return EntryInvalid, parseError("expected value after metric", t)
		}
		if p.val, err = strconv.ParseFloat(yoloString(p.l.buf()[1:]), 64); err != nil {
			return EntryInvalid, err
		}
		if math.IsNaN(p.val) {
			p.val = math.Float64frombits(value.NormalNaN)
		}
		p.hasTS = false
		switch p.nextToken() {
		case tLinebreak:
			break
		case tTimestamp:
			p.hasTS = true
			var ts float64
			if ts, err = strconv.ParseFloat(yoloString(p.l.buf()[1:]), 64); err != nil {
				return EntryInvalid, err
			}
			p.ts = int64(ts * 1000)
			if t2 := p.nextToken(); t2 != tLinebreak {
				return EntryInvalid, parseError("expected next entry after timestamp", t)
			}
		default:
			return EntryInvalid, parseError("expected timestamp or new record", t)
		}
		return EntrySeries, nil
	default:
		err = fmt.Errorf("%q %q is not a valid start token", t, string(p.l.cur()))
	}
	return EntryInvalid, err
}
func (p *OpenMetricsParser) parseLVals() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	first := true
	for {
		t := p.nextToken()
		switch t {
		case tBraceClose:
			return nil
		case tComma:
			if first {
				return parseError("expected label name or left brace", t)
			}
			t = p.nextToken()
			if t != tLName {
				return parseError("expected label name", t)
			}
		case tLName:
			if !first {
				return parseError("expected comma", t)
			}
		default:
			if first {
				return parseError("expected label name or left brace", t)
			}
			return parseError("expected comma or left brace", t)
		}
		first = false
		p.offsets = append(p.offsets, p.l.start, p.l.i)
		if t := p.nextToken(); t != tEqual {
			return parseError("expected equal", t)
		}
		if t := p.nextToken(); t != tLValue {
			return parseError("expected label value", t)
		}
		if !utf8.Valid(p.l.buf()) {
			return fmt.Errorf("invalid UTF-8 label value")
		}
		p.offsets = append(p.offsets, p.l.start+1, p.l.i-1)
	}
}
