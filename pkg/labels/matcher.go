package labels

import (
	"fmt"
	"regexp"
)

type MatchType int

const (
	MatchEqual	MatchType	= iota
	MatchNotEqual
	MatchRegexp
	MatchNotRegexp
)

func (m MatchType) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	typeToStr := map[MatchType]string{MatchEqual: "=", MatchNotEqual: "!=", MatchRegexp: "=~", MatchNotRegexp: "!~"}
	if str, ok := typeToStr[m]; ok {
		return str
	}
	panic("unknown match type")
}

type Matcher struct {
	Type	MatchType
	Name	string
	Value	string
	re	*regexp.Regexp
}

func NewMatcher(t MatchType, n, v string) (*Matcher, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := &Matcher{Type: t, Name: n, Value: v}
	if t == MatchRegexp || t == MatchNotRegexp {
		re, err := regexp.Compile("^(?:" + v + ")$")
		if err != nil {
			return nil, err
		}
		m.re = re
	}
	return m, nil
}
func (m *Matcher) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s%s%q", m.Name, m.Type, m.Value)
}
func (m *Matcher) Matches(s string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch m.Type {
	case MatchEqual:
		return s == m.Value
	case MatchNotEqual:
		return s != m.Value
	case MatchRegexp:
		return m.re.MatchString(s)
	case MatchNotRegexp:
		return !m.re.MatchString(s)
	}
	panic("labels.Matcher.Matches: invalid match type")
}
