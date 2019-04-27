package runtime

import (
	"syscall"
)

func Uname() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	buf := syscall.Utsname{}
	err := syscall.Uname(&buf)
	if err != nil {
		panic("syscall.Uname failed: " + err.Error())
	}
	str := "(" + charsToString(buf.Sysname[:])
	str += " " + charsToString(buf.Release[:])
	str += " " + charsToString(buf.Version[:])
	str += " " + charsToString(buf.Machine[:])
	str += " " + charsToString(buf.Nodename[:])
	str += " " + charsToString(buf.Domainname[:]) + ")"
	return str
}
