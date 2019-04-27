package strutil

import (
	"testing"
)

type linkTest struct {
	expression		string
	expectedGraphLink	string
	expectedTableLink	string
}

var linkTests = []linkTest{{"sum(incoming_http_requests_total) by (system)", "/graph?g0.expr=sum%28incoming_http_requests_total%29+by+%28system%29&g0.tab=0", "/graph?g0.expr=sum%28incoming_http_requests_total%29+by+%28system%29&g0.tab=1"}, {"sum(incoming_http_requests_total{system=\"trackmetadata\"})", "/graph?g0.expr=sum%28incoming_http_requests_total%7Bsystem%3D%22trackmetadata%22%7D%29&g0.tab=0", "/graph?g0.expr=sum%28incoming_http_requests_total%7Bsystem%3D%22trackmetadata%22%7D%29&g0.tab=1"}}

func TestLink(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, tt := range linkTests {
		if graphLink := GraphLinkForExpression(tt.expression); graphLink != tt.expectedGraphLink {
			t.Errorf("GraphLinkForExpression failed for expression (%#q), want %q got %q", tt.expression, tt.expectedGraphLink, graphLink)
		}
		if tableLink := TableLinkForExpression(tt.expression); tableLink != tt.expectedTableLink {
			t.Errorf("TableLinkForExpression failed for expression (%#q), want %q got %q", tt.expression, tt.expectedTableLink, tableLink)
		}
	}
}
func TestSanitizeLabelName(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	actual := SanitizeLabelName("fooClientLABEL")
	expected := "fooClientLABEL"
	if actual != expected {
		t.Errorf("SanitizeLabelName failed for label (%s), want %s got %s", "fooClientLABEL", expected, actual)
	}
	actual = SanitizeLabelName("barClient.LABEL$$##")
	expected = "barClient_LABEL____"
	if actual != expected {
		t.Errorf("SanitizeLabelName failed for label (%s), want %s got %s", "barClient.LABEL$$##", expected, actual)
	}
}
