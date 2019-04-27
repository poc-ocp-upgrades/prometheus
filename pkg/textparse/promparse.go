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
	"unsafe"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/value"
)

type promlexer struct {
	b	[]byte
	i	int
	start	int
	err	error
	state	int
}
type token int

const (
	tInvalid	token	= -1
	tEOF		token	= 0
	tLinebreak	token	= iota
	tWhitespace
	tHelp
	tType
	tUnit
	tEofWord
	tText
	tComment
	tBlank
	tMName
	tBraceOpen
	tBraceClose
	tLName
	tLValue
	tComma
	tEqual
	tTimestamp
	tValue
)

func (t token) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch t {
	case tInvalid:
		return "INVALID"
	case tEOF:
		return "EOF"
	case tLinebreak:
		return "LINEBREAK"
	case tWhitespace:
		return "WHITESPACE"
	case tHelp:
		return "HELP"
	case tType:
		return "TYPE"
	case tUnit:
		return "UNIT"
	case tEofWord:
		return "EOFWORD"
	case tText:
		return "TEXT"
	case tComment:
		return "COMMENT"
	case tBlank:
		return "BLANK"
	case tMName:
		return "MNAME"
	case tBraceOpen:
		return "BOPEN"
	case tBraceClose:
		return "BCLOSE"
	case tLName:
		return "LNAME"
	case tLValue:
		return "LVALUE"
	case tEqual:
		return "EQUAL"
	case tComma:
		return "COMMA"
	case tTimestamp:
		return "TIMESTAMP"
	case tValue:
		return "VALUE"
	}
	return fmt.Sprintf("<invalid: %d>", t)
}
func (l *promlexer) buf() []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return l.b[l.start:l.i]
}
func (l *promlexer) cur() byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return l.b[l.i]
}
func (l *promlexer) next() byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l.i++
	if l.i >= len(l.b) {
		l.err = io.EOF
		return byte(tEOF)
	}
	for l.b[l.i] == 0 && (l.state == sLValue || l.state == sMeta2 || l.state == sComment) {
		l.i++
	}
	return l.b[l.i]
}
func (l *promlexer) Error(es string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l.err = errors.New(es)
}

type PromParser struct {
	l	*promlexer
	series	[]byte
	text	[]byte
	mtype	MetricType
	val	float64
	ts	int64
	hasTS	bool
	start	int
	offsets	[]int
}

func NewPromParser(b []byte) Parser {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &PromParser{l: &promlexer{b: append(b, '\n')}}
}
func (p *PromParser) Series() ([]byte, *int64, float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if p.hasTS {
		return p.series, &p.ts, p.val
	}
	return p.series, nil, p.val
}
func (p *PromParser) Help() ([]byte, []byte) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := p.l.b[p.offsets[0]:p.offsets[1]]
	if strings.IndexByte(yoloString(p.text), byte('\\')) >= 0 {
		return m, []byte(helpReplacer.Replace(string(p.text)))
	}
	return m, p.text
}
func (p *PromParser) Type() ([]byte, MetricType) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.l.b[p.offsets[0]:p.offsets[1]], p.mtype
}
func (p *PromParser) Unit() ([]byte, []byte) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
func (p *PromParser) Comment() []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.text
}
func (p *PromParser) Metric(l *labels.Labels) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
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
func (p *PromParser) nextToken() token {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for {
		if tok := p.l.Lex(); tok != tWhitespace {
			return tok
		}
	}
}
func parseError(exp string, got token) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Errorf("%s, got %q", exp, got)
}
func (p *PromParser) Next() (Entry, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	p.start = p.l.i
	p.offsets = p.offsets[:0]
	switch t := p.nextToken(); t {
	case tEOF:
		return EntryInvalid, io.EOF
	case tLinebreak:
		return p.Next()
	case tHelp, tType:
		switch t := p.nextToken(); t {
		case tMName:
			p.offsets = append(p.offsets, p.l.start, p.l.i)
		default:
			return EntryInvalid, parseError("expected metric name after HELP", t)
		}
		switch t := p.nextToken(); t {
		case tText:
			if len(p.l.buf()) > 1 {
				p.text = p.l.buf()[1:]
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
			case "summary":
				p.mtype = MetricTypeSummary
			case "untyped":
				p.mtype = MetricTypeUnknown
			default:
				return EntryInvalid, fmt.Errorf("invalid metric type %q", s)
			}
		case tHelp:
			if !utf8.Valid(p.text) {
				return EntryInvalid, fmt.Errorf("help text is not a valid utf8 string")
			}
		}
		if t := p.nextToken(); t != tLinebreak {
			return EntryInvalid, parseError("linebreak expected after metadata", t)
		}
		switch t {
		case tHelp:
			return EntryHelp, nil
		case tType:
			return EntryType, nil
		}
	case tComment:
		p.text = p.l.buf()
		if t := p.nextToken(); t != tLinebreak {
			return EntryInvalid, parseError("linebreak expected after comment", t)
		}
		return EntryComment, nil
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
		if p.val, err = strconv.ParseFloat(yoloString(p.l.buf()), 64); err != nil {
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
			if p.ts, err = strconv.ParseInt(yoloString(p.l.buf()), 10, 64); err != nil {
				return EntryInvalid, err
			}
			if t2 := p.nextToken(); t2 != tLinebreak {
				return EntryInvalid, parseError("expected next entry after timestamp", t)
			}
		default:
			return EntryInvalid, parseError("expected timestamp or new record", t)
		}
		return EntrySeries, nil
	default:
		err = fmt.Errorf("%q is not a valid start token", t)
	}
	return EntryInvalid, err
}
func (p *PromParser) parseLVals() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t := p.nextToken()
	for {
		switch t {
		case tBraceClose:
			return nil
		case tLName:
		default:
			return parseError("expected label name", t)
		}
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
		if t = p.nextToken(); t == tComma {
			t = p.nextToken()
		}
	}
}

var lvalReplacer = strings.NewReplacer(`\"`, "\"", `\\`, "\\", `\n`, "\n")
var helpReplacer = strings.NewReplacer(`\\`, "\\", `\n`, "\n")

func yoloString(b []byte) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return *((*string)(unsafe.Pointer(&b)))
}
