package strutil

import (
	"fmt"
	"net/url"
	"regexp"
)

var (
	invalidLabelCharRE = regexp.MustCompile(`[^a-zA-Z0-9_]`)
)

func TableLinkForExpression(expr string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	escapedExpression := url.QueryEscape(expr)
	return fmt.Sprintf("/graph?g0.expr=%s&g0.tab=1", escapedExpression)
}
func GraphLinkForExpression(expr string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	escapedExpression := url.QueryEscape(expr)
	return fmt.Sprintf("/graph?g0.expr=%s&g0.tab=0", escapedExpression)
}
func SanitizeLabelName(name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return invalidLabelCharRE.ReplaceAllString(name, "_")
}
