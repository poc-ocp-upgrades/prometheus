package runtime

import "runtime"

func Uname() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "(" + runtime.GOOS + ")"
}
