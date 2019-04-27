package runtime

import (
	"syscall"
)

func VmLimits() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return getLimits(syscall.RLIMIT_DATA, "b")
}
