package testutil

import (
	"io/ioutil"
	"os"
)

const (
	defaultDirectory		= ""
	NilCloser			= nilCloser(true)
	temporaryDirectoryRemoveRetries	= 2
)

type (
	Closer			interface{ Close() }
	nilCloser		bool
	TemporaryDirectory	interface {
		Closer
		Path() string
	}
	temporaryDirectory	struct {
		path	string
		tester	T
	}
	callbackCloser	struct{ fn func() }
	T		interface {
		Fatal(args ...interface{})
		Fatalf(format string, args ...interface{})
	}
)

func (c nilCloser) Close() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (c callbackCloser) Close() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.fn()
}
func NewCallbackCloser(fn func()) Closer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &callbackCloser{fn: fn}
}
func (t temporaryDirectory) Close() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	retries := temporaryDirectoryRemoveRetries
	err := os.RemoveAll(t.path)
	for err != nil && retries > 0 {
		switch {
		case os.IsNotExist(err):
			err = nil
		default:
			retries--
			err = os.RemoveAll(t.path)
		}
	}
	if err != nil {
		t.tester.Fatal(err)
	}
}
func (t temporaryDirectory) Path() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.path
}
func NewTemporaryDirectory(name string, t T) (handler TemporaryDirectory) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		directory	string
		err		error
	)
	directory, err = ioutil.TempDir(defaultDirectory, name)
	if err != nil {
		t.Fatal(err)
	}
	handler = temporaryDirectory{path: directory, tester: t}
	return
}
