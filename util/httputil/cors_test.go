package httputil

import (
	"net/http"
	"regexp"
	"testing"
)

func getCORSHandlerFunc() http.Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hf := func(w http.ResponseWriter, r *http.Request) {
		reg := regexp.MustCompile(`^https://foo\.com$`)
		SetCORS(w, reg, r)
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(hf)
}
func TestCORSHandler(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tearDown := setup()
	defer tearDown()
	client := &http.Client{}
	ch := getCORSHandlerFunc()
	mux.Handle("/any_path", ch)
	dummyOrigin := "https://foo.com"
	req, err := http.NewRequest("OPTIONS", server.URL+"/any_path", nil)
	if err != nil {
		t.Error("could not create request")
	}
	req.Header.Set("Origin", dummyOrigin)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("client get failed with unexpected error")
	}
	AccessControlAllowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	if AccessControlAllowOrigin != dummyOrigin {
		t.Fatalf("%q does not match %q", dummyOrigin, AccessControlAllowOrigin)
	}
	req, err = http.NewRequest("OPTIONS", server.URL+"/any_path", nil)
	if err != nil {
		t.Error("could not create request")
	}
	req.Header.Set("Origin", "https://not-foo.com")
	resp, err = client.Do(req)
	if err != nil {
		t.Error("client get failed with unexpected error")
	}
	AccessControlAllowOrigin = resp.Header.Get("Access-Control-Allow-Origin")
	if AccessControlAllowOrigin != "" {
		t.Fatalf("Access-Control-Allow-Origin should not exist but it was set to: %q", AccessControlAllowOrigin)
	}
}
