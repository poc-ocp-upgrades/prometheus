package promql

import "github.com/prometheus/prometheus/pkg/textparse"

const (
	fuzzInteresting	= 1
	fuzzMeh		= 0
	fuzzDiscard	= -1
)

func FuzzParseMetric(in []byte) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := textparse.New(in)
	for p.Next() {
	}
	if p.Err() == nil {
		return fuzzInteresting
	}
	return fuzzMeh
}
func FuzzParseMetricSelector(in []byte) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := ParseMetricSelector(string(in))
	if err == nil {
		return fuzzInteresting
	}
	return fuzzMeh
}
func FuzzParseExpr(in []byte) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := ParseExpr(string(in))
	if err == nil {
		return fuzzInteresting
	}
	return fuzzMeh
}
func FuzzParseStmts(in []byte) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := ParseStmts(string(in))
	if err == nil {
		return fuzzInteresting
	}
	return fuzzMeh
}
